package main
import "fmt"

func main() {
    myString := "Hello Golang"
    for i := 0; i < len(myString); i++ {
		fmt.Printf("%c ", myString[i])
	}
}