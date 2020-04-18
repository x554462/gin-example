package validator

import (
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func integer(v interface{}, key string, params []string) error {
	pLen := len(params)
	var min, max int64 = 0, 2147483647
	if pLen > 0 {
		min, _ = strconv.ParseInt(params[0], 10, 64)
	}
	if pLen > 1 {
		max, _ = strconv.ParseInt(params[1], 10, 64)
	}
	err := exception.New(fmt.Sprintf("%s必须为整数(%d~%d)", key, min, max), exception.ValidateError)
	var number int64
	switch integer := v.(type) {
	case int, int64, uint, uint64, int32, int16, int8, uint32, uint16, uint8:
		number, _ = strconv.ParseInt(fmt.Sprintf("%d", integer), 10, 64)
	case string:
		var parseErr error
		if number, parseErr = strconv.ParseInt(integer, 10, 64); parseErr != nil {
			return exception.New(parseErr.Error(), exception.ValidateError)
		}
	default:
		return err
	}
	if number < min || number > max {
		return err
	}
	return nil
}

func varchar(v interface{}, key string, params []string) error {
	varchar, ok := v.(string)
	var min, max int64 = 0, 4096
	pLen := len(params)
	if pLen > 0 {
		min, _ = strconv.ParseInt(params[0], 10, 64)
	}
	if pLen > 1 {
		max, _ = strconv.ParseInt(params[1], 10, 64)
	}
	err := exception.New(fmt.Sprintf("%s必须为字符串，长度为%d~%d", key, min, max), exception.ValidateError)
	if !ok {
		return err
	}
	count := int64(utf8.RuneCountInString(varchar))
	if count < min || count > max {
		return err
	}
	return nil
}

var phonePattern *regexp.Regexp = regexp.MustCompile(`/1[3456789]{1}\d{9}$/`)

func phone(v interface{}, key string, params []string) error {
	phone, ok := v.(string)
	if key == "" {
		key = "手机号"
	}
	err := exception.New(fmt.Sprintf("%s错误，请填写正确的手机号", key), exception.ValidateError)
	if !ok {
		return err
	}
	if phonePattern.MatchString(phone) {
		return nil
	}
	return err
}

func numberchar(v interface{}, key string, params []string) error {
	numberchar, ok := v.(string)
	var min, max int64 = 0, 4096
	pLen := len(params)
	if pLen > 0 {
		min, _ = strconv.ParseInt(params[0], 10, 64)
	}
	if pLen > 1 {
		max, _ = strconv.ParseInt(params[1], 10, 64)
	}
	err := exception.New(fmt.Sprintf("%s必须为数字字符，长度为%d~%d", key, min, max), exception.ValidateError)
	if !ok {
		return err
	}
	var counter int
	mmin := int(min)
	mmax := int(max)
	for _, chr := range numberchar {
		if chr < '0' || chr > '9' {
			return err
		}
		counter += 1
		if counter >= mmax {
			return err
		}
	}
	if counter < mmin {
		return err
	}
	return nil
}

func bankcard(v interface{}, key string, params []string) error {
	bankcard, ok := v.(string)
	if key == "" {
		key = "银行卡"
	}
	err := exception.New(fmt.Sprintf("%s错误，请填写正确的银行卡号", key), exception.ValidateError)
	if !ok {
		return err
	}
	var numbers = make([]int32, 0)
	for _, chr := range bankcard {
		numbers = append(numbers, chr-'0')
	}
	arr := []int32{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
	l := len(numbers)
	if l > 19 || l < 15 {
		return err
	}
	var totalNum, num, counter int32 = 0, 0, 0
	for l -= 1; l >= 0; l -= 1 {
		num = numbers[l]
		counter += 1
		if counter%2 == 0 {
			num = arr[num]
		}
		totalNum += num
	}
	if totalNum%10 == 0 {
		return nil
	}
	return err
}

func date(v interface{}, key string, params []string) error {
	date, ok := v.(string)
	if key == "" {
		key = "日期"
	}
	err := exception.New(fmt.Sprintf("%s格式错误", key), exception.ValidateError)
	if !ok {
		return err
	}
	dateArr := strings.Split(date, "-")
	if len(dateArr) != 3 {
		return err
	}
	year, convErr := strconv.ParseInt(dateArr[0], 10, 32)
	if convErr != nil {
		return err
	}
	month, convErr := strconv.ParseInt(dateArr[1], 10, 32)
	if convErr != nil {
		return err
	}
	day, convErr := strconv.ParseInt(dateArr[2], 10, 32)
	if convErr != nil {
		return err
	}
	if year < 2000 || month < 0 || month > 12 || day < 0 || day > 31 {
		return err
	}
	return nil
}

func regex(v interface{}, key string, params []string) error {
	s, ok := v.(string)
	err := exception.New(fmt.Sprintf("%s格式错误", key), exception.ValidateError)
	if !ok {
		sptr, ok := v.(*string)
		if !ok {
			return err
		}
		if sptr == nil {
			return nil
		}
		s = *sptr
	}
	if len(params) < 1 {
		return err
	}

	re, regErr := regexp.Compile(params[0])
	if regErr != nil {
		return err
	}

	if !re.MatchString(s) {
		return err
	}
	return nil
}
