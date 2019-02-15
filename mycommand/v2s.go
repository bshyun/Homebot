package mycommand

import (
	"fmt"
)

func (cp *commandproc) V2S() (string, error) {
	fmt.Println("v2s")

	cp.Command = cp.Command.(string)

	return cp.Command.(string), nil
}
