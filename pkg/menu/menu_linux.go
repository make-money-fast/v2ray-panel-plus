//go:build linux

package menu

var (
	UIAddress       = ""
	ServerUIAddress = ""
	CloseChan       = make(chan struct{})
)

func Init() {

}

func shutdown() {

}

func Shutdown() {
}
