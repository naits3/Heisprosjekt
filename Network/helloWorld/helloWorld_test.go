package helloWorld1


/*
Skal ligge i samme mappe som pakka
*/

import ("testing")

func TestReturnHW(t *testing.T){
	var str string = ReturnHW()
	if(str != "Hello World" ){
		t.Error("Expected Hello World, got"+" "+ str)
	}
}

