package main

import (
	"fmt"

	"github.com/groovypotato/gator/internal/config"
)
func main() {
	currConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = currConfig.SetUser("groovypotato")
	if err != nil {
		fmt.Println(err)
		return
	}
	newConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(newConfig)
	}
}