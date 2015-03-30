package src

const N_FLOORS = 4


const (
	ORDER = 1
	EMPTY = 0
	DELETE_ORDER = -1
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
	SetValue int
	Floor int
	ButtonType int
}

const(
	SET_BUTTON_LAMP				= 0
	SET_MOTOR_DIR				= 1
	SET_FLOOR_INDICATOR_LAMP	= 2
	SET_DOOR_OPEN_LAMP			= 3
)

//Values
const(
	DIR_DOWN 	= -1
	DIR_STOP    =  0
	DIR_UP      =  1
)

const(
	OFF			=  0
	ON 			=  1
)

//Floors
const(
	FLOOR_NONE	= -1
	FLOOR_1		=  0
	FLOOR_2		=  1
	FLOOR_3		=  2
	FLOOR_4		=  3
)

//ButtonType
const(
	BUTTON_NONE		=-1
	BUTTON_UP 		= 0
	BUTTON_DOWN		= 1
	BUTTON_INSIDE	= 2
)


