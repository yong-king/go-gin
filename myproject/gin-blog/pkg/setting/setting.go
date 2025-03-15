package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var Cfg *ini.File

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	PrefixUrl      string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath    string
	LogSaveName    string
	LogFileExt     string
	TimeFormat     string
	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string
}

var AppSetting = &App{}

type Serve struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServeSetting = &Serve{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host         string
	Password     string
	MaxIdle      int
	MaxActive    int
	IdleTimeoout time.Duration
}

var RedisSetting = &Redis{}

func Setup() {
	var err error
	Cfg, err = ini.Load("gin-blog/conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'gin-blog/conf/app.ini': %v", err)
	}

	//err = Cfg.Section("app").MapTo(AppSetting)
	//if err != nil {
	//	log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	//}
	//err = Cfg.Section("server").MapTo(ServeSetting)
	//if err != nil {
	//	log.Fatalf("Cfg.MapTo ServeSetting err: %v", err)
	//}
	//err = Cfg.Section("database").MapTo(DatabaseSetting)
	//if err != nil {
	//	log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	//}
	mapTo("app", AppSetting)
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	mapTo("database", DatabaseSetting)
	mapTo("server", ServeSetting)
	ServeSetting.ReadTimeout = ServeSetting.ReadTimeout * time.Second
	ServeSetting.WriteTimeout = ServeSetting.WriteTimeout * time.Second
	mapTo("redis", RedisSetting)
	RedisSetting.IdleTimeoout = RedisSetting.IdleTimeoout * time.Second

	//LoadBase()
	//LoadServer()
	//LoadApp()
}

func mapTo(section string, v interface{}) {
	err := Cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo section %s err: %v", section, err)
	}
}

//func LoadApp() {
//	sec, err := Cfg.GetSection("app")
//	if err != nil {
//		log.Fatalf("Fail to get section 'app': %v", err)
//	}
//	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
//	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!")
//}
//
//func LoadServer() {
//	sec, err := Cfg.GetSection("server")
//	if err != nil {
//		log.Fatalf("Fail to get section 'server': %v", err)
//	}
//	HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
//	// time.Duration 代表 时间间隔，本质是 int64，单位是 纳秒
//	// time.Duration(60) * time.Second = 60s
//	ReadTimeOut = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
//	WriteTimeOut = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
//}
//
//func LoadBase() {
//	RunMode = Cfg.Section("").Key("RunMode").MustString("debug")
//}
