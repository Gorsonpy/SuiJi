// Code generated by hertz generator.

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/XZ0730/runFzu/biz/dal"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"

	"github.com/XZ0730/runFzu/config"
	"github.com/XZ0730/runFzu/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzUtils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/hertz-contrib/gzip"
)

func main() {

	path := flag.String("config", "./config", "config path")
	klog.Info(*path)
	config.Init(*path, "config.yaml", "runnerFzu")
	dal.Init()
	r := server.New(
		server.WithHostPorts("0.0.0.0:8087"),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(1<<31),
	)

	r.Use(recovery.Recovery(recovery.WithRecoveryHandler(recoveryHandler)))

	// Gzip
	r.Use(gzip.Gzip(gzip.BestSpeed))

	// Sentinel 流量治理
	r.Use(hertzSentinel.SentinelServerMiddleware(
		hertzSentinel.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return "server_test"
		}),
		hertzSentinel.WithServerBlockFallback(func(ctx context.Context, c *app.RequestContext) {
			hlog.CtxInfof(ctx, "frequent requests have been rejected by the gateway. clientIP: %v\n", c.ClientIP())
			c.AbortWithStatusJSON(400, hertzUtils.H{
				"status_msg":  "too many request; the quota used up",
				"status_code": -1,
			})
		}),
	))
	register(r)
	r.Spin()
}

func recoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {

	hlog.CtxInfof(ctx, "[Recovery] InternalServiceError err=%v\n stack=%s\n", err, stack)
	c.JSON(consts.StatusInternalServerError, map[string]interface{}{
		"code":    errno.ServiceErrorCode,
		"message": fmt.Sprintf("[Recovery] err=%v\nstack=%s", err, stack),
	})
}
