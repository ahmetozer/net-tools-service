package functions

import "fmt"

func Recover(Where string) {
	if r := recover(); r != nil {
		fmt.Println("Recovered from ", Where, r)
	}
}
