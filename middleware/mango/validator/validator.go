package validator

import (
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"reflect"
	"regexp"
	"strings"
)

var (
	ErrUnsupported    = exception.New("unsupported type", exception.ValidateError)
	ErrUnknownTag     = exception.New("unknown tag", exception.ValidateError)
	ErrCannotValidate = exception.New("cannot validate unexported struct", exception.ValidateError)
)

// 定义校验函数
type ValidationFunc func(v interface{}, key string, params []string) error

// 校验实体
type Validator struct {
	tagName         string // tag name
	validationFuncs map[string]ValidationFunc
}

func (mv *Validator) ValidateStruct(v interface{}) error {
	return mv.validateStruct(reflect.ValueOf(v))
}

func (v *Validator) Engine() interface{} {
	return defaultValidator
}

var defaultValidator = NewValidator()

func NewValidator() *Validator {
	return &Validator{
		tagName: "validate",
		validationFuncs: map[string]ValidationFunc{
			"integer":    integer,
			"varchar":    varchar,
			"phone":      phone,
			"numberchar": numberchar,
			"bankcard":   bankcard,
			"date":       date,
			"regex":      regex,
		},
	}
}

// 设置tag
func SetTag(tag string) {
	defaultValidator.SetTag(tag)
}

func (mv *Validator) SetTag(tag string) {
	mv.tagName = tag
}

// 使用tag校验，不改变原tag
func WithTag(tag string) *Validator {
	return defaultValidator.WithTag(tag)
}

func (mv *Validator) WithTag(tag string) *Validator {
	newFuncs := map[string]ValidationFunc{}
	for k, f := range mv.validationFuncs {
		newFuncs[k] = f
	}
	v := &Validator{
		tagName:         mv.tagName,
		validationFuncs: newFuncs,
	}
	v.SetTag(tag)
	return v
}

// 设置校验方法
func SetValidationFunc(name string, vf ValidationFunc) error {
	return defaultValidator.SetValidationFunc(name, vf)
}

func (mv *Validator) SetValidationFunc(name string, vf ValidationFunc) error {
	if name == "" {
		return exception.New("name cannot be empty", exception.RuntimeError)
	}
	if vf == nil {
		delete(mv.validationFuncs, name)
		return nil
	}
	mv.validationFuncs[name] = vf
	return nil
}

// 校验
func Validate(v interface{}) {
	if err := defaultValidator.ValidateStruct(v); err != nil {
		exception.Throw(err)
	}
}

func (mv *Validator) validateStruct(sv reflect.Value) error {
	kind := sv.Kind()
	if (kind == reflect.Ptr || kind == reflect.Interface) && !sv.IsNil() {
		return mv.validateStruct(sv.Elem())
	}
	if kind != reflect.Struct && kind != reflect.Interface {
		// 不支持的验证类型
		return ErrUnsupported
	}

	st := sv.Type()
	nfields := st.NumField()
	for i := 0; i < nfields; i++ {
		// 遍历struct
		if err := mv.validateField(st.Field(i), sv.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

// validateField validates the field of fieldVal referred to by fieldDef.
// If fieldDef refers to an anonymous/embedded field,
// validateField will walk all of the embedded type's fields and validate them on sv.
func (mv *Validator) validateField(fieldDef reflect.StructField, fieldVal reflect.Value) error {
	tag := fieldDef.Tag.Get(mv.tagName)
	if tag == "-" {
		return nil
	}
	// deal with pointers
	for (fieldVal.Kind() == reflect.Ptr || fieldVal.Kind() == reflect.Interface) && !fieldVal.IsNil() {
		fieldVal = fieldVal.Elem()
	}

	// ignore private structs unless Anonymous
	if !fieldDef.Anonymous && fieldDef.PkgPath != "" {
		return nil
	}

	if tag != "" {
		var err error
		if fieldDef.PkgPath != "" {
			err = ErrCannotValidate
		} else {
			err = mv.validValue(fieldVal, tag)
		}
		if err != nil {
			return err
		}
	}

	// no-op if field is not a struct, interface, array, slice or map
	return mv.deepValidateCollection(fieldVal, func() string {
		return fieldDef.Name
	})
}

func (mv *Validator) deepValidateCollection(f reflect.Value, fnameFn func() string) error {
	switch f.Kind() {
	case reflect.Interface, reflect.Ptr:
		if f.IsNil() {
			return nil
		}
		if err := mv.deepValidateCollection(f.Elem(), fnameFn); err != nil {
			return err
		}
	case reflect.Struct:
		err := mv.validateStruct(f)
		if err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		// we don't need to loop over every byte in a byte slice so we only end up
		// looping when the kind is something we care about
		switch f.Type().Elem().Kind() {
		case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice:
			for i := 0; i < f.Len(); i++ {
				if err := mv.deepValidateCollection(f.Index(i), func() string {
					return fmt.Sprintf("%s[%d]", fnameFn(), i)
				}); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		for _, key := range f.MapKeys() {
			if err := mv.deepValidateCollection(key, func() string {
				return fmt.Sprintf("%s[%+v](key)", fnameFn(), key.Interface())
			}); err != nil {
				return err
			}
			if err := mv.deepValidateCollection(f.MapIndex(key), func() string {
				return fmt.Sprintf("%s[%+v](value)", fnameFn(), key.Interface())
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// validValue is like Valid but takes a Value instead of an interface
func (mv *Validator) validValue(v reflect.Value, tags string) error {
	if v.Kind() == reflect.Invalid {
		return mv.validateVar(nil, tags)
	}
	return mv.validateVar(v.Interface(), tags)
}

// validateVar validates one single variable
func (mv *Validator) validateVar(v interface{}, tag string) error {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown tag found, give up.
		return err
	}
	for _, t := range tags {
		if err := t.Fn(v, t.Key, t.Param); err != nil {
			return err
		}
	}
	return nil
}

// tag represents one of the tag items
type tag struct {
	Name  string // name of the tag
	Key   string
	Fn    ValidationFunc // validation function to call
	Param []string       // parameter to send to the validation function
}

// separate by no escaped pipe
var sepPattern *regexp.Regexp = regexp.MustCompile(`((?:^|[^\\])(?:\\\\)*)\|`)

func splitUnescapedPipe(str string) []string {
	ret := []string{}
	indexes := sepPattern.FindAllStringIndex(str, -1)
	last := 0
	for _, is := range indexes {
		ret = append(ret, str[last:is[1]-1])
		last = is[1]
	}
	ret = append(ret, str[last:])
	return ret
}

// parseTags parses all individual tags found within a struct tag.
func (mv *Validator) parseTags(t string) ([]tag, error) {
	tl := splitUnescapedPipe(t)
	tags := make([]tag, 0, len(tl))
	for _, i := range tl {
		tg := tag{}
		v := strings.SplitN(i, "=", 2)
		tg.Name = strings.Trim(v[0], " ")
		if tg.Name == "" {
			return []tag{}, ErrUnknownTag
		}
		if len(v) > 1 {
			ps := strings.Split(v[1], ",")
			for pi, p := range ps {
				p = strings.Trim(p, " ")
				if pi == 0 {
					tg.Key = p
				} else {
					tg.Param = append(tg.Param, p)
				}
			}
		}
		var found bool
		if tg.Fn, found = mv.validationFuncs[tg.Name]; !found {
			return []tag{}, ErrUnknownTag
		}
		tags = append(tags, tg)
	}
	return tags, nil
}
