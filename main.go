package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"github.com/jgimenez/deltachat-webhook-bot/deltachat"
)

type ServerOptions struct {
	ListenAddr    string
	ImportAccount string
}

func main() {
	opts := parseOptions()

	deltaChatBot, err := deltachat.New()
	if err != nil {
		slog.Error("could not create deltachat client", "err", err)
		return
	}

	if opts.ImportAccount != "" {
		err := deltaChatBot.ImportBackup(opts.ImportAccount)
		if err != nil {
			slog.Error("could not import backup", "err", err)
			return
		}
	}

	err = deltaChatBot.Start()
	if err != nil {
		slog.Error("could not start deltachat client", "err", err)
		return
	}

	// make cancellable context that cancels with ctrl+c
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		slog.Info("ctrl+c received, shutting down")
		cancel()
	}()

	server := NewServer(opts.ListenAddr, deltaChatBot)
	server.Serve(ctx) // blocks

	err = deltaChatBot.Close()
	if err != nil {
		slog.Error("could not close deltachat client", "error", err)
	}
}

func parseOptions() ServerOptions {
	opts := ServerOptions{
		ListenAddr:    getEnvOr("DELTA_CHAT_BOT_LISTEN_ADDR", ":8080"),
		ImportAccount: getEnv("DELTA_CHAT_BOT_IMPORT_ACCOUNT"),
	}
	return opts
}

// getEnv reads an environment variable or a file with the same name and the suffix _FILE (for example, ASANA_ARCHIVAL_CLIENT_SECRET_FILE)
// if the file exists, it is read and the content is returned
// if the file does not exist, the environment variable is returned
// if the environment variable is not set, an empty string is returned
func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := os.Getenv(key + "_FILE"); value != "" {
		content, err := os.ReadFile(value)
		if err != nil {
			log.Fatal(err)
		}
		return string(content)
	}
	return ""
}

func getEnvOr(key string, defaultValue string) string {
	result := getEnv(key)
	if result == "" {
		return defaultValue
	}
	return result
}
