package main

import (
	"context"
	"deltachat-bot/deltachat"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
)

type Server struct {
	addr                         string
	deltaChatClient              *deltachat.Client
	deltaChatNotificationAddress string
}

func NewServer(addr string, deltaChatClient *deltachat.Client, deltaChatNotificationAddress string) *Server {
	return &Server{addr: addr, deltaChatClient: deltaChatClient, deltaChatNotificationAddress: deltaChatNotificationAddress}
}

func (s *Server) Serve(ctx context.Context) {
	http.HandleFunc("/", s.sendMessageHandler)
	server := &http.Server{
		Addr:    s.addr,
		Handler: nil,
	}
	go func() {
		<-ctx.Done()
		log.Println("shutting down server")
		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("error shutting down server", "error", err)
		}
		log.Println("server shut down")
	}()
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", "error", err)
	}
}

type SendMessageRequest struct {
	Text string `json:"text"`
}

func (s *Server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var request SendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.deltaChatClient.SendMessage(s.deltaChatNotificationAddress, request.Text)
	if err != nil {
		slog.Error("error sending message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
