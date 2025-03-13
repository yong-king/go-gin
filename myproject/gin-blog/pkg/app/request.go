package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/youngking/gin-blog/pkg/logging"
)

func MarkErrors(errors []*validation.Error) {
	for _, e := range errors {
		logging.Info(e.Key, e.Message)
	}
	return
}
