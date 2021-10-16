/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 11:17
 */
package api

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gochat/api/router"
	"gochat/api/rpc"
	"gochat/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Chat struct {
}

func New() *Chat {
	return &Chat{}
}

// Run api server,Also, you can use gin,echo ... framework wrap
func (c *Chat) Run() {
	//init rpc client
	rpc.InitLogicRpcClient()

	r := router.Register()
	runMode := config.GetGinRunMode()
	logrus.Info("server start , now run mode is ", runMode)
	gin.SetMode(runMode)
	apiConfig := config.Conf.Api
	port := apiConfig.ApiBase.ListenPort
	// @tip 这里无需再调用 flag.Parse() 了吧
	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()
	// @tip 这里就会阻塞住
	// if have two quit signal , this signal will priority capture ,also can graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")

	// @tip 这会让主线程也退出吧; os.Exit() 会让 defer 不被执行
	os.Exit(0)
}
