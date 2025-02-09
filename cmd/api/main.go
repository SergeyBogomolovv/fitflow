package main

import (
	"fmt"

	"github.com/SergeyBogomolovv/fitflow/config"
)

func main() {
	conf := config.MustNewConfig("./config/config.yml")
	fmt.Printf("%+v\n", conf)
}
