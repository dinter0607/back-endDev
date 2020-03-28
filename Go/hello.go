package main

import "fmt"

func myFunction() {
	var a int
	a = 1
	for i := 0; i < 10 ; i++ {
		a = a+i
		if a % 2 == 0 {
			fmt.Println(a, "Chia hết cho 2")
		} else {
			fmt.Println(a, "không chia hết cho 2")
		}
	}  
}

func main() {
    myFunction()
    fmt.Println("Hello, World!")
}