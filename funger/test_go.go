package main 
/*
#include "test.h"
*/
import "C"
import "fmt"

func main(){
	i := C.test_func()
	fmt.Println(i)

}

//kjøring: go install in folder.
