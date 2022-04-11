package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("Exit!")
	if ok := initLogger(); !ok {
		return
	}
	if ok := initConfig(); !ok {
		return
	}
	if ok := initDB(); !ok {
		return
	}
	initServer()
}
