

package main

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/

import "C"

func set_button_lamp(button int, floor int, value int){
	if (value) > 0 {
		C.io_set_bit(lamp_channel_matrix[floor][button])
	}
}

func main(){
	set_button_lamp(0,1,1)
}
