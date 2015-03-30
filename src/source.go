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

const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
)

const(
	OFF			=  0
	ON 			=  1
)

type ButtonOrder struct{
	Floor 		int
	ButtonType 	int 
}

type ElevatorData struct {
	ID 				int
	Floor 			int
	Direction 		int
	OutsideOrders [N_FLOORS][2]int
	InsideOrders  [N_FLOORS]int
}

type Command struct{
	CommandType int
	Value int	
}

const(
	SET_BUTTON_LAMP				= 0
	SET_MOTOR_DIR				= 1
	SET_FLOOR_INDICATOR_LAMP	= 2
	SET_DOOR_OPEN_LAMP			= 3
)



