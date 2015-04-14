package main

import "Heisprosjekt/controller"

func main() {
	stop := make(chan int)
	controller.InitController()
	<- stop
}