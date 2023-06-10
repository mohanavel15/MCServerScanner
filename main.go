package main

import (
	"MCServerScanner/pkg/mcstatus"
	"fmt"
	"time"
)

func main() {
	timeout := time.Second * 3
	data, err := mcstatus.Lookup("0.0.0.0", 25565, timeout)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)
}
