package validator

import (
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"math"
	"regexp"
	"strconv"
)

type ValueInterface interface {
	Name(name string) ValueInterface
	NoValidate() ValueInterface
	Min(min int64) ValueInterface
	Max(max int64) ValueInterface
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	UInt() uint
	UInt8() uint8
	UInt16() uint16
	UInt32() uint32
	UInt64() uint64
	String() string
}

func NewValue(data string) Value {
	return Value{data: data}
}

type Value struct {
	min, max   *int64
	noValidate bool
	name       string
	data       string
}

func (v Value) Name(name string) ValueInterface {
	v.name = name
	return v
}

func (v Value) NoValidate() ValueInterface {
	v.noValidate = true
	return v
}

func (v Value) Min(min int64) ValueInterface {
	v.min = &min
	return v
}

func (v Value) Max(max int64) ValueInterface {
	v.max = &max
	return v
}

func (v Value) Int() int {
	return int(v.Int64())
}
func (v Value) Int8() int8 {
	x := v.Int()
	if math.MinInt8 > x || math.MaxInt8 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return int8(x)
}
func (v Value) Int16() int16 {
	x := v.Int()
	if math.MinInt16 > x || math.MaxInt16 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return int16(x)
}
func (v Value) Int32() int32 {
	x := v.Int()
	if math.MinInt32 > x || math.MaxInt32 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return int32(x)
}
func (v Value) Int64() int64 {
	i, err := strconv.ParseInt(v.data, 10, 64)
	if err != nil {
		exception.ThrowMsg(v.name+" 非合法数字", exception.ValidateError)
	}
	if !v.noValidate {
		if v.min != nil && i < *v.min {
			exception.ThrowMsg(fmt.Sprintf("%s不能小于%d", v.name, *v.min), exception.ValidateError)
		}
		if v.max != nil && i > *v.max {
			exception.ThrowMsg(fmt.Sprintf("%s不能大于%d", v.name, *v.max), exception.ValidateError)
		}
	}
	return i
}

func (v Value) UInt() uint {
	return uint(v.UInt64())
}
func (v Value) UInt8() uint8 {
	x := v.UInt()
	if math.MaxUint8 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return uint8(x)
}
func (v Value) UInt16() uint16 {
	x := v.UInt()
	if math.MaxUint16 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return uint16(x)
}
func (v Value) UInt32() uint32 {
	x := v.Int()
	if math.MaxInt32 < x {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	return uint32(x)
}
func (v Value) UInt64() uint64 {
	i, err := strconv.ParseInt(v.data, 10, 64)
	if err != nil {
		exception.ThrowMsg(v.name+" 非合法数字", exception.ValidateError)
	}
	if i < 0 {
		exception.ThrowMsg(fmt.Sprintf("%s超出值范围", v.name), exception.ValidateError)
	}
	if !v.noValidate {
		if v.min != nil && i < *v.min {
			exception.ThrowMsg(fmt.Sprintf("%s不能小于%d", v.name, *v.min), exception.ValidateError)
		}
		if v.max != nil && i > *v.max {
			exception.ThrowMsg(fmt.Sprintf("%s不能大于%d", v.name, *v.max), exception.ValidateError)
		}
	}
	return uint64(i)
}

var (
	specialChar = `\<|\>|\\|/`
	re, _       = regexp.Compile(specialChar)
)

func (v Value) String() string {
	l := len([]rune(v.data))
	if !v.noValidate {
		if re.MatchString(v.data) {
			exception.ThrowMsg(fmt.Sprintf("%s不能包含特殊字符%s", v.name, specialChar), exception.ValidateError)
		}
		xl := int64(l)
		if v.min != nil && xl < *v.min {
			exception.ThrowMsg(fmt.Sprintf("%s长度不能小于%d", v.name, *v.min), exception.ValidateError)
		}
		if v.max != nil && xl > *v.max {
			exception.ThrowMsg(fmt.Sprintf("%s长度不能大于%d", v.name, *v.min), exception.ValidateError)
		}
	}
	return v.data
}

type Nil struct {
	must bool
	name string
}

func NewNil(must bool) Nil {
	return Nil{must: must}
}

func (n Nil) Int() int {
	if n.must {
		exception.ThrowMsg(fmt.Sprintf("%s不能为空", n.name), exception.ValidateError)
	}
	return 0
}
func (n Nil) Int8() int8 {
	return int8(n.Int())
}
func (n Nil) Int16() int16 {
	return int16(n.Int())
}
func (n Nil) Int32() int32 {
	return int32(n.Int())
}
func (n Nil) Int64() int64 {
	return int64(n.Int())
}
func (n Nil) UInt() uint {
	if n.must {
		exception.ThrowMsg(fmt.Sprintf("%s不能为空", n.name), exception.ValidateError)
	}
	return 0
}
func (n Nil) UInt8() uint8 {
	return uint8(n.UInt())
}
func (n Nil) UInt16() uint16 {
	return uint16(n.UInt())
}
func (n Nil) UInt32() uint32 {
	return uint32(n.UInt())
}
func (n Nil) UInt64() uint64 {
	return uint64(n.UInt())
}
func (n Nil) String() string {
	if n.must {
		exception.ThrowMsg(fmt.Sprintf("%s不能为空", n.name), exception.ValidateError)
	}
	return ""
}
func (n Nil) Min(min int64) ValueInterface {
	return n
}
func (n Nil) Max(max int64) ValueInterface {
	return n
}
func (n Nil) Name(name string) ValueInterface {
	n.name = name
	return n
}
func (n Nil) NoValidate() ValueInterface {
	return n
}
