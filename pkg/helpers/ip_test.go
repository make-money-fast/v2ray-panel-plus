package helpers

import (
	"fmt"
	"testing"
)

func TestGetPublicIP(t *testing.T) {
	fmt.Println("my ip address is:", GetPublicIP())
}

func TestGetInternalIP(t *testing.T) {
	fmt.Println("my ip address is:", GetInternalIP())
}
