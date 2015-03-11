package main

import "Heisprosjekt/driver"
import "time"
import "fmt"
import "runtime"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//done chanel
	done := make(chan bool)

	n := driver.Init() 
	if n < 0{
		fmt.Println("Could not initialize hardware")
	}



	// go runElevator(done)
 //    go floorsensor()
 //    go buttonSignalUp(done)
	// go buttonSignalDown(done)
	//go buttonSignalInside(done)
	go setButtonLamp()
	// go setFloorIndicatorLamp()

	<-done
	<-done
	<-done
	<-done
}

func runElevator(done chan bool){
		driver.SetMotorDir(1)
		time.Sleep(1*time.Second)
		driver.SetMotorDir(0)
		done <- true
}

func floorsensor(){
	for{
		if floor:=driver.GetFloorSensor(); floor > -1{
			fmt.Println("floor sensor responed\n")
			break
		}
	}
}

func buttonSignalUp(done chan bool){
	floor_nr := 0
	for{
		if(driver.GetButtonSignal(0,floor_nr)==1){
			fmt.Println("Button is pushed up on floor ",floor_nr)
			floor_nr = floor_nr + 1
			if(floor_nr==3){
				break
			}
		}
	}

	done <- true
}

func buttonSignalDown(done chan bool){
	floor_nr := 1
	for{
		if(driver.GetButtonSignal(1,floor_nr)==1){
			fmt.Println("Button is pushed down on floor ",floor_nr)
			floor_nr = floor_nr + 1
			if(floor_nr==4){
				break
			}
		}
	}

	done <- true
}

func buttonSignalInside(done chan bool){
	floor_nr := 0
	for{
		if(driver.GetButtonSignal(2,floor_nr)==1){
			fmt.Println("Button is pushed inside for floor ",floor_nr)
			floor_nr = floor_nr + 1
			if(floor_nr==4){
				break
			}
		}
	}

	done <- true
}

func setButtonLamp(){
	driver.SetButtonLamp(0,0,1)
	driver.SetButtonLamp(0,1,1)
	driver.SetButtonLamp(0,2,1)
	driver.SetButtonLamp(1,1,1)
	driver.SetButtonLamp(1,2,1)
	driver.SetButtonLamp(1,3,1)
	driver.SetButtonLamp(2,0,1)
	driver.SetButtonLamp(2,1,1)
	driver.SetButtonLamp(2,2,1)
	driver.SetButtonLamp(2,3,1)
}

func setFloorIndicatorLamp(){
	driver.SetFloorIndicatorLamp(0)
	time.Sleep(time.Second)
	driver.SetFloorIndicatorLamp(1)
	time.Sleep(time.Second)
	driver.SetFloorIndicatorLamp(2)
	time.Sleep(time.Second)
	driver.SetFloorIndicatorLamp(3)
}