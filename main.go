package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"solace.dev/go/messaging"
	"solace.dev/go/messaging/pkg/solace/config"
	"solace.dev/go/messaging/pkg/solace/message"
	"solace.dev/go/messaging/pkg/solace/resource"
)

// Define Topic Prefix
const TopicPrefix = "events/payment-service"

func MessageHandlerEuro(message message.InboundMessage) {
	var messageBody string

	if payload, ok := message.GetPayloadAsString(); ok {
		messageBody = payload
	} else if payload, ok := message.GetPayloadAsBytes(); ok {
		messageBody = string(payload)
	}

	fmt.Printf("Received Message Body %s \n", messageBody)

	api := slack.New(getEnv("BOT_TOKEN", "xoxb-3517990543552-3498517567635-CX5hwEl01DUYCDXtxSzF40zp"))

	api.PostMessage(getEnv("CHANNEL_ID", "C03EJ6VUTKL"), slack.MsgOptionText("A new user bought a product using card visa with currency is EURO", false))
	api.PostMessage(getEnv("CHANNEL_ID", "C03EJ6VUTKL"), slack.MsgOptionText(messageBody, false))
}

func MessageHandlerUsd(message message.InboundMessage) {
	var messageBody string

	if payload, ok := message.GetPayloadAsString(); ok {
		messageBody = payload
	} else if payload, ok := message.GetPayloadAsBytes(); ok {
		messageBody = string(payload)
	}

	fmt.Printf("Received Message Body %s \n", messageBody)

	api := slack.New(getEnv("BOT_TOKEN", "xoxb-3517990543552-3498517567635-CX5hwEl01DUYCDXtxSzF40zp"))

	api.PostMessage(getEnv("CHANNEL_ID", "C03EJ6VUTKL"), slack.MsgOptionText("A new user bought a product using card visa with currency is USD", false))
	api.PostMessage(getEnv("CHANNEL_ID", "C03EJ6VUTKL"), slack.MsgOptionText(messageBody, false))
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {

	// Configuration parameters
	brokerConfig := config.ServicePropertyMap{
		config.TransportLayerPropertyHost:                getEnv("TransportLayerPropertyHost", "tcps://mrbhn5fvgw72c.messaging.solace.cloud:55443"),
		config.ServicePropertyVPNName:                    getEnv("ServicePropertyVPNName", "payment-broker"),
		config.AuthenticationPropertySchemeBasicUserName: getEnv("AuthenticationPropertySchemeBasicUserName", "solace-cloud-client"),
		config.AuthenticationPropertySchemeBasicPassword: getEnv("AuthenticationPropertySchemeBasicPassword", "sp6c596qno9oq3cdsm80dp4eo4"),
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
		WithSubscriptions(resource.TopicSubscriptionOf(TopicPrefix + "/*/EUR/pm_card_visa/>")).
		Build()

	if err != nil {
		panic(err)
	}

	// Start Direct Message Receiver
	if err := directReceiver.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Direct Receiver running? ", directReceiver.IsRunning())

	//  Build a Direct Message Receiver
	anotherDirectReceiver, err := messagingService.CreateDirectMessageReceiverBuilder().
		WithSubscriptions(resource.TopicSubscriptionOf(TopicPrefix + "/*/USD/pm_card_visa/>")).
		Build()

	if err != nil {
		panic(err)
	}

	// Start another Direct Message Receiver
	if err := anotherDirectReceiver.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Direct Receiver running? ", anotherDirectReceiver.IsRunning())

	for 1 != 0 {

		if regErr := directReceiver.ReceiveAsync(MessageHandlerEuro); regErr != nil {
			panic(regErr)
		}

		if regErr := anotherDirectReceiver.ReceiveAsync(MessageHandlerUsd); regErr != nil {
			panic(regErr)
		}

	}
}
