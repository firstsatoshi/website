package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/internal/handler"
	"github.com/firstsatoshi/website/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/website-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	env := "pro"
	if len(os.Getenv("BITEAGLE_TESTNET")) != 0 {
		env = "dev"
	}

	c.LogConf.ServiceName = "website-" + env
	logx.SetUp(c.LogConf)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
