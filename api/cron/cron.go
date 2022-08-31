package cron

import (
	"github.com/robfig/cron"

	"zendea/service"
	"zendea/util"
	"zendea/util/log"
)

func Setup() {
	if !util.IsProd() {
		log.Info("Not in a production enviroment!")
		return
	}

	log.Info("Cron setup")
	// 开启定时任务
	startSchedule()
}

func startSchedule() {
	c := cron.New()

	// Generate Staticstic Data
	addCronFunc(c, "@every 10m", func() {
		service.StatisticService.GenerateData()
	})

	// Generate RSS
	addCronFunc(c, "@every 30m", func() {
		service.ArticleService.GenerateRss()
		service.TopicService.GenerateRss()
	})

	// Generate sitemap
	addCronFunc(c, "@every 45m", func() {
		service.SitemapService.GenerateToday()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		log.Error(err.Error())
	}
}
