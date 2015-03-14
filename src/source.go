package src

const n_floors = 4

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
	floor 		int
	ButtonType 	int 
}

func GetNFloors() int{
	return n_floors
}

type ElevatorData struct {
	floor 			int
	direction 		int
	outsideOrders [n_floors][2]int
	insideOrders  [n_floors]int
}