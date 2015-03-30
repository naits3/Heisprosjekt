package main

import "Heisprosjekt/driver"
import "time"
import "fmt"
import "Heisprosjekt/src"

func main() {
	
	n := driver.Init() 
	if n < 0{
		fmt.Println("Could not initialize hardware")
	}

	driver.SetMotorDir(src.DIR_UP)
	time.Sleep(1*time.Second)
	driver.SetMotorDir(src.DIR_STOP)
}