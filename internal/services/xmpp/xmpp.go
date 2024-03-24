package xmpp

import (
	"github.com/mattn/go-xmpp"
)

type XmppService struct {
	username     string
	password     string
	server       string
	recipientJID string
}

func NewXmppService(opts XmppServiceOptions) XmppService {
	return XmppService{
		username:     opts.Username,
		password:     opts.Password,
		server:       opts.Server,
		recipientJID: opts.RecipientJID,
	}
}

type XmppServiceOptions struct {
	Username     string
	Password     string
	Server       string
	RecipientJID string
}

func (xs *XmppService) SendMessage(message string) error {
	// Modify the options to use DIGEST-MD5
	options := xmpp.Options{
		Host:     xs.server,
		User:     xs.username,
		Password: xs.password,
		NoTLS:    true, // Set to true if your server doesn't support TLS
		Debug:    false,
		Session:  true,
		// Mechanism: xmpp.DIGESTMD5,
	}

	client, err := options.NewClient()
	if err != nil {
		return err
	}

	// Send a chat message
	if _, err := client.Send(xmpp.Chat{
		Remote: xs.recipientJID,
		Type:   "chat",
		Text:   message,
	}); err != nil {
		return err
	}

	// Close the connection when done
	return client.Close()
}
