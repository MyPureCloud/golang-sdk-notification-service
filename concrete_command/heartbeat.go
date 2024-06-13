package concrete_command

import "fmt"

type heartbeat struct {
	message string
}

func NewHeartbeat(message string) *heartbeat {
	return &heartbeat{
		message: message,
	}
}

func (hb *heartbeat) Execute() {
	fmt.Println("calling heartbeat Execute()")
	fmt.Println(hb.message)
}