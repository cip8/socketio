package callback

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	seri "github.com/njones/socketio/serialize"
)

type ErrorWrap func() error

func (fn ErrorWrap) Callback(data ...interface{}) error { return fn() }
func (ErrorWrap) Serialize() (string, error) {
	return "", ErrStubSerialize
}
func (ErrorWrap) Unserialize(string) error {
	return ErrStubUnserialize
}

type FuncString func(string)

func (fn FuncString) Callback(v ...interface{}) error {
	if len(v) == 0 {
		v = append(v, "unknown")
	}
	if val, ok := v[0].(string); ok {
		fn(val)
	} else {
		fn("undefined")
	}
	return nil
}
func (FuncString) Serialize() (string, error) {
	return "", ErrStubSerialize
}
func (FuncString) Unserialize(string) error {
	return ErrStubUnserialize
}

type Wrap struct {
	Func       func() interface{} // func([T]...) error
	Parameters []seri.Serializable
}

func (fn Wrap) Callback(data ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case string:
				err = errors.New(e)
			case error:
				err = e
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = ErrUnknownPanic
			}
		}
	}()

	f := reflect.ValueOf(fn.Func())

	if len(data) != f.Type().NumIn() {
		return ErrInvalidDataInParams
	}

	if len(fn.Parameters) != f.Type().NumIn() {
		return ErrInvalidFuncInParams
	}

	if f.Type().NumOut() != 1 {
		return ErrSingleOutParam
	}

	in := make([]reflect.Value, f.Type().NumIn())
	for i, val := range fn.Parameters {
		switch v := data[i].(type) {
		case string:
			val.Unserialize(v)
		case float64:
			vStr := strconv.FormatFloat(v, 'f', 10, 64)
			vStr = strings.TrimRight(vStr, ".0")
			val.Unserialize(vStr)
		default:
			return ErrBadParamType
		}
		if vv, ok := val.(interface{ Interface() interface{} }); ok {
			in[i] = reflect.ValueOf(vv.Interface())
		} else {
			return ErrInterfaceNotFound
		}
	}

	res := f.Call(in)
	erro := res[0].Interface()
	if erro != nil {
		return erro.(error)
	}

	return nil
}

func (Wrap) Serialize() (string, error) {
	return "", ErrStubSerialize
}
func (Wrap) Unserialize(string) error {
	return ErrStubUnserialize
}
