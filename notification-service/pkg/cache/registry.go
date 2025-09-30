package cache

import (
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
	"sync"
)

type Envelope struct {
	T string             `msgpack:"t"`
	V msgpack.RawMessage `msgpack:"v"`
}

var reg sync.Map

func keyForType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}

func Register(ptr any) {
	rt := reflect.TypeOf(ptr)
	if rt == nil || rt.Kind() != reflect.Ptr {
		panic("cache.Register: expected nil pointer, e.g. (*MyType)(nil)")
	}
	t := rt.Elem()
	key := t.PkgPath() + "." + t.Name()
	reg.Store(key, t)
}

func NewPtrByKey(key string) (any, bool) {
	if v, ok := reg.Load(key); ok {
		t := v.(reflect.Type)
		return reflect.New(t).Interface(), true
	}
	return nil, false
}

func KeyOf(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + "." + t.Name()
}
