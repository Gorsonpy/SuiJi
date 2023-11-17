// Code generated by hertz generator.

package asr

import (
	"context"

	"github.com/XZ0730/runFzu/biz/model/base"
	"github.com/XZ0730/runFzu/biz/pack"
	"github.com/XZ0730/runFzu/pkg/errno"
	"github.com/XZ0730/runFzu/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		func(c context.Context, ctx *app.RequestContext) {
			token_byte := ctx.GetHeader("token")
			if _, err := utils.CheckToken(string(token_byte)); err != nil {
				resp := base.NewBaseResponse()
				pack.PackBase(resp, errno.AuthorizationFailedErrCode, errno.AuthorizationFailedError.ErrorMsg)
				ctx.JSON(consts.StatusOK, resp)
				ctx.Abort()
				return
			}
			ctx.Next(c)
		},
	}
}
func _apiMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _asrtotextMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _fileuploadMw() []app.HandlerFunc {
	// your code...
	return nil
}
