package main

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	//api文档自动生成开始
	yaag.Init(&yaag.Config{On: true, DocTitle: "go-bbs", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Stage": ""}})
	app.Use(irisyaag.New())
	//api文档自动生成结束

	//集成prometheus监控开始,访问/metrics
	m := prometheusMiddleware.New("go-bbs", 300, 1200, 5000)
	app.Use(m.ServeHTTP)
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		//错误代码处理程序不与其他路由共享相同的中间件，所以单独执行错误
		m.ServeHTTP(ctx)
		_, _ = ctx.Writef("Not Found")
	})

	app.Handle("GET", "/", func(ctx iris.Context) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		_, _ = ctx.Writef("Slept for %d milliseconds", sleep)
	})
	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.WriteString("pong")
	})
	app.Get("/hello", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{"message": "hello"})
	})
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	_ = app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
