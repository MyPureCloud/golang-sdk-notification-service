package concrete_command

import "fmt"

type routingQueuesIdUsers struct {
	message string
}

func NewRoutingQueuesIdUsers(message string) *routingQueuesIdUsers {
	return &routingQueuesIdUsers{
		message: message,
	}
}

func (rqiu *routingQueuesIdUsers) Execute() {
	fmt.Println("calling routing.queues.{id}.users Execute()")
	fmt.Println(rqiu.message)
}
