package app

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  errCode,
		"data": data,
	})
	return
}
