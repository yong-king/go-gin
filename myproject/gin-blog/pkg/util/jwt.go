package utills

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/youngking/gin-blog/pkg/setting"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

// Claims结构体
type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 生成JWT
func GenerateToken(username string, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err

}

//func GenerateToken(username, password string) (string, error) {
//	nowTime := time.Now()
//	expireTime := nowTime.Add(3 * time.Hour)
//	claims := Claims{
//		username,
//		password,
//		jwt.StandardClaims{
//			ExpiresAt: expireTime.Unix(),
//			Issuer:    "gin-blog",
//		},
//	}
//	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	token, err := tokenClaims.SignedString(jwtSecret)
//	return token, err
//}

// 解析JWT
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
