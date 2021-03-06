package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paked/configure"
	"github.com/paked/messenger"
)

var (
	conf        = configure.New()
	verifyToken = conf.String("verify-token", "", "The token used to verify facebook")
	verify      = conf.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken   = conf.String("page-token", "", "The token that is used to verify the page on facebook")
	appSecret   = conf.String("app-secret", "", "The app secret from the facebook developer portal")
	port        = conf.Int("port", 8080, "The port used to serve the messenger bot")
)

func main() {
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	conf.Use(configure.NewJSONFromFile("config.json"))

	conf.Parse()

	// Create a new messenger client
	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	// Setup a handler to be triggered when a message is received
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		/*
			p, err := client.ProfileByID(m.Sender.ID)
			if err != nil {
				fmt.Println("Something went wrong!", err)
			}
		*/

		/*
			button := messenger.StructuredMessageButton{}
			button.Title = "title"
			button.FallbackURL = "fallback URL"
			button.MessengerExtensions = false
			button.Payload = "what is a payload?"
			button.Type = "what type?"
			button.URL = "http://whuut.com"
			button.WebviewHeightRatio = "3"
			button.WebviewShareButton = "share?"
			buttons := []messenger.StructuredMessageButton{button}
			r.ButtonTemplate("Hello", &buttons)
		*/

		reply := messenger.QuickReply{}
		reply.Payload = "Payload"
		reply.Title = "Title"
		reply.ContentType = "text"
		replies := []messenger.QuickReply{reply, reply, reply}
		r.TextWithReplies("Replies", replies)

		button := messenger.StructuredMessageButton{}
		button.Type = "web_url"
		button.URL = "http://nambataxi.kg"
		button.Title = "Namba Taxi Website"
		buttons := []messenger.StructuredMessageButton{button}
		r.ButtonTemplate("Buttons", &buttons)

	})

	// Setup a handler to be triggered when a message is delivered
	client.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		fmt.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
	})

	// Setup a handler to be triggered when a message is read
	client.HandleRead(func(m messenger.Read, r *messenger.Response) {
		fmt.Println("Read at:", m.Watermark().Format(time.UnixDate))
	})

	fmt.Printf("Serving messenger bot on localhost:%d\n", *port)

	http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), client.Handler())
}
