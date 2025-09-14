package utils

import (
	"fmt"
	"net"
	"strings"
)

type InfluxDB struct {
	Conn *net.UDPConn
}

func InitInfluxUDP(host string) (*InfluxDB, error) {
	raddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return &InfluxDB{Conn: connection}, nil
}

func (db *InfluxDB) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}

	return nil
}

func (db *InfluxDB) Send(measurement string, tags map[string]string, fields map[string]interface{}, timestamp int64) error {
	if db.Conn == nil {
		return fmt.Errorf("udp conn is nil")
	}

	line := buildLine(measurement, tags, fields, timestamp)

	_, err := db.Conn.Write([]byte(line))
	return err
}

func buildLine(measurement string, tags map[string]string, fields map[string]interface{}, ts int64) string {
	sb := &strings.Builder{}

	// measurement
	sb.WriteString(escape(measurement))

	// tags
	for k, v := range tags {
		sb.WriteByte(',')
		sb.WriteString(escape(k))
		sb.WriteByte('=')
		sb.WriteString(escape(v))
	}

	// fields
	sb.WriteByte(' ')
	first := true
	for k, v := range fields {
		if !first {
			sb.WriteByte(',')
		}
		first = false
		sb.WriteString(escape(k))
		sb.WriteByte('=')

		switch val := v.(type) {
		case string:
			sb.WriteString("\"")
			sb.WriteString(strings.ReplaceAll(val, `"`, `\"`))
			sb.WriteString("\"")
		case int, int64:
			fmt.Fprintf(sb, "%di", v)
		case float32, float64:
			fmt.Fprintf(sb, "%f", v)
		case bool:
			if val {
				sb.WriteString("t")
			} else {
				sb.WriteString("f")
			}
		default:
			fmt.Fprintf(sb, "\"%v\"", v)
		}
	}

	if ts > 0 {
		fmt.Fprintf(sb, " %d", ts)
	}

	return sb.String()
}

func escape(s string) string {
	s = strings.ReplaceAll(s, " ", `\ `)
	s = strings.ReplaceAll(s, ",", `\,`)
	s = strings.ReplaceAll(s, "=", `\=`)
	return s
}
