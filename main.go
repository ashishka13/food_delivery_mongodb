package main

import (
	"foodDelivery/controller"
	"foodDelivery/processor"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go processor.MainProcessor() // very useful for interval polling

	controller.MyController()
	// select {} // this will cause the program to run forever
}
