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
	logx.Info("========Btc Inscribe Task Start======")

	chainCfg := &chaincfg.MainNetParams
	apiHost := ""
	if len(os.Getenv("BITEAGLE_TESTNET")) != 0 {
		chainCfg = &chaincfg.SigNetParams
		apiHost = "TODO" // TODO
	}

	inscribeTask := inscribe.NewBtcInscribeTask(apiHost, &c, chainCfg)

	startup.TaskStartup([]task.Task{inscribeTask})
}
