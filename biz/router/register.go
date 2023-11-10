// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	base "github.com/XZ0730/runFzu/biz/router/base"
	goal "github.com/XZ0730/runFzu/biz/router/goal"
	multiledger "github.com/XZ0730/runFzu/biz/router/multiledger"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	multiledger.Register(r)

	goal.Register(r)

	base.Register(r)

}
