package main

import (
	"flag"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/startup"
	"github.com/firstsatoshi/website/common/task"
	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/tasks/inscribe"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "./etc/website-api.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	chainCfg := &chaincfg.MainNetParams
	apiHost := "https://mempool.space/api"
	env := "pro"
	if len(os.Getenv("BITEAGLE_TESTNET")) != 0 {
		chainCfg = &chaincfg.TestNet3Params
		apiHost = "https://mempool.space/testnet/api"
		env = "dev"
	}

	c.LogConf.ServiceName = "inscribe-" + env
	logx.SetUp(c.LogConf)

	logx.Info("========Btc Inscribe Task Start======")
	inscribeTask := inscribe.NewBtcInscribeTask(apiHost, &c, chainCfg)

	startup.TaskStartup([]task.Task{inscribeTask})
}
