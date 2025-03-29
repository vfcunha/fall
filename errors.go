package fall

import (
	"fmt"
	"log"
)

func PanicMsgIfError(err error, message string) {
	if err != nil {
		if message != "" {
			panic(fmt.Errorf("%s %w", message, err))
		}
	}
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogFatalIfError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func LogIfError(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}

// type ErrorInterceptor[R any] =

func InterceptErrorp[R any](interceptor func() (R, error), errMessage string) (R, error) {
	algo, err := interceptor()
	if err != nil {
		return algo, fmt.Errorf("%s %w", errMessage, err)
	}
	return algo, nil
}
