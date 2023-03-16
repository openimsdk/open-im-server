package specialerror

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var handlers []func(err error) errs.CodeError

func AddErrHandler(h func(err error) errs.CodeError) {
	if h == nil {
		panic("nil handler")
	}
	handlers = append(handlers, h)
}

func AddReplace(target error, codeErr errs.CodeError) {
	AddErrHandler(func(err error) errs.CodeError {
		if err == target {
			return codeErr
		}
		return nil
	})
}

func ErrCode(err error) errs.CodeError {
	if codeErr, ok := err.(errs.CodeError); ok {
		return codeErr
	}
	for i := 0; i < len(handlers); i++ {
		if codeErr := handlers[i](err); codeErr != nil {
			return codeErr
		}
	}
	return nil
}
