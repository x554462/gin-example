package exception

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

func New(msg string, parent error) ErrorWrap {
	return ErrorWrap{Msg: msg, Err: parent}
}

func CheckError(err error) {
	if err != nil {
		if !errors.Is(err, RootError) {
			err = New(err.Error(), RootError)
		}
		throwWithCallerDepth(err, 2)
	}
}

func Throw(err error) {
	throwWithCallerDepth(err, 2)
}

func throwWithCallerDepth(err error, callerDepth int) {
	var fileLine string
	if _, file, line, ok := runtime.Caller(callerDepth); ok {
		fileLine = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	if e, ok := err.(ErrorWrapInterface); ok {
		err = ErrorWrap{Msg: e.Error(), Err: e.Unwrap(), Code: e.GetCode(), fileLine: fileLine}
	} else {
		err = ErrorWrap{Msg: err.Error(), Err: RootError, fileLine: fileLine}
	}
	panic(err)
}

func ThrowMsg(msg string, parent error) {
	throwWithCallerDepth(New(msg, parent), 2)
}

func ThrowWithCallerDepth(err error, callerDepth int) {
	throwWithCallerDepth(err, callerDepth)
}

func ThrowMsgWithCallerDepth(msg string, parent error, callerDepth int) {
	throwWithCallerDepth(New(msg, parent), callerDepth)
}

func TryCatch(try func(), catch func(err error), errs ...error) {
	defer func() {
		if recv := recover(); recv != nil {
			if e, ok := recv.(error); ok {
				for _, err := range errs {
					if errors.Is(e, err) {
						catch(e)
						return
					}
				}
			}
			panic(recv)
		}
	}()
	try()
}
