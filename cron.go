package main

import (
	"github.com/linbihuan/gocron/spiders/calendar"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 */1 * * * *", func() {
		calendar.Crawl()
	})

	//启动计划任务
	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务
	defer c.Stop()

	select {}
}
