package sendxmpp

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattn/go-xmpp"
)

func SendXMPP(message string) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}
	username := os.Getenv("XMPP_USERNAME")
	password := os.Getenv("XMPP_PASSWORD")
	server := os.Getenv("XMPP_SERVER")
	recipientJID := os.Getenv("XMPP_RECIPIENT")

	if username == "" || password == "" || server == "" || recipientJID == "" {
		fmt.Println("Please set XMPP_USERNAME, XMPP_PASSWORD, XMPP_SERVER or XMPP_RECIPIENT environment variables")
		os.Exit(1)
	}

	// Modify the options to use DIGEST-MD5
	options := xmpp.Options{
		Host:     server,
		User:     username,
		Password: password,
		NoTLS:    true, // Set to true if your server doesn't support TLS
		Debug:    false,
		Session:  true,
		// Mechanism: xmpp.DIGESTMD5,
	}

	client, err := options.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Send a chat message
	_, err = client.Send(xmpp.Chat{
		Remote: recipientJID,
		Type:   "chat",
		Text:   message,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message sent successfully!")

	// Close the connection when done
	client.Close()
}
