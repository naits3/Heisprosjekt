package helloWorld1

/* 
Hvis funksjonen skal testes må den være en pakke som vil si:

- Ligger under ~Go\src\Xxxx 
- Kommandoen go install må kjøres for å compilere opp filene til
	bibliotek.

- link packages: http://www.golang-book.com/11/index.htm
- link testing:  http://www.golang-book.com/12/index.htm

Remember capital letters if a function on the outside is going
to see the function. It means "public" in the golang.
*/



func ReturnHW() string{
	return("Hello World")
}