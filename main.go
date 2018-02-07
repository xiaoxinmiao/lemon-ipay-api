package main

import (
	"flag"
	"fmt"
	"lemon-ipay-api/alipay"
	"lemon-ipay-api/core"
	"lemon-ipay-api/ipay"
	"lemon-ipay-api/model"
	"lemon-ipay-api/wechat"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	appEnv      = flag.String("APP_ENV", os.Getenv("APP_ENV"), "APP_ENV")
	connEnv     = flag.String("IPAY_CONN", os.Getenv("IPAY_CONN"), "IPAY_CONN")
	bmappingUrl = flag.String("BMAPPING_URL", os.Getenv("BMAPPING_URL"), "BMAPPING_URL")
	hostUrl     = flag.String("IPAY_HOST", os.Getenv("IPAY_HOST"), "IPAY_HOST")
	jwtEnv      = flag.String("JWT_SECRET", os.Getenv("JWT_SECRET"), "JWT_SECRET")
)

func init() {
	flag.Parse()
	envParam := &core.EnvParamDto{
		AppEnv:      *appEnv,
		ConnEnv:     *connEnv,
		BmappingUrl: *bmappingUrl,
		HostUrl:     *hostUrl,
	}
	core.InitEnv(envParam)
	model.Db = InitDB("mysql", envParam.ConnEnv)
	model.Db.Sync(new(model.WxAccount), new(model.NotifyWechat), new(model.AlAccount), new(model.NotifyAlipay))
}

func InitDB(dialect, conn string) (newDb *xorm.Engine) {
	newDb, err := xorm.NewEngine(dialect, conn)
	if err != nil {
		panic(err)
	}
	return
}

func main() {

	e := echo.New()
	e.Use(middleware.CORS())
	RegisterApi(e)
	e.Start(":5000")
}

func RegisterApi(e *echo.Echo) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "lemon epay")
	})
	e.GET("/ping", func(c echo.Context) error {
		fmt.Println("pong")
		return c.String(http.StatusOK, "pong")
	})

	v3 := e.Group("/v3")
	v3.GET("/record/:Id", ipay.GetRecord)

	wx := v3.Group("/wx")
	wx.POST("/pay", wechat.Pay)
	wx.POST("/query", wechat.Query)
	wx.POST("/reverse", wechat.Reverse)
	wx.POST("/refund", wechat.Refund)
	wx.POST("/prepay", wechat.Prepay)
	wx.POST("/notify", wechat.Notify)
	wx.GET("/prepayeasy", wechat.PrepayEasy)
	wx.GET("/prepayopenid", wechat.PrepayOpenId)

	al := v3.Group("/al")
	al.POST("/pay", alipay.Pay)
	al.POST("/query", alipay.Query)
	al.POST("/reverse", alipay.Reverse)
	al.POST("/refund", alipay.Refund)
	al.POST("/prepay", alipay.Prepay)
	al.POST("/notify", alipay.Notify)

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(*jwtEnv),
		Skipper: func(c echo.Context) bool {
			ignore := []string{
				"/ping",
				"/v3/wx",
				"/v3/al",
			}
			for _, i := range ignore {
				if strings.HasPrefix(c.Request().URL.Path, i) {
					return true
				}
			}

			return false
		},
	}))

	wxEland := v3.Group("/jwt/wx")
	wxEland.POST("/pay", wechat.Pay)
	wxEland.POST("/query", wechat.Query)
	wxEland.POST("/reverse", wechat.Reverse)
	wxEland.POST("/refund", wechat.Refund)
	wxEland.POST("/prepay", wechat.Prepay)
	wxEland.POST("/notify", wechat.Notify)
	wxEland.GET("/prepayeasy", wechat.PrepayEasy)
	wxEland.GET("/prepayopenid", wechat.PrepayOpenId)

	alEland := v3.Group("/jwt/al")
	alEland.POST("/pay", alipay.Pay)
	alEland.POST("/query", alipay.Query)
	alEland.POST("/reverse", alipay.Reverse)
	alEland.POST("/refund", alipay.Refund)
	alEland.POST("/prepay", alipay.Prepay)
	alEland.POST("/notify", alipay.Notify)

}
