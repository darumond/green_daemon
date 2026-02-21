package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "client" {
		fmt.Println("running client")
		clientMain()
	} else {
		fmt.Println("running server")
		serverMain()
	}
}
