package main

import (
	"flag"
	"github.com/CloudStriver/cloudmind-mq/internal/config"
	"github.com/CloudStriver/cloudmind-mq/internal/listen"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/config.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config
	logx.DisableStat()
	conf.MustLoad(*configFile, &c)
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range listen.Mqs(c) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
