package ecoscript

import (
	"fmt"
)

func Guard(err error, msgAndArgs ...interface{}) {
	var msg string
	if err != nil {
		if len(msgAndArgs) == 0 || msgAndArgs == nil {
			msg = err.Error()
		} else {
			if len(msgAndArgs) == 1 {
				msg = msgAndArgs[0].(string)
			}
			if len(msgAndArgs) > 1 {
				msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
			}
			msg = fmt.Sprintf("%s: %s", msg, err)
		}
		panic(msg)
	}
}
