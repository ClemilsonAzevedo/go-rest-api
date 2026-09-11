package main

import (
	"fmt"

	"github.com/clemilsonazevedo/9-apis/configs"
)

func main() {
	cfg, err := configs.LoadConfigs(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg)
}
