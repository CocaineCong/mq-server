package main

import (
	"mq-server/conf"
	"mq-server/services"
)

func main() {
	conf.Init()

	forever := make(chan bool)
	services.MQ2MySQL()
	<-forever
}
