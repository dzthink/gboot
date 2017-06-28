package main

import (
	"fmt"
	"boot"
)

func main() {
	app := &Test{}
	booter := boot.NewBooter(app)
	booter.Run()
}

type Test struct {

}

func(t *Test) Name() string {
	return "test"
}

func(t *Test) Start(conf boot.ConfigInterface) (error) {
	fmt.Println("start...")
	return nil
}

func(t *Test) Reload(conf boot.ConfigInterface) (error) {
	fmt.Println("reload...")
	return nil
}

func(t *Test) Stop() error {
	fmt.Println("stop...")
	return nil
}
