//package main
//
//import (
//	"fmt"
//	"github.com/fvbock/endless"
//	_ "github.com/youngking/gin-blog/docs"
//	"github.com/youngking/gin-blog/models"
//	"github.com/youngking/gin-blog/pkg/gredis"
//	"github.com/youngking/gin-blog/pkg/logging"
//	"github.com/youngking/gin-blog/pkg/setting"
//	"github.com/youngking/gin-blog/routers"
//	"log"
//	"syscall"
//)
//
//func main() {
//	setting.Setup()
//	models.Setup()
//	logging.Setup()
//	gredis.SetUp()
//	//endles server
//	//设置服务器默认参数
//	endless.DefaultReadTimeOut = setting.ServeSetting.ReadTimeout
//	endless.DefaultWriteTimeOut = setting.ServeSetting.WriteTimeout
//	endless.DefaultMaxHeaderBytes = 1 << 20
//	endPoint := fmt.Sprintf(":%d", setting.ServeSetting.HttpPort)
//
//	// 创建服务器
//	server := endless.NewServer(endPoint, routers.InitRouter())
//	// 获取进程id
//	server.BeforeBegin = func(add string) {
//		log.Printf("Actual pid is %d", syscall.Getpid())
//	}
//	// 	启动服务
//	err := server.ListenAndServe()
//	if err != nil {
//		log.Printf("Server err: %v", err)
//	}
//
//}

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
	"sync"
	"syscall"
)

func startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 设置端口号
	endPoint := fmt.Sprintf(":%d", port)

	// 创建服务器
	server := endless.NewServer(endPoint, routers.InitRouter())

	// 获取进程ID
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	// 启动服务器
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err on port %d: %v", port, err)
	}
}

func main() {
	// 设置初始化
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.SetUp()

	// 设置服务器默认参数
	endless.DefaultReadTimeOut = setting.ServeSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServeSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	var wg sync.WaitGroup

	// 启动两个服务器，分别监听8001和8002端口
	ports := []int{8001, 8002}
	for _, port := range ports {
		wg.Add(1)
		go startServer(port, &wg)
	}

	// 等待所有服务器完成
	wg.Wait()
}
