package deltachat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat/option"

	"github.com/deltachat/deltachat-rpc-client-go/deltachat"
	"github.com/deltachat/deltachat-rpc-client-go/deltachat/transport"
)

type Client struct {
	bot   *deltachat.Bot
	trans *transport.IOTransport
	me    deltachat.AccountId
}

func New() (*Client, error) {
	trans := transport.NewIOTransport()
	err := trans.Open()
	if err != nil {
		return nil, fmt.Errorf("opening transport: %w", err)
	}
	rpc := &deltachat.Rpc{Context: context.Background(), Transport: trans}
	bot := deltachat.NewBot(rpc)
	me := deltachat.GetAccount(rpc)

	return &Client{bot: bot, trans: trans, me: me}, nil
}

func (c *Client) Start() error {
	go func() {
		err := runBot(c.bot, c.me)
		if err != nil {
			slog.Error("error running bot", "err", err)
		}
	}()
	return nil
}

func (c *Client) Close() error {
	c.bot.Stop()
	c.trans.Close()
	return nil
}

func logEvent(bot *deltachat.Bot, accId deltachat.AccountId, event deltachat.Event) {
	switch ev := event.(type) {
	case deltachat.EventInfo:
		slog.Debug("DeltaChat", "msg", ev.Msg)
	case deltachat.EventWarning:
		slog.Warn("DeltaChat", "msg", ev.Msg)
	case deltachat.EventError:
		slog.Error("DeltaChat", "msg", ev.Msg)
	}
}

func runBot(bot *deltachat.Bot, accId deltachat.AccountId) error {
	sysinfo, _ := bot.Rpc.GetSystemInfo()
	for k, v := range sysinfo {
		slog.Info("Deltachat core info", "key", k, "value", v)
	}

	bot.On(deltachat.EventInfo{}, logEvent)
	bot.On(deltachat.EventWarning{}, logEvent)
	bot.On(deltachat.EventError{}, logEvent)
	bot.OnUnhandledEvent(logEvent)
	bot.OnNewMsg(func(bot *deltachat.Bot, accId deltachat.AccountId, msgId deltachat.MsgId) {
		msg, _ := bot.Rpc.GetMessage(accId, msgId)
		if msg.FromId > deltachat.ContactLastSpecial {
			slog.Info("Received message", "from", msg.FromId, "text", msg.Text)
		}
	})

	if isConf, _ := bot.Rpc.IsConfigured(accId); !isConf {
		return errors.New("account not configured")
	}

	addr, _ := bot.Rpc.GetConfig(accId, "configured_addr")
	slog.Info("DeltaChat listening at", "addr", addr.Unwrap())
	err := bot.Run()
	if err != nil {
		return err
	}
	return nil
}

var ErrNotFound = errors.New("not found")

func (c *Client) ImportBackup(filename string) error {
	if isConf, _ := c.bot.Rpc.IsConfigured(c.me); isConf {
		// already configured
		return nil
	}

	err := c.bot.Rpc.ImportBackup(c.me, filename, option.None[string]())
	if err != nil {
		return fmt.Errorf("importing backup: %w", err)
	}
	err = c.bot.Rpc.SetConfig(c.me, "bot", option.Some("1"))
	if err != nil {
		log.Fatalln(err)
	}
	err = c.bot.Rpc.Configure(c.me)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func (c *Client) FindContact(address string) (deltachat.ContactId, error) {
	ids, err := c.bot.Rpc.GetContactIds(c.me, uint(deltachat.ContactFlagVerifiedOnly), option.None[string]())
	if err != nil {
		return 0, fmt.Errorf("getting contact ids: %w", err)
	}
	for _, id := range ids {
		contact, err := c.bot.Rpc.GetContact(c.me, id)
		if err != nil {
			return 0, fmt.Errorf("error getting contact: %w", err)
		}
		if contact.Address == address {
			return id, nil
		}
	}
	return 0, ErrNotFound
}

func (c *Client) SendMessage(toAddress string, text string) error {
	slog.Info("Sending message", "text", text)
	contactID, err := c.FindContact(toAddress)
	if err != nil {
		return fmt.Errorf("looking up contact: %w", err)
	}
	chatID, err := c.bot.Rpc.CreateChatByContactId(c.me, contactID)
	_, err = c.bot.Rpc.MiscSendTextMessage(c.me, chatID, text)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}
	return nil
}
