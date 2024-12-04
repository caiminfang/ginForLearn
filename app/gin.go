package app

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hello/configs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const timeFormatter = "2006-01-02 15:04:05.999999"

// g *gin.Engine核心结构体，路由引擎
func StartHttpServer(router func(g *gin.Engine)) {
	conf := Conf.Server
	srv := newHttpServer(conf, newGin(router))
	go runServer(srv)
	//Logger.Infof(">>>>>> Http Server started on %s", srv.Addr)
	//Logger.Infof(">>>>>> current version: %s", version)
	gracefulShutdown(conf, srv)

}

func gracefulShutdown(conf configs.ServerConf, srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	//等待信号中断 关闭server
	quit := make(chan os.Signal)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("quit")

	<-quit
	//Logger.Infof(">>>>>> Shutdown WebServer ...")

	// HttpServer退出等待时间，默认等待10秒
	timeout := 10 * time.Second
	if conf.ShutdownTimeout > 0 {
		timeout = time.Duration(conf.ShutdownTimeout) * time.Second
	}
	// 设置退出等待超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// 关闭HttpServer
	if err := srv.Shutdown(ctx); err != nil {
		// os.Exit() deferred functions are not run.
		//Logger.Infof(">>>>>> WebServer  timeout shutdown:", err)
		return
	}
	//Logger.Infof(">>>>>> WebServer closed")
}

func runServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v\n", err)
	}
}

// 处理跨域请求,支持options访问
func newCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowCredentials: true,
	})
}

// newGin 创建一个自定义的GIN引擎
func newGin(router func(g *gin.Engine)) *gin.Engine {
	// 禁用控制台颜色,非必须设置
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	// 初始化GIN
	//创建新的 Gin 引擎：使用 gin.New() 创建一个新的 Gin 路由引擎实例 r。与 gin.Default() 不同，gin.New() 不会自动注册日志和恢复中间件。
	r := gin.New()
	//使用 newCors() 中间件来处理跨源资源共享 (CORS) 的请求，常用于允许来自不同源的请求。
	if Conf.Server.OpenCors {
		r.Use(newCors())
	}
	//注册中间件
	r.Use(ignoreIndexAndFavicon(), ginZap(zapLogger, timeFormatter, false), ginzap.RecoveryWithZap(zapLogger, true))
	// 限制单个请求大小（4M）
	r.Use(limits.RequestSizeLimiter(4 << 20))
	// 注册路由
	router(r)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": "404", "data": "", "msg": "请求的API不存在111"})
	})
	return r
}

// ginZap 日志记录，服务处理时间超过200ms才需要输出日志
func ginZap(logger *zap.Logger, timFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else if latency.Milliseconds() >= 200 {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timFormat)),
				zap.Duration("latency", latency),
			)
		}
	}
}

// 忽略index和favicon.ico请求
func ignoreIndexAndFavicon() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/favicon.ico" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if path == "/" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func newHttpServer(conf configs.ServerConf, h http.Handler) *http.Server {
	fmt.Println("'=========================='")
	fmt.Println(conf.Port)
	//addr := fmt.Sprintf(":#{conf.Port}")
	addr := fmt.Sprintf(":%d", conf.Port) // 假设 conf.Port 是整数
	//addr := conf.Port
	fmt.Println(addr)
	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}
