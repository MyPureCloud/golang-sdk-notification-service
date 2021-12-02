package main

import (
	"flag"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/concurrent_notification_service"
	"github.com/MyPureCloud/golang-sdk-notification-service.git/sequential_notification_service"
	"os"
)

func main() {
	topic1 := os.Getenv("USERS_ID_PRESENCE")
	topic2 := os.Getenv("ROUTING_QUEUES_ID_USERS")

	// NOTE: need to improve. as of now topics are just a list of strings.
	// maybe we could supply topics from a directory containing JSON files
	topics := []string{topic1, topic2}

	ns := flag.String("ns", "", "notification service: supported values: concurrent, sequential")

	flag.Parse()

	if *ns == "" {
		flag.Usage()
	}

	if *ns == "sequential" {
		sns := sequential_notification_service.NewSequentialNotificationService(topics)
		sns.RunSequentialNotificationService()
	}
	if *ns == "concurrent" {
		cns := concurrent_notification_service.NewConcurrentNotificationService(topics)
		cns.RunConcurrentNotificationService()
	}
}
