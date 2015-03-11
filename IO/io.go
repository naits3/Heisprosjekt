package IO

import "Heisprosjekt/driver"
import "runtime"
import"fmt"

var N_FLOOR int

//cmd is shorten with command 
//TODO: Fix this struct to be more generical when it comes to commands
type Cmd struct{
	cmdType int
	floor int
	value int
	buttonType int
}

//TODO implement a order type whitch is formated by the formate functions

const(
	SET_BUTTON_LAMP		= iota
	SET_MOTOR_DIR
	SET_FLOOR_INDICATOR_LAMP 
)

func InitializeIo(chCmdCtrl chan Cmd, chCtrl chan int){
	
	runtime.GOMAXPROCS(runtime.NumCPU())

	N_FLOOR 		= driver.GetN_FLOOR()
	chButtonOrder 	:= make(chan []int)
	chFloorSensor  	:= make(chan int)
	
	if err:=driver.Init(); err<0{
		fmt.Println("Could not initialize hardware")
		//TODO Handle error, exit. Try one more time
	}

	go ioHandler(chCmdCtrl,chCtrl,chButtonOrder,chFloorSensor)	
	go pollButtonOrders(chButtonOrder)
	go pollFLoorSensors(chFloorSensor)
}

func ioHandler(chCmdCtrl chan Cmd, chCtrl chan int, chButtonOrder chan []int, chFloorSensor chan int){
	for{
		select{
			case <- chButtonOrder:
				//format order
				//send order
				chCtrl <- 1

			case <- chFloorSensor:
				//format newFloor
				//send newfloor
				chCtrl <- 2

			case cmdCtrl := <- chCmdCtrl:
				doCmd(cmdCtrl)
		}
	}
}

func pollButtonOrders(chButtonOrder chan []int){
	//TODO: Test funksjonen og fiks problemet med slices!!!
	for{
		for floor := 0; floor < N_FLOOR; floor++ {
			
			if(floor > 0 && driver.GetButtonSignal(driver.BUTTON_CALL_UP,floor)==1){
				chButtonOrder <- []int{0,floor}
			}
			
			if(floor < N_FLOOR-2 && driver.GetButtonSignal(driver.BUTTON_CALL_DOWN,floor)==1){
				chButtonOrder <- []int{1,floor}
			}
			
			if(driver.GetButtonSignal(driver.BUTTON_COMMAND,floor)==1){
				chButtonOrder <- []int{2,floor}
			}
		}
	}
}

func pollFLoorSensors(chFloorSensor chan int){
	//TODO: Test funksjonen
	for{	
		if floor := driver.GetFloorSensor(); floor != -1{
		chFloorSensor <- floor
		}
	}
}

func doCmd(cmdCtrl Cmd){
	//TODO  Fiks enum problemet med typer.
	if(SET_BUTTON_LAMP == cmdCtrl.cmdType){
		driver.SetButtonLamp(1, cmdCtrl.floor, cmdCtrl.value)
	}

	if(SET_MOTOR_DIR == cmdCtrl.cmdType){
		driver.SetMotorDir(1)
	}

	if(SET_FLOOR_INDICATOR_LAMP == cmdCtrl.cmdType){
		driver.SetFloorIndicatorLamp(2)
	}
}

func formatOrder(order int) int{
	//Format the order
	return order
}

func formatNewFloor(newFloor int) int{
	//Format the newFLoor
	return newFloor
}




