package main

import (
	"flag"
	"fmt"
	"git.zc0901.com/go/god/api/httpx"
	"github.com/xqk/cymzjs-api/pkg"

	"github.com/xqk/cymzjs-api/api/cymzjs/internal/config"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/handler"
	"github.com/xqk/cymzjs-api/api/cymzjs/internal/svc"

	"git.zc0901.com/go/god/api"
	"git.zc0901.com/go/god/lib/conf"
)

var configFile = flag.String("f", "etc/cymzjs-api.yaml", "配置文件")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := api.MustNewServer(c.ServerConf)
	defer server.Stop()

	httpx.SetErrorHandler(pkg.ApiErrHandler)
	httpx.SetOkJsonHandler(pkg.ApiOKHandler)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
