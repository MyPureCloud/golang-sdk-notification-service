package concrete_command

import "fmt"

type usersIdPresence struct {
	message string
}

func NewUsersIdPresence(message string) *usersIdPresence {
	return &usersIdPresence{
		message: message,
	}
}

func (uip *usersIdPresence) Execute() {
	fmt.Println("calling users.{id}.presence Execute()")
	fmt.Println(uip.message)
}
