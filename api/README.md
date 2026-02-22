# 🚀 WhatsApp API - Código Fonte

Código Go para o WhatsApp Multi-Tenant REST API.

## 📁 Arquivos

```
api/
├── main.go              # Servidor HTTP (porta 5000)
├── models.go            # Estruturas (Tenant, Session)
├── session.go           # Gerenciador de sessões
├── handlers.go          # Endpoints HTTP
├── go.mod               # Dependências Go
├── README.md            # Este arquivo
├── MANUAL.md            # Referência API completa ⭐
├── EXEMPLOS.md          # Exemplos cURL
├── .env.example         # Template variáveis
└── data/                # Banco dados (SQLite)
```

## 🚀 Iniciar Rápido

```bash
cd api

# Windows PowerShell
$env:CGO_ENABLED="1"
go run .

# Linux/Mac
CGO_ENABLED=1 go run .
```

API estará em: **http://localhost:5000**

Testar health:
```bash
curl http://localhost:5000/health
```

### 2. Criar um Tenant (Conta)

```bash
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "id": "empresa1",
    "name": "Empresa 1",
    "email": "admin@empresa1.com"
  }'
```

Resposta:
```json
{
  "success": true,
  "message": "Tenant criado com sucesso",
  "data": {
    "id": "empresa1",
    "name": "Empresa1",
    "email": "admin@empresa1.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 3. Criar Sessão (Usuário/Conta WhatsApp)

```bash
curl -X POST http://localhost:8080/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "usuario1"
  }'
```

### 4. Obter QR Code para Conectar

```bash
curl http://localhost:8080/api/tenants/empresa1/sessions/usuario1/qr
```

Resposta com QR Code em Base64:
```json
{
  "success": true,
  "message": "QR Code gerado. Escaneie com seu WhatsApp em 30 segundos",
  "data": {
    "qr_code": "data:image/png;base64,iVBOR...",
    "status": "pending"
  }
}
```

**Escaneie o QR Code com seu WhatsApp para conectar.**

### 5. Verificar Status da Sessão

```bash
curl http://localhost:8080/api/tenants/empresa1/sessions/usuario1/status
```

### 6. Listar Sessões do Tenant

```bash
curl http://localhost:8080/api/tenants/empresa1/sessions
```

### 7. Enviar Mensagem Simples

Após conectar (status = "connected"):

```bash
curl -X POST http://localhost:8080/api/tenants/empresa1/sessions/usuario1/send-message \
  -H "Content-Type: application/json" \
  -d '{
    "number": "5585988888888",
    "message": "Olá, teste da API!"
  }'
```

**Nota:** Use o número completo com código do país (ex: 55 = Brasil, 34 = Espanha, etc)

### 8. Listar Tenants

```bash
curl http://localhost:8080/api/tenants/empresa1
```

### 9. Deletar Sessão

```bash
curl -X DELETE http://localhost:8080/api/tenants/empresa1/sessions/usuario1
```

### 10. Deletar Tenant

```bash
curl -X DELETE http://localhost:8080/api/tenants/empresa1
```

### 11. Health Check

```bash
curl http://localhost:8080/health
```

## Estrutura de Dados

### Tenant
```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Session
```json
{
  "id": "string (user_id)",
  "tenant_id": "string",
  "phone_number": "string",
  "status": "disconnected|qr_pending|connected",
  "created_at": "2024-01-15T10:30:00Z",
  "last_connected_at": "2024-01-15T10:30:00Z"
}
```

## Fluxo de Uso

1. Criar Tenant
2. Criar Session dentro do Tenant
3. Obter QR Code e escanear com WhatsApp
4. Verificar status até estar "connected"
5. Enviar mensagens quando conectado
6. Remover session ou tenant quando não precisar

## Multi-Tenant Example

```bash
# Tenant 1 - Empresa A
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id": "empresaA", "name": "Empresa A", "email": "admin@a.com"}'

# Tenant 2 - Empresa B
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id": "empresaB", "name": "Empresa B", "email": "admin@b.com"}'

# Sessão do Usuário 1 na Empresa A
curl -X POST http://localhost:8080/api/tenants/empresaA/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id": "usuario1"}'

# Sessão de um Usuário da Empresa B
curl -X POST http://localhost:8080/api/tenants/empresaB/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id": "vendedor1"}'
```

## Deployment em Linux

1. Compilar para Linux:
```bash
GOOS=linux GOARCH=amd64 go build -o whatsapp-api ./api
```

2. Transferir o binário para o VPS Linux

3. Execute:
```bash
./whatsapp-api
```

## Notas

- As sessões são persistidas em banco de dados SQLite (`data/sessions/{tenantId}/`)
- Cada sessão mantém seu próprio estado do WhatsApp
- O servidor aguarda conexões em `0.0.0.0:8080` (acessível de até máquinas remotas)
- Para enviar para grupos, use o número do grupo no formato: `120363000000000000-1234567890@g.us`

## Próximas Funcionalidades

- [ ] Suporte a envio de mídia (imagens, vídeos, áudio)
- [ ] Webhook para eventos (mensagens recebidas, etc)
- [ ] Autenticação e autorização
- [ ] Logging persistido
- [ ] Métricas e monitoramento
- [ ] Rate limiting
