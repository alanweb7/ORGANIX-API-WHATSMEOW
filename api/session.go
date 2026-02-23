package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

// SessionManager gerencia tenants e suas sessões
type SessionManager struct {
	mu       sync.RWMutex
	tenants  map[string]*Tenant
	sessions map[string]map[string]*Session // map[tenantID]map[userID]->Session
	dbPath   string
}

func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		tenants:  make(map[string]*Tenant),
		sessions: make(map[string]map[string]*Session),
		dbPath:   "./data/sessions",
	}

	// Criar diretório de dados se não existir
	os.MkdirAll(sm.dbPath, os.ModePerm)

	return sm
}

// CreateTenant cria um novo tenant
func (sm *SessionManager) CreateTenant(id, name, email string) (*Tenant, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.tenants[id]; exists {
		return nil, fmt.Errorf("tenant %s já existe", id)
	}

	tenant := &Tenant{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: sm.now(),
		UpdatedAt: sm.now(),
	}

	sm.tenants[id] = tenant
	sm.sessions[id] = make(map[string]*Session)

	// Criar diretório para o tenant
	tenantDir := filepath.Join(sm.dbPath, id)
	os.MkdirAll(tenantDir, os.ModePerm)

	log.Printf("Tenant criado: %s (%s)", id, name)
	return tenant, nil
}

// GetTenant obtém um tenant
func (sm *SessionManager) GetTenant(id string) (*Tenant, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	tenant, exists := sm.tenants[id]
	if !exists {
		return nil, fmt.Errorf("tenant %s não encontrado", id)
	}

	return tenant, nil
}

// DeleteTenant remove um tenant e suas sessões
func (sm *SessionManager) DeleteTenant(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Desconectar todas as sessões
	if sessions, exists := sm.sessions[id]; exists {
		for _, session := range sessions {
			if session.Client != nil {
				session.Client.Disconnect()
			}
		}
	}

	delete(sm.tenants, id)
	delete(sm.sessions, id)

	// Remover diretório
	tenantDir := filepath.Join(sm.dbPath, id)
	os.RemoveAll(tenantDir)

	log.Printf("Tenant removido: %s", id)
	return nil
}

// CreateSession cria uma nova sessão para um tenant
func (sm *SessionManager) CreateSession(tenantID, userID string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Verificar se tenant existe
	if _, exists := sm.tenants[tenantID]; !exists {
		return nil, fmt.Errorf("tenant %s não encontrado", tenantID)
	}

	// Verificar se sessão já existe
	if sessions, exists := sm.sessions[tenantID]; exists {
		if _, exists := sessions[userID]; exists {
			return nil, fmt.Errorf("sessão %s já existe para tenant %s", userID, tenantID)
		}
	}

	// Preparar cliente WhatsApp
	tenantDir := filepath.Join(sm.dbPath, tenantID)
	ctx := context.Background()
	// Use URI mode and enable foreign keys for SQLite
	dbPath := fmt.Sprintf("file:%s?_foreign_keys=on", filepath.Join(tenantDir, "wa-"+userID+".db"))
	container, err := sqlstore.New(ctx, "sqlite3", dbPath, waLog.Noop)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, err
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	session := &Session{
		ID:        userID,
		TenantID:  tenantID,
		Status:    "disconnected",
		Client:    client,
		CreatedAt: sm.now(),
	}

	sm.sessions[tenantID][userID] = session
	log.Printf("Sessão criada: %s para tenant %s", userID, tenantID)

	return session, nil
}

// GetSession obtém uma sessão
func (sm *SessionManager) GetSession(tenantID, userID string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions, exists := sm.sessions[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant %s não encontrado", tenantID)
	}

	session, exists := sessions[userID]
	if !exists {
		return nil, fmt.Errorf("sessão %s não encontrada", userID)
	}

	return session, nil
}

// GetSessions obtém todas as sessões de um tenant
func (sm *SessionManager) GetSessions(tenantID string) ([]*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions, exists := sm.sessions[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant %s não encontrado", tenantID)
	}

	result := make([]*Session, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, session)
	}

	return result, nil
}

// DeleteSession remove uma sessão
func (sm *SessionManager) DeleteSession(tenantID, userID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessions, exists := sm.sessions[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s não encontrado", tenantID)
	}

	session, exists := sessions[userID]
	if !exists {
		return fmt.Errorf("sessão %s não encontrada", userID)
	}

	if session.Client != nil {
		session.Client.Disconnect()
	}

	delete(sessions, userID)
	log.Printf("Sessão removida: %s de tenant %s", userID, tenantID)

	return nil
}

// SendMessage envia uma mensagem
// TODO: Implementar envio de mensagem com a estrutura correta de waE2E.Message
func (sm *SessionManager) SendMessage(tenantID, userID, phoneNumber, message string) error {
	session, err := sm.GetSession(tenantID, userID)
	if err != nil {
		return err
	}

	if session.Client == nil || session.Status != "connected" {
		return fmt.Errorf("sessão não conectada")
	}

	// Converter número para JID
	jid, err := types.ParseJID(phoneNumber + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("número inválido: %v", err)
	}

	// Validar que o número foi convertido
	if jid.IsEmpty() {
		return fmt.Errorf("JID vazio após conversão")
	}

	log.Printf("Enviando mensagem para %s", jid.String())
	// TODO: Implementar com waE2E.Message
	return fmt.Errorf("funcionalidade em desenvolvimento")
}

func (sm *SessionManager) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, sessions := range sm.sessions {
		for _, session := range sessions {
			if session.Client != nil {
				session.Client.Disconnect()
			}
		}
	}
}

func (sm *SessionManager) now() time.Time {
	return time.Now().UTC()
}
