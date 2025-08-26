package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/jgimenez/deltachat-webhook-bot/deltachat"
)

type Server struct {
	addr            string
	deltaChatClient *deltachat.Client
}

func NewServer(addr string, deltaChatClient *deltachat.Client) *Server {
	return &Server{addr: addr, deltaChatClient: deltaChatClient}
}

func (s *Server) Serve(ctx context.Context) {
	http.HandleFunc("/{destination}", s.sendMessageHandler)
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
	destination := r.PathValue("destination")
	// destination check to avoid spam
	_, err := s.deltaChatClient.FindContact(destination)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var request SendMessageRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if request.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.deltaChatClient.SendMessage(destination, request.Text)
	if err != nil {
		slog.Error("error sending message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
