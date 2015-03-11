package main

import "Heisprosjekt/driver"
import "time"
import "fmt"

func main() {
	
	n := driver.Init() 
	if n < 0{
		fmt.Println("Could not initialize hardware")
	}

	driver.SetMotorDir(1)
	time.Sleep(1*time.Second)
	driver.SetMotorDir(0)

}