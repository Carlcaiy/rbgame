package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func main() {
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "tongits-x",
	}, option.WithCredentialsFile("refresh_token.json"))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	fmt.Println("app")
	cli, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error initializing messaging: %v\n", err)
	}
	fmt.Println("cli")
	res, err := cli.SendAll(context.Background(), []*messaging.Message{
		{
			Token: "c3lGFNwfSi6fH2nNGxIsIY:APA91bGj-jzW_e1Ue8MNbunuY8hvqU-BJnr_ZqKgjgHelPQe3ogw8lSSNAdFlFgY_cG67Wpq_3VVct8_n5f13reE8EdDGx6ApcRf9XggbdR_oZ-KkqEnHFgmtXNqveqwKVq8RurkbxNv",
			Notification: &messaging.Notification{
				Title: "hello",
				Body:  "this is test",
			},
			Data: map[string]string{
				"Id":   "1",
				"Type": "anytype",
				"AB":   "",
			},
		},
	})
	if err != nil {
		log.Fatalf("error send messaging: %v\n", err)
	}
	fmt.Printf("%#v", res)
}
