package sys

import (
	"fmt"
	"os"

	"github.com/haochend413/muninx/config"
)

func LogError(err error) {
	// this needs to be controlled by the overall config.
	err_file := config.ErrorFilePathDefualt()
	f, openErr := os.OpenFile(err_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if openErr != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, err.Error())
}
