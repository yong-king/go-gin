package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/youngking/gin-blog/pkg/app"
	"github.com/youngking/gin-blog/pkg/e"
	utills "github.com/youngking/gin-blog/pkg/util"
	"github.com/youngking/gin-blog/service/auth_service"
	"net/http"
)

type auth struct {
	UserName string `valid:"Required;MaxSize(50)" json:"user_name"`
	Password string `valid:"Required;MaxSize(50)" json:"password"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}

	userName := c.Query("user_name")
	password := c.Query("password")

	a := auth{UserName: userName, Password: password}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: userName, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := utills.GenerateToken(userName, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
