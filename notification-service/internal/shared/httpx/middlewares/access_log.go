package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"mime"
	"net/http"
	"net/url"
	"notification-service-api/internal/shared/httpx"
	"strconv"
	"strings"
	"time"
)

const maxBodyLogBytes = 2048

type bodyCaptureWriter struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqBody := snapshotRequestBody(c.Request)
		bw := &bodyCaptureWriter{ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		latency := time.Since(start)
		latencyMs := latency.Milliseconds()

		respBody := truncateForLog(bw.buf.Bytes())

		msg := "HTTP " + strconv.Itoa(status) +
			" in " + latency.Truncate(time.Millisecond).String() +
			" (" + strconv.Itoa(size) + "B)" +
			" | in: " + reqBody +
			" | out: " + respBody

		log := httpx.FromLogger(c)
		fields := []zap.Field{
			zap.Int("status", status),
			zap.Int("size_bytes", size),
			zap.Duration("latency", latency),
			zap.Int64("latency_ms", latencyMs),
		}

		switch {
		case status >= 500:
			log.Error(msg, fields...)
		case status >= 400:
			log.Warn(msg, fields...)
		default:
			log.Info(msg, fields...)
		}
	}
}

func snapshotRequestBody(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	ctype := r.Header.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ctype)

	switch mt {
	case "application/json", "text/plain", "":
		limited := io.LimitReader(r.Body, maxBodyLogBytes+1)
		b, _ := io.ReadAll(limited)
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(b))
		return truncateForLog(b)

	case "application/x-www-form-urlencoded":
		_ = r.ParseForm()
		return truncateForLog([]byte(redactQuery(r.Form).Encode()))

	default:
		return ""
	}
}

func truncateForLog(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if len(b) > maxBodyLogBytes {
		b = b[:maxBodyLogBytes]
	}
	s := string(b)
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return redactSecrets(s)
}

func redactSecrets(s string) string {
	replacements := []string{
		`"password":"`, `"password":"***`,
		`"passwd":"`, `"passwd":"***`,
		`"token":"`, `"token":"***`,
		`"authorization":"`, `"authorization":"***`,
	}
	r := strings.NewReplacer(replacements...)
	return r.Replace(s)
}

func redactQuery(v url.Values) url.Values {
	out := url.Values{}
	for k, vals := range v {
		kl := strings.ToLower(k)
		if kl == "password" || kl == "passwd" || kl == "token" || kl == "authorization" {
			out[k] = []string{"***"}
			continue
		}
		out[k] = vals
	}
	return out
}
