package cli

import (
	"fmt"
	"os"
)

func (a *App) Present(obj interface{}) error {
	return a.Output.Write(os.Stdout, obj)
}

func (a *App) Fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
