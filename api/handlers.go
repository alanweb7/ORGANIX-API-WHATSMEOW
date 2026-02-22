package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// createTenant cria um novo tenant
func createTenant(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		tenant, err := sm.CreateTenant(req.ID, req.Name, req.Email)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusCreated, StandardResponse{
			Success: true,
			Message: "Tenant criado com sucesso",
			Data:    tenant,
		})
	}
}

// getTenant obtém informações de um tenant
func getTenant(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]

		tenant, err := sm.GetTenant(tenantID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Data:    tenant,
		})
	}
}

// deleteTenant remove um tenant
func deleteTenant(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]

		if err := sm.DeleteTenant(tenantID); err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Message: "Tenant removido com sucesso",
		})
	}
}

// getSessions obtém todas as sessões de um tenant
func getSessions(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]

		sessions, err := sm.GetSessions(tenantID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Data:    sessions,
		})
	}
}

// createSession cria uma nova sessão
func createSession(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]

		var req struct {
			UserID string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		session, err := sm.CreateSession(tenantID, req.UserID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusCreated, StandardResponse{
			Success: true,
			Message: "Sessão criada, acesse /qr para obter QR Code",
			Data:    session,
		})
	}
}

// getQRCode obtém o QR Code para conectar
func getQRCode(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]
		userID := mux.Vars(r)["userId"]

		session, err := sm.GetSession(tenantID, userID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		if session.Client == nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   "Cliente não inicializado",
			})
			return
		}

		// Conectar e obter QR Code
		qrChan, err := session.Client.GetQRChannel(r.Context())
		if err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   fmt.Sprintf("Erro ao obter QR Code: %v", err),
			})
			return
		}

		err = session.Client.Connect()
		if err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   fmt.Sprintf("Erro ao conectar: %v", err),
			})
			return
		}

		// Aguardar QR Code
		for evt := range qrChan {
			if evt.Code != "" {
				qrBase64 := base64.StdEncoding.EncodeToString([]byte(evt.Code))

				writeJSON(w, http.StatusOK, StandardResponse{
					Success: true,
					Message: "QR Code gerado. Escaneie com seu WhatsApp em 30 segundos",
					Data: QRCodeResponse{
						QRCode: qrBase64,
						Status: "pending",
					},
				})
				return
			}
		}

		writeJSON(w, http.StatusInternalServerError, StandardResponse{
			Success: false,
			Error:   "Falha ao gerar QR Code",
		})
	}
}

// getSessionStatus obtém status da sessão
func getSessionStatus(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]
		userID := mux.Vars(r)["userId"]

		session, err := sm.GetSession(tenantID, userID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		status := "disconnected"
		if session.Client != nil && session.Client.IsConnected() {
			status = "connected"
			session.Status = "connected"
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Data: map[string]interface{}{
				"user_id": userID,
				"status":  status,
				"jid":     session.JID.String(),
			},
		})
	}
}

// deleteSession remove uma sessão
func deleteSession(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]
		userID := mux.Vars(r)["userId"]

		if err := sm.DeleteSession(tenantID, userID); err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Message: "Sessão removida com sucesso",
		})
	}
}

// sendMessage envia uma mensagem
func sendMessage(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]
		userID := mux.Vars(r)["userId"]

		var req MessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		if err := sm.SendMessage(tenantID, userID, req.Number, req.Message); err != nil {
			log.Printf("Erro ao enviar mensagem: %v", err)
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   fmt.Sprintf("Erro ao enviar mensagem: %v", err),
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Message: "Mensagem enviada com sucesso",
		})
	}
}

// sendMedia envia mídia
func sendMedia(sm *SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := mux.Vars(r)["tenantId"]
		userID := mux.Vars(r)["userId"]

		var req MediaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		// Verificar se arquivo existe
		_, err := os.ReadFile(req.FilePath)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   fmt.Sprintf("Arquivo não encontrado: %v", err),
			})
			return
		}

		session, err := sm.GetSession(tenantID, userID)
		if err != nil {
			writeJSON(w, http.StatusNotFound, StandardResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		if session.Client == nil || session.Status != "connected" {
			writeJSON(w, http.StatusBadRequest, StandardResponse{
				Success: false,
				Error:   "Sessão não conectada",
			})
			return
		}

		writeJSON(w, http.StatusOK, StandardResponse{
			Success: true,
			Message: "Funcionalidade em desenvolvimento",
		})
	}
}

// Função auxiliar para escrever JSON
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
