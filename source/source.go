package source

const n_floors = 4

const(
	BUTTON_UP 		= 0
	BUTTON_DOWN		= 1
	BUTTON_INSIDE	= 2
)
	
const(
	DIRN_DOWN 	 = -1
	DIRN_STOP    = 0
	DIRN_UP      = 1
)

type Order struct{
	floor 		int
	ButtonType 	int 
}


func GetNFloors() int{
	return n_floors
}
