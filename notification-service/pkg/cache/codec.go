package cache

import (
	"github.com/vmihailenco/msgpack/v5"
)

func Pack(obj any) ([]byte, error) {
	switch obj.(type) {
	case string, []byte, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return msgpack.Marshal(obj)
	default:
		key := KeyOf(obj)
		payload, err := msgpack.Marshal(obj)
		if err != nil {
			return nil, err
		}
		env := Envelope{T: key, V: payload}
		return msgpack.Marshal(&env)
	}
}

func Unpack(raw []byte) (any, bool, error) {
	var env Envelope
	if err := msgpack.Unmarshal(raw, &env); err == nil && (env.T != "" || len(env.V) > 0) {
		if env.T != "" {
			if dst, ok := NewPtrByKey(env.T); ok {
				if err := msgpack.Unmarshal(env.V, dst); err != nil {
					return nil, false, err
				}

				return dst, true, nil
			}
			var v any
			if err := msgpack.Unmarshal(env.V, &v); err != nil {
				return nil, false, err
			}
			return v, false, nil
		}
		var v any
		if err := msgpack.Unmarshal(env.V, &v); err != nil {
			return nil, false, err
		}
		return v, false, nil
	}

	var v any
	if err := msgpack.Unmarshal(raw, &v); err != nil {
		return nil, false, err
	}
	return v, false, nil
}
