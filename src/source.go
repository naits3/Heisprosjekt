package src

const N_FLOORS = 4

const(
	BUTTON_UP 		= 0
	BUTTON_DOWN		= 1
	BUTTON_INSIDE	= 2
)
	
const(
	DIR_DOWN 	= -1
	DIR_STOP    =  0
	DIR_UP      =  1
)

type ButtonOrder struct{
	Floor 		int
	ButtonType 	int 
}

type ElevatorData struct {
	Floor 			int
	Direction 		int
	OutsideOrders [N_FLOORS][2]int
	InsideOrders  [N_FLOORS]int
}

