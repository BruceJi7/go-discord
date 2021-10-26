package helpers

import "fmt"

func RaiseError(err error) {
	if err != nil {
		panic(err)
	}
}
func RaiseOrPrint(err error, msg string) {
	if err != nil {
		panic(err)
	} else {
		fmt.Println(msg)
	}
}
