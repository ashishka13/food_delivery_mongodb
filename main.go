package main

import (
	"food-delivery/controller"
	"food-delivery/worker"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go worker.MainWorker() // very useful for interval polling

	controller.MyController()
	// select {} // this will cause the program to run forever
}
