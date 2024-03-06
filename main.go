package main

import (
	"aks-demo/router"
	"fmt"
)

func main() {
	err := router.Router.Run(fmt.Sprintf(":%d", 8080))
	if err != nil {
		panic(err)
	}
}
