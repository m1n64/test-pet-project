package cache

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"strconv"
)

type MultiCacheValue struct {
	v any
}

func (m *MultiCacheValue) Bytes() []byte {
	return m.v.([]byte)
}

func (m *MultiCacheValue) String() string {
	switch v := m.v.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case msgpack.RawMessage:
		var x any
		if err := msgpack.Unmarshal([]byte(v), &x); err == nil {
			switch s := x.(type) {
			case string:
				return s
			default:
				return fmt.Sprint(s)
			}
		}
		return string(v)
	default:
		return fmt.Sprint(v)
	}
}

func (m *MultiCacheValue) Int() int64 {
	switch v := m.v.(type) {
	case nil:
		return 0
	case int:
		return int64(v)
	case int8, int16, int32, int64:
		return reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(v).Uint())
	case float32, float64:
		return int64(reflect.ValueOf(v).Float())
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	case []byte:
		n, _ := strconv.ParseInt(string(v), 10, 64)
		return n
	default:
		n, _ := strconv.ParseInt(fmt.Sprint(v), 10, 64)
		return n
	}
}

func (m *MultiCacheValue) Interface() interface{} {
	return m.v
}
