package main

import (
	"flag"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/startup"
	"github.com/firstsatoshi/website/common/task"
	"github.com/firstsatoshi/website/internal/config"
	"github.com/firstsatoshi/website/tasks/deposit"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "./etc/website-api.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	env := "pro"
	chainCfg := &chaincfg.MainNetParams
	apiHost := "https://mempool.space/api"
	if len(os.Getenv("BITEAGLE_TESTNET")) != 0 {
		chainCfg = &chaincfg.TestNet3Params
		apiHost = "https://mempool.space/testnet/api"
		env = "dev"
	}

	c.LogConf.ServiceName = "deposit-" + env
	logx.SetUp(c.LogConf)

	logx.Info("========Btc Deposit Task Start======")
	depositTask := deposit.NewBtcDepositTask(apiHost, &c, chainCfg)

	startup.TaskStartup([]task.Task{depositTask})
}
