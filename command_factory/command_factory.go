package command_factory

import (
	"fmt"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/command_interface"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/concrete_command"

)

type commandFactory struct {}

func NewCommandFactory() *commandFactory {
	return &commandFactory{}
}

func (cf *commandFactory) GetCommand(messageType string, message string) (command_interface.Command, error) {
	if messageType == "users.{id}.presence" {
		return concrete_command.NewUsersIdPresence(message), nil
	}
	if messageType == "routing.queues.{id}.users" {
		return concrete_command.NewRoutingQueuesIdUsers(message), nil
	}
	if messageType == "WebSocket Heartbeat" {
		return concrete_command.NewHeartbeat(message), nil
	}
	return nil, fmt.Errorf("%s", "command not found")
}
