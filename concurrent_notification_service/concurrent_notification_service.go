package concurrent_notification_service

import (
	"fmt"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/command_factory"
	"github.com/gorilla/websocket"
	"github.com/tidwall/pretty"
	"strings"
	//"github.com/MyPureCloud/golang-sdk-notification-service.git/command_factory"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/command_interface"
	"github.com/mypurecloud/platform-client-sdk-go/v56/platformclientv2"
	//"github.com/tidwall/pretty"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	//"strings"
)

type concurrentNotificationService struct {
	topics []string
}

func NewConcurrentNotificationService(topicList []string) *concurrentNotificationService {
	return &concurrentNotificationService{
		topics: topicList,
	}
}

func (cns *concurrentNotificationService) RunConcurrentNotificationService() {
	fmt.Println("=== running concurrent notification service ===")
	config := platformclientv2.GetDefaultConfiguration()
	cns.Authenticate(config)
	notificationsApi := platformclientv2.NewNotificationsApi()
	channelId, connectUri := cns.CreateChannel(notificationsApi)
	cns.SubscribeToTopic(notificationsApi, channelId)
	cns.Listen(connectUri)
}

func (cns *concurrentNotificationService) Authenticate(config *platformclientv2.Configuration) {
	fmt.Println("authenticating")
	err := config.AuthorizeClientCredentials(os.Getenv("GENESYSCLOUD_OAUTHCLIENT_ID"), os.Getenv("GENESYSCLOUD_OAUTHCLIENT_SECRET"))
	if err != nil {
		log.Fatalf("error authenticating: %v", err)
	}
}

func (cns *concurrentNotificationService) CreateChannel(notificationsApi *platformclientv2.NotificationsApi) (string, string) {
	fmt.Println("creating the notifications channel")
	channel, _, err := notificationsApi.PostNotificationsChannels()
	if err != nil {
		log.Fatalf("error creating notifications channel: %v", err)
	}
	return *channel.Id, *channel.ConnectUri
}

func (cns *concurrentNotificationService) SubscribeToTopic(notificationsApi *platformclientv2.NotificationsApi, channelId string) {
	var reqBody []platformclientv2.Channeltopic
	for i, topic := range cns.topics {
		fmt.Println("subscribing to topic: " + topic)
		reqBody = append(reqBody, platformclientv2.Channeltopic{Id: &cns.topics[i], SelfUri: nil})
	}
	_, _, err := notificationsApi.PostNotificationsChannelSubscriptions(channelId, reqBody)
	if err != nil {
		log.Fatalf("error subscribing to topics: %v", err)
	}
}

func (cns *concurrentNotificationService) Listen(connectUri string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: connectUri[0:3], Host: connectUri[6:31], Path: connectUri[32:]}

	connection, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("handshake failed with status %v: error %v", resp.StatusCode, err)
	}
	fmt.Println("connected to server")

	// close the connection when function returns
	defer connection.Close()

	commandsChan := make(chan command_interface.Command, 10)

	defer close(commandsChan)

	// spin off 10 go routines and listen for commands on the commands channel
	for i := 0; i < 10; i++ {
		go func(commandsChan chan command_interface.Command) {
			for command := range commandsChan  {
				command.Execute()
			}
		}(commandsChan)
	}

	// process incoming messages with factory and command pattern and send command on channel
	go readMessages(connection, commandsChan)

	for  {
		select {
		case <-interrupt:
			// send close message to the server
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("error closing connection: %v", err)
				return
			}
		}
	}

}

func readMessages(connection *websocket.Conn, commandsChan chan command_interface.Command) {
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Fatalf("error reading message: %v", err)
		}
		// format message
		prettyMsg := fmt.Sprintf("%s", pretty.Pretty(message))
		msg := strings.TrimSuffix(prettyMsg, "\n")
		messageType, err := getMessageType(msg)
		if err != nil {
			log.Fatalf("error getting message type: %v", err)
		}
		commandFactory := command_factory.NewCommandFactory()
		command, err := commandFactory.GetCommand(messageType, msg)
		if err != nil {
			log.Fatalf("error getting command: %v", err)
		}
		commandsChan <-command
	}
}

func getMessageType(message string) (string, error) {
	if match, _ := regexp.MatchString(`routing\.queues\.(.+)\.users`, message); match {
		return "routing.queues.{id}.users", nil
	}
	if match, _ := regexp.MatchString(`users\.(.+)\.presence`, message); match {
		return "users.{id}.presence", nil
	}
	if match, _ := regexp.MatchString(`WebSocket\sHeartbeat`, message); match {
		return "WebSocket Heartbeat", nil
	}
	return "", fmt.Errorf("%v", "message type not found")
}