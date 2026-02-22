package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Inicializar gerenciador de sessões
	sessionManager := NewSessionManager()
	defer sessionManager.Close()

	// Criar router
	router := mux.NewRouter()

	// Endpoints do gerenciador de tenants
	router.HandleFunc("/api/tenants", createTenant(sessionManager)).Methods("POST")
	router.HandleFunc("/api/tenants/{tenantId}", getTenant(sessionManager)).Methods("GET")
	router.HandleFunc("/api/tenants/{tenantId}", deleteTenant(sessionManager)).Methods("DELETE")

	// Endpoints de sessões WhatsApp
	router.HandleFunc("/api/tenants/{tenantId}/sessions", getSessions(sessionManager)).Methods("GET")
	router.HandleFunc("/api/tenants/{tenantId}/sessions", createSession(sessionManager)).Methods("POST")
	router.HandleFunc("/api/tenants/{tenantId}/sessions/{userId}/qr", getQRCode(sessionManager)).Methods("GET")
	router.HandleFunc("/api/tenants/{tenantId}/sessions/{userId}/status", getSessionStatus(sessionManager)).Methods("GET")
	router.HandleFunc("/api/tenants/{tenantId}/sessions/{userId}", deleteSession(sessionManager)).Methods("DELETE")

	// Endpoints de mensagens
	router.HandleFunc("/api/tenants/{tenantId}/sessions/{userId}/send-message", sendMessage(sessionManager)).Methods("POST")
	router.HandleFunc("/api/tenants/{tenantId}/sessions/{userId}/send-media", sendMedia(sessionManager)).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	}).Methods("GET")

	// Configurar servidor
	server := &http.Server{
		Addr:         ":5000",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Servidor iniciado em http://localhost:5000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinais de encerramento
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Encerrando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}

	log.Println("Servidor encerrado")
}
