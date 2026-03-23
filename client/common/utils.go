package common

import (
	"fmt"
	"os"
)

func ListenForSigTerm(channel chan os.Signal) bool {
	select {
	case sig := <-channel:
		fmt.Println("Received signal", sig)
		return true
	default:
		return false
	}
}
