package helpers

import (
	"fmt"
	"testing"
)

func TestGetMyIP(t *testing.T) {
	fmt.Println("my ip address is:", GetMyIP())
}
