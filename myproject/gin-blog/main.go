package main

import (
	"fmt"
	"github.com/fvbock/endless"
	_ "github.com/youngking/gin-blog/docs"
	"github.com/youngking/gin-blog/models"
	"github.com/youngking/gin-blog/pkg/gredis"
	"github.com/youngking/gin-blog/pkg/logging"
	"github.com/youngking/gin-blog/pkg/setting"
	"github.com/youngking/gin-blog/routers"
	"log"
	"syscall"
)

func main() {
	//// 创建默认路由引擎
	//router := gin.Default()
	//
	//// 定义test API
	//router.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	//router := routers.InitRouter()
	//
	//// 创建http请求
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HttpPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeOut,
	//	WriteTimeout:   setting.WriteTimeOut,
	//	MaxHeaderBytes: 1 << 20, // 设置 HTTP 请求头的最大字节数，1M
	//}
	//s.ListenAndServe()
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.SetUp()
	//endles server
	//设置服务器默认参数
	endless.DefaultReadTimeOut = setting.ServeSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServeSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServeSetting.HttpPort)

	// 创建服务器
	server := endless.NewServer(endPoint, routers.InitRouter())
	// 获取进程id
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	// 	启动服务
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

	//router := routers.InitRouter()
	//
	//// 创建http请求
	//s := &http.Server{
	//	Addr:           fmt.Sprintf(":%d", setting.HttpPort),
	//	Handler:        router,
	//	ReadTimeout:    setting.ReadTimeOut,
	//	WriteTimeout:   setting.WriteTimeOut,
	//	MaxHeaderBytes: 1 << 20, // 设置 HTTP 请求头的最大字节数，1M
	//}
	//
	//go func() {
	//	if err := s.ListenAndServe(); err != nil {
	//		log.Printf("Listen err %s\n:", err)
	//	}
	//}()
	//
	//quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt)
	//<-quit
	//
	//log.Println("Shutdown Server ...")
	//
	//context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := s.Shutdown(context); err != nil {
	//	log.Fatal("Server Shutdown:", err)
	//}
	//log.Println("Server exiting")
}
