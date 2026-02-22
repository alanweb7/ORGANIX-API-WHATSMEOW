package main

import (
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// Tenant representa uma conta/organização
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session representa um usuário conectado ao WhatsApp dentro de um Tenant
type Session struct {
	ID              string            `json:"id"`
	TenantID        string            `json:"tenant_id"`
	PhoneNumber     string            `json:"phone_number"`
	Status          string            `json:"status"` // "disconnected", "qr_pending", "connected"
	JID             types.JID         `json:"-"`
	Client          *whatsmeow.Client `json:"-"`
	CreatedAt       time.Time         `json:"created_at"`
	LastConnectedAt *time.Time        `json:"last_connected_at"`
}

// QRCodeResponse para responder com dados do QR Code
type QRCodeResponse struct {
	QRCode string `json:"qr_code"`
	Status string `json:"status"`
}

// MessageRequest para requisição de envio de mensagem
type MessageRequest struct {
	Number  string `json:"number"` // Número do destinatário (59899999999)
	Message string `json:"message"`
}

// MediaRequest para requisição de envio de mídia
type MediaRequest struct {
	Number    string `json:"number"`
	FilePath  string `json:"file_path"`
	Caption   string `json:"caption"`
}

// StandardResponse resposta padrão da API
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
