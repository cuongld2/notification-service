package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/message"
	"solace.dev/go/messaging/pkg/solace/resource"
)

// Define Topic Prefix
const TopicPrefix = "solace/payment"

var messageTransfer string

func MessageHandler(message message.InboundMessage) {
	var messageBody string

	if payload, ok := message.GetPayloadAsString(); ok {
		messageBody = payload
	} else if payload, ok := message.GetPayloadAsBytes(); ok {
		messageBody = string(payload)
	}

	fmt.Printf("Received Message Body %s \n", messageBody)
	messageTransfer = messageBody
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {

	api := slack.New("xoxb-3517990543552-3498517567635-CX5hwEl01DUYCDXtxSzF40zp")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))
	groups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}

	// Configuration parameters
	brokerConfig := config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                "tcps://mrbhn5fvgw72c.messaging.solace.cloud:55443",
		config.ServicePropertyVPNName:                    "payment-broker",
		config.AuthenticationPropertySchemeBasicUserName: "solace-cloud-client",
		config.AuthenticationPropertySchemeBasicPassword: "sp6c596qno9oq3cdsm80dp4eo4",
	}
	messagingService, err := messaging.NewMessagingServiceBuilder().FromConfigurationProvider(brokerConfig).WithTransportSecurityStrategy(config.NewTransportSecurityStrategy().WithoutCertificateValidation()).
		Build()

	if err != nil {
		panic(err)
	}

	// Connect to the messaging serice
	if err := messagingService.Connect(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the broker? ", messagingService.IsConnected())

	//  Build a Direct Message Receiver
	directReceiver, err := messagingService.CreateDirectMessageReceiverBuilder().
		WithSubscriptions(resource.TopicSubscriptionOf(TopicPrefix + "/*/hello/>")).
		Build()

	if err != nil {
		panic(err)
	}

	// Start Direct Message Receiver
	if err := directReceiver.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Direct Receiver running? ", directReceiver.IsRunning())

	messageBody := "Payment intent confirmed has id is :"
	for 1 != 0 {

		if regErr := directReceiver.ReceiveAsync(MessageHandler); regErr != nil {
			panic(regErr)
		}

		if strings.Contains(messageTransfer, messageBody) {

			api.PostMessage("C03EJ6VUTKL", slack.MsgOptionText("A new user bought a product", false))
			api.PostMessage("C03EJ6VUTKL", slack.MsgOptionText(messageTransfer, false))
		}

		messageTransfer = ""

	}

	// // Terminate the Direct Receiver
	// directReceiver.Terminate(2 * time.Second)
	// fmt.Println("\nDirect Receiver Terminated? ", directReceiver.IsTerminated())

	// // Disconnect the Message Service
	// messagingService.Disconnect()
	// fmt.Println("Messaging Service Disconnected? ", !messagingService.IsConnected())
}
