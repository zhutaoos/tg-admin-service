package conv

import (
	"app/tools/logger"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BuiltinT interface {
	~string | ~int | ~uint | ~float32 | ~float64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~bool | ~complex64 | ~complex128
}

// Conv 参考 https://github.com/dablelv/cyan/tree/master/conv
func Conv[T any](a any) (T, error) {
	var t T
	switch any(t).(type) {
	case bool:
		v, err := ToBoolE(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case int:
		v, err := ToIntE(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case int8:
		v, err := ToInt8E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case int16:
		v, err := ToInt16E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case int32:
		v, err := ToInt32E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case int64:
		v, err := ToInt64E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case uint:
		v, err := ToUintE(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case uint8:
		v, err := ToUint8E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case uint16:
		v, err := ToUint16E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case uint32:
		v, err := ToUint32E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case uint64:
		v, err := ToUint64E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
	case float32:
		v, err := ToFloat32E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case float64:
		v, err := ToFloat64E(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case time.Duration:
		v, err := ToDurationE(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	case string:
		v, err := ToStringE(a)
		if err != nil {
			return t, err
		}
		t = any(v).(T)
		break
	default:
		return t, fmt.Errorf("the type %T isn't supported", t)
	}
	return t, nil
}

var errNegativeNotAllowed = errors.New("unable to cast negative value")

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// ToBoolE casts any type to a bool type.
func ToBoolE(a any) (bool, error) {
	a = indirect(a)

	switch b := a.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float64, float32, uintptr, complex64, complex128:
		return !reflect.ValueOf(a).IsZero(), nil
	case string:
		return strconv.ParseBool(a.(string))
	case time.Duration:
		return b != 0, nil
	case json.Number:
		v, err := b.Float64()
		return v != 0, err
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", a, a)
	}
}

// ToIntE casts any type to an int type.
func ToIntE(i any) (int, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return intVal, nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case uint:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint8:
		return int(s), nil
	case float64:
		return int(s), nil
	case float32:
		return int(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return int(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		v, err := s.Int64()
		return int(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}

// ToInt8E casts any type to an int8 type.
func ToInt8E(i any) (int8, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return int8(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		return int8(s), nil
	case int32:
		return int8(s), nil
	case int16:
		return int8(s), nil
	case int8:
		return s, nil
	case uint:
		return int8(s), nil
	case uint64:
		return int8(s), nil
	case uint32:
		return int8(s), nil
	case uint16:
		return int8(s), nil
	case uint8:
		return int8(s), nil
	case float64:
		return int8(s), nil
	case float32:
		return int8(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return int8(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	case json.Number:
		v, err := s.Int64()
		return int8(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int8", i, i)
	}
}

// ToInt16E casts any type to an int16 type.
func ToInt16E(i any) (int16, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return int16(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		return int16(s), nil
	case int32:
		return int16(s), nil
	case int16:
		return s, nil
	case int8:
		return int16(s), nil
	case uint:
		return int16(s), nil
	case uint64:
		return int16(s), nil
	case uint32:
		return int16(s), nil
	case uint16:
		return int16(s), nil
	case uint8:
		return int16(s), nil
	case float64:
		return int16(s), nil
	case float32:
		return int16(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return int16(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	case json.Number:
		v, err := s.Int64()
		return int16(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int16", i, i)
	}
}

// ToInt32E casts any type to an int32 type.
func ToInt32E(i any) (int32, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return int32(intVal), nil
	}

	switch s := i.(type) {
	case int64:
		return int32(s), nil
	case int32:
		return s, nil
	case int16:
		return int32(s), nil
	case int8:
		return int32(s), nil
	case uint:
		return int32(s), nil
	case uint64:
		return int32(s), nil
	case uint32:
		return int32(s), nil
	case uint16:
		return int32(s), nil
	case uint8:
		return int32(s), nil
	case float64:
		return int32(s), nil
	case float32:
		return int32(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return int32(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	case json.Number:
		v, err := s.Int64()
		return int32(v), err
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int32", i, i)
	}
}

// ToInt64E casts any to an int64 type.
func ToInt64E(i any) (int64, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return int64(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		return s, nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return v, nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		return s.Int64()
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	}
}

// ToUintE casts any type to a uint type.
func ToUintE(i any) (uint, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		if intVal < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case uint:
		return s, nil
	case uint64:
		return uint(s), nil
	case uint32:
		return uint(s), nil
	case uint16:
		return uint(s), nil
	case uint8:
		return uint(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint", i, i)
	}
}

// ToUint8E casts any type to a uint type.
func ToUint8E(i any) (uint8, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		if intVal < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case uint:
		return uint8(s), nil
	case uint64:
		return uint8(s), nil
	case uint32:
		return uint8(s), nil
	case uint16:
		return uint8(s), nil
	case uint8:
		return s, nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint8(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint8", i, i)
	}
}

// ToUint16E casts any type to a uint16 type.
func ToUint16E(i any) (uint16, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		if intVal < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case uint:
		return uint16(s), nil
	case uint64:
		return uint16(s), nil
	case uint32:
		return uint16(s), nil
	case uint16:
		return s, nil
	case uint8:
		return uint16(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint16(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint16", i, i)
	}
}

// ToUint32E casts any type to a uint32 type.
func ToUint32E(i any) (uint32, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		if intVal < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case uint:
		return uint32(s), nil
	case uint64:
		return uint32(s), nil
	case uint32:
		return s, nil
	case uint16:
		return uint32(s), nil
	case uint8:
		return uint32(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint32(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint32", i, i)
	}
}

// ToUint64E casts any type to a uint64 type.
func ToUint64E(i any) (uint64, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		if intVal < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(intVal), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case uint:
		return uint64(s), nil
	case uint64:
		return s, nil
	case uint32:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return uint64(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	}
}

// ToFloat32E casts any type to a float32 type.
func ToFloat32E(i any) (float32, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return float32(intVal), nil
	}

	switch s := i.(type) {
	case float64:
		return float32(s), nil
	case float32:
		return s, nil
	case int64:
		return float32(s), nil
	case int32:
		return float32(s), nil
	case int16:
		return float32(s), nil
	case int8:
		return float32(s), nil
	case uint:
		return float32(s), nil
	case uint64:
		return float32(s), nil
	case uint32:
		return float32(s), nil
	case uint16:
		return float32(s), nil
	case uint8:
		return float32(s), nil
	case string:
		v, err := strconv.ParseFloat(s, 32)
		return float32(v), err
	case json.Number:
		v, err := s.Float64()
		return float32(v), err
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float32", i, i)
	}
}

// ToFloat64E casts any type to a float64 type.
func ToFloat64E(i any) (float64, error) {
	i = indirect(i)

	intVal, ok := toInt(i)
	if ok {
		return float64(intVal), nil
	}

	switch s := i.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case string:
		return strconv.ParseFloat(s, 64)
	case json.Number:
		return s.Float64()
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	}
}

// ToStringE casts any type to a string type.
func ToStringE(i any) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(s, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.JSStr:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.Srcset:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// ToDurationE casts any type to time.Duration type.
func ToDurationE(i any) (time.Duration, error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		v, _ := Conv[int64](s)
		return time.Duration(v), nil
	case float32, float64:
		v, _ := Conv[float64](s)
		return time.Duration(v), nil
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			return time.ParseDuration(s)
		}
		return time.ParseDuration(s + "ns")
	case json.Number:
		v, err := s.Float64()
		return time.Duration(v), err
	default:
		return time.Duration(0), fmt.Errorf("unable to cast %#v of type %T to Duration", i, i)
	}
}

// toInt returns the int value of v if v or v's underlying type is an int.
// Note that this will return false for int64 etc. types.
func toInt(v any) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}

// Copied from html/template/content.go.
func indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		// Avoid creating a reflection.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// Copied from html/template/content.go.
func indirectToStringerOrError(a any) any {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// trimZeroDecimal trims the zero decimal.
// example 12.00 to 12 while 12.01 still to be 12.01.
func trimZeroDecimal(s string) string {
	var foundZero bool
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				return s[:i-1]
			}
		case '0':
			foundZero = true
		default:
			return s
		}
	}
	return s
}

// Map2Struct map 转 struct
// m 导入的 map
// s 导出的 struct
// withJsonTag 为 true 根据 json 注解导入
func Map2Struct(m map[string]any, s any, withJsonTag bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)

		var name string
		if withJsonTag {
			name = fi.Tag.Get("json")
		} else {
			name = fi.Name
		}

		if name == "" {
			continue
		}
		switch fi.Type.Name() {
		case "int64":
		case "int32":
		case "int8":
		case "int":
			a, err := Conv[int64](m[name])
			if err == nil {
				v.Field(i).SetInt(a)
			}
			break
		case "uint64":
		case "uint32":
		case "uint8":
		case "uint":
			a, err := Conv[uint64](m[name])
			if err == nil {
				v.Field(i).SetUint(a)
			}
			break
		case "float32":
		case "float64":
			a, err := Conv[float64](m[name])
			if err == nil {
				v.Field(i).SetFloat(a)
			}
			break
		case "string":
			a, err := Conv[string](m[name])
			if err == nil {
				v.Field(i).SetString(a)
			}
			break
		case "bool":
			a, err := Conv[bool](m[name])
			if err == nil {
				v.Field(i).SetBool(a)
			}
			break
		}
	}
}

// Struct2Map struct 转 Map
// s 结构体
// withJsonTag true 导出字段 key 为 json 注解
func Struct2Map(s any, withJsonTag bool) map[string]any {
	out := make(map[string]any)

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil
	}
	t := v.Type()

	if withJsonTag == true {
		for i := 0; i < v.NumField(); i++ {
			fi := t.Field(i)
			key := fi.Tag.Get("json") // 设置 json 再导出
			if key != "" {
				out[key] = v.Field(i).Interface()
			}
		}
	} else {
		for i := 0; i < v.NumField(); i++ {
			fi := t.Field(i)
			key := fi.Tag.Get("json") // 设置 json 再导出
			if key != "" {
				out[fi.Name] = v.Field(i).Interface()
			}
		}
	}

	return out
}

// Map2AnyMap 将任意类型 map 转为 map[string]any
func Map2AnyMap[T any](before map[string]T) map[string]any {
	m := make(map[string]any)
	for k, v := range before {
		var i any
		i = v
		m[k] = i
	}
	return m
}

// InSlice 判断切片中是否存在某元素
func InSlice[T BuiltinT](list []T, item T) (int, T) {
	for k, v := range list {
		if v == item {
			return k, v
		}
	}
	return -1, item
}

// RemoveSlice 从切片中删除指定元素
func RemoveSlice[T BuiltinT](a []T, elem T) []T {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

func Timestamp2Str(timestamp int64) string {
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}
	t := time.Unix(timestamp, 0)
	return t.Format(time.DateTime)
}

// Explode 字符串转切片
func Explode[T BuiltinT](sep, str string) ([]T, error) {
	split := strings.Split(sep, str)
	n := make([]T, 0)
	for _, v := range split {
		conv, err := Conv[T](v)
		if err != nil {
			return nil, err
		}
		n = append(n, conv)
	}
	return n, nil
}

func PostForm[T BuiltinT](ctx *gin.Context, key string) T {
	v := ctx.PostForm(key)
	conv, err := Conv[T](v)
	if err != nil {
		logger.Error("PostForm fail", "key", key, "val", v)
	}
	return conv
}
