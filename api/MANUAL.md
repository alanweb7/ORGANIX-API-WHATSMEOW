# WhatsApp Multi-Tenant REST API - Manual Completo

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura](#arquitetura)
3. [Instalação](#instalação)
4. [Como Iniciar](#como-iniciar)
5. [Docker](#docker)
6. [Endpoints](#endpoints)
7. [Exemplos de Uso](#exemplos-de-uso)
8. [Fluxo de Conexão](#fluxo-de-conexão)
9. [Multi-Tenant e Multi-Sessão](#multi-tenant-e-multi-sessão)
10. [Deployment em Linux](#deployment-em-linux)
11. [Próximos Passos](#próximos-passos)

---

## 🎯 Visão Geral

API REST completa para gerenciar múltiplas sessões do WhatsApp por tenant (conta).

### Caso de Uso

- **Multi-Tenant**: Uma empresa pode ter N contas WhatsApp (uma para cada filial, departamento, etc)
- **Multi-Sessão**: Cada tenant pode ter N usuários conectados simultaneamente com números WhatsApp diferentes
- **Escalável**: Pronta para VPS Linux com Docker

### Características

✅ Gerenciamento de múltiplas contas (Tenants)  
✅ Múltiplos usuários por conta (Sessões)  
✅ QR Code automático para conexão  
✅ Persistência em SQLite  
✅ API RESTful com JSON  
✅ Health check integrado  
✅ Pronta para produção em Linux  

---

## 🏗️ Arquitetura

### Estrutura de Pastas

```
ORGANIX-API-WHATSMEOW/
├── api/
│   ├── main.go              # Servidor HTTP principal
│   ├── models.go            # Estruturas de dados (Tenant, Session)
│   ├── session.go           # SessionManager (núcleo da lógica)
│   ├── handlers.go          # Handlers HTTP (endpoints)
│   ├── go.mod               # Dependências Go
│   ├── go.sum               # Hash das dependências
│   ├── test.ps1             # Script de teste automático
│   ├── MANUAL.md            # Este arquivo
│   ├── README.md            # Documentação técnica
│   └── EXEMPLOS.md          # Exemplos de requisições
├── whatsmeow/               # Biblioteca WhatsApp (incluída)
└── data/                    # Dados persistidos (criado automaticamente)
    └── sessions/
        ├── empresa1/
        │   └── wa-usuario1.db
        └── empresa2/
            └── wa-usuario2.db
```

### Fluxo de Dados

```
Cliente HTTP
    ↓
API (localhost:5000)
    ↓
Router (Gorilla Mux)
    ↓
Handlers HTTP
    ↓
SessionManager
    ↓ ↓ ↓
[Tenant 1]  [Tenant 2]  [Tenant 3]
   ↓            ↓           ↓
[Session A]  [Session X]  [Session Z]
   ↓            ↓           ↓
WhatsApp    WhatsApp    WhatsApp
(55 98888-1) (55 99999-2) (55 97777-3)
```

---

## 💻 Instalação

### Pré-requisitos

- **Windows**: Go 1.21+, Git
- **Linux**: Go 1.21+, Git, GCC (para compilar SQLite)
- **VPS**: Portas 80, 443 ou 5000 abertas

### Passos de Instalação

#### 1. Clonar Repositório

```bash
git clone https://github.com/tulir/whatsmeow.git
cd ORGANIX-API-WHATSMEOW
```

#### 2. Baixar Dependências

```bash
cd api
go mod download
go mod tidy
```

#### 3. Compilar (Opcional)

```bash
# Windows
go build -o api.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o whatsapp-api .
```

---

## 🚀 Como Iniciar

### Opção 1: Com `go run` (Recomendado para Desenvolvimento)

```bash
cd api
go run .
```

Servidor estará em: **`http://localhost:5000`**

### Opção 2: Com Executável Compilado

```bash
# Windows
.\api.exe

# Linux
./whatsapp-api
```

### Opção 3: Script de Teste Automático (PowerShell)

```powershell
cd api
.\test.ps1
```

O script fará todos os testes automaticamente!

---

## � Docker

Guia rápido para usar Docker. Para documentação completa, veja [DOCKER.md](DOCKER.md).

### Pré-requisitos

```bash
# Instalar Docker
# Windows/Mac: https://www.docker.com/products/docker-desktop
# Linux: sudo apt install docker.io docker-compose
```

### Desenvolvimento com Docker

```bash
# Iniciar em modo desenvolvimento
docker-compose -f docker-compose.dev.yml up

# Em outro terminal
curl http://localhost:5000/health
```

### Build com Docker Compose

```bash
# Produção
docker-compose up -d

# Verificar status
docker-compose ps

# Ver logs
docker-compose logs -f
```

### Build Manual

#### Linux/Mac

```bash
# Build
./build.sh v1.0.0

# Rodar
docker run -p 5000:5000 whatsapp-api:v1.0.0
```

#### Windows

```powershell
# Build
.\build.ps1 -Version "v1.0.0"

# Rodar
docker run -p 5000:5000 whatsapp-api:v1.0.0
```

### Usar Makefile (Recomendado)

```bash
# Build
make build

# Build com versão
make build VERSION=v1.0.0

# Iniciar
make up

# Ver logs
make logs

# Health check
make health

# Deploy
make deploy HOST=user@seu-vps.com

# Ver todas as opções
make help
```

### Push para Registry

```bash
# Docker Hub
docker tag whatsapp-api:v1.0.0 seu-usuario/whatsapp-api:v1.0.0
docker login
docker push seu-usuario/whatsapp-api:v1.0.0

# Registry Private
make push REGISTRY=seu-registry.com/seu-usuario VERSION=v1.0.0
```

### Estrutura Docker

```
api/
├── Dockerfile              # Multi-stage build otimizado
├── .dockerignore          # Arquivos ignorados no build
├── docker-compose.yml     # Produção
├── docker-compose.dev.yml # Desenvolvimento
├── Makefile               # Automação
├── build.sh              # Build script (Linux/Mac)
├── build.ps1             # Build script (Windows)
├── deploy.sh             # Deploy automático
├── k8s-deployment.yaml   # Kubernetes
└── DOCKER.md             # Documentação completa
```

### Variáveis de Ambiente

Copie `.env.example` para `.env`:

```bash
cp .env.example .env
nano .env
```

Depois use no docker-compose:

```yaml
env_file:
  - .env
```

Para mais informações, veja [DOCKER.md](DOCKER.md).

---

## �📡 Endpoints

### 1. Health Check

**GET** `/health`

Verificar se servidor está online.

```bash
curl http://localhost:5000/health
```

**Resposta:**
```json
{
  "status": "ok"
}
```

---

### 2. Gerenciar Tenants

#### Criar Tenant

**POST** `/api/tenants`

Criar uma nova conta/organização.

```bash
curl -X POST http://localhost:5000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "id": "empresa1",
    "name": "Empresa 1",
    "email": "admin@empresa1.com"
  }'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Tenant criado com sucesso",
  "data": {
    "id": "empresa1",
    "name": "Empresa 1",
    "email": "admin@empresa1.com",
    "created_at": "2026-02-22T17:52:21Z",
    "updated_at": "2026-02-22T17:52:21Z"
  }
}
```

---

#### Obter Informações do Tenant

**GET** `/api/tenants/{tenantId}`

```bash
curl http://localhost:5000/api/tenants/empresa1
```

---

#### Deletar Tenant

**DELETE** `/api/tenants/{tenantId}`

Remove o tenant e TODAS as suas sessões.

```bash
curl -X DELETE http://localhost:5000/api/tenants/empresa1
```

---

### 3. Gerenciar Sessões

#### Criar Sessão

**POST** `/api/tenants/{tenantId}/sessions`

Criar uma nova sessão (usuário com WhatsApp).

```bash
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id": "usuario1"}'
```

**Resposta:**
```json
{
  "success": true,
  "message": "Sessão criada, acesse /qr para obter QR Code",
  "data": {
    "id": "usuario1",
    "tenant_id": "empresa1",
    "phone_number": "",
    "status": "disconnected",
    "created_at": "2026-02-22T17:52:21Z",
    "last_connected_at": null
  }
}
```

---

#### Listar Sessões

**GET** `/api/tenants/{tenantId}/sessions`

Listar todas as sessões de um tenant.

```bash
curl http://localhost:5000/api/tenants/empresa1/sessions
```

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": "usuario1",
      "tenant_id": "empresa1",
      "status": "connected",
      "created_at": "2026-02-22T17:52:21Z"
    },
    {
      "id": "usuario2",
      "tenant_id": "empresa1",
      "status": "disconnected",
      "created_at": "2026-02-22T17:53:00Z"
    }
  ]
}
```

---

### 4. QR Code para Conexão

#### Obter QR Code

**GET** `/api/tenants/{tenantId}/sessions/{userId}/qr`

Gera um QR Code em Base64 para conectar o WhatsApp.

```bash
curl http://localhost:5000/api/tenants/empresa1/sessions/usuario1/qr
```

**Resposta:**
```json
{
  "success": true,
  "message": "QR Code gerado. Escaneie com seu WhatsApp em 30 segundos",
  "data": {
    "qr_code": "iVBORw0KGgoAAAANSUhEUgAAAXEAAAFxCAYAAAC7...",
    "status": "pending"
  }
}
```

### Como Usar o QR Code

1. Copie o valor de `qr_code`
2. Abra https://www.base64-image-decoder.com/
3. Cole a string base64
4. Clique "Decode to Image"
5. Escaneie com seu WhatsApp em 30 segundos
6. Verifique o status da conexão

---

#### Verificar Status da Sessão

**GET** `/api/tenants/{tenantId}/sessions/{userId}/status`

Verifica se a sessão está conectada.

```bash
curl http://localhost:5000/api/tenants/empresa1/sessions/usuario1/status
```

**Resposta (Desconectado):**
```json
{
  "success": true,
  "data": {
    "user_id": "usuario1",
    "status": "disconnected",
    "jid": ""
  }
}
```

**Resposta (Conectado):**
```json
{
  "success": true,
  "data": {
    "user_id": "usuario1",
    "status": "connected",
    "jid": "5585988888888@s.whatsapp.net"
  }
}
```

---

### 5. Enviar Mensagens

#### Enviar Mensagem de Texto

**POST** `/api/tenants/{tenantId}/sessions/{userId}/send-message`

Envia uma mensagem de texto para um contato.

```bash
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions/usuario1/send-message \
  -H "Content-Type: application/json" \
  -d '{
    "number": "5585999999999",
    "message": "Olá! Teste da API WhatsApp"
  }'
```

**Parâmetros:**
- `number`: Número do WhatsApp (com código do país, ex: 55 = Brasil)
- `message`: Texto da mensagem

**Resposta:**
```json
{
  "success": true,
  "message": "Mensagem enviada com sucesso"
}
```

> **⚠️ Nota**: Envio de mensagens está em desenvolvimento. Será implementado em breve.

---

#### Deletar Sessão

**DELETE** `/api/tenants/{tenantId}/sessions/{userId}`

Remove uma sessão e desconecta o WhatsApp.

```bash
curl -X DELETE http://localhost:5000/api/tenants/empresa1/sessions/usuario1
```

---

## 📋 Exemplos de Uso

### Exemplo 1: Uma Empresa com Um Usuário

```bash
# 1. Criar empresa
curl -X POST http://localhost:5000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa1","name":"Minha Empresa","email":"admin@empresa.com"}'

# 2. Criar usuário/sessão
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"vendedor1"}'

# 3. Obter QR Code
curl http://localhost:5000/api/tenants/empresa1/sessions/vendedor1/qr

# 4. Escanear com WhatsApp (30 segundos)

# 5. Verificar conexão
curl http://localhost:5000/api/tenants/empresa1/sessions/vendedor1/status

# 6. Enviar mensagem
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions/vendedor1/send-message \
  -H "Content-Type: application/json" \
  -d '{"number":"5585999999999","message":"Olá!"}'
```

---

### Exemplo 2: Uma Empresa com Múltiplos Usuários

```bash
# Criar empresa
curl -X POST http://localhost:5000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa1","name":"Empresa","email":"admin@empresa.com"}'

# Sessão Vendedor 1
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions \
  -d '{"user_id":"vendedor1"}'
curl http://localhost:5000/api/tenants/empresa1/sessions/vendedor1/qr
# Escanear...

# Sessão Vendedor 2
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions \
  -d '{"user_id":"vendedor2"}'
curl http://localhost:5000/api/tenants/empresa1/sessions/vendedor2/qr
# Escanear...

# Sessão Suporte
curl -X POST http://localhost:5000/api/tenants/empresa1/sessions \
  -d '{"user_id":"suporte1"}'
curl http://localhost:5000/api/tenants/empresa1/sessions/suporte1/qr
# Escanear...

# Resultado: 3 números WhatsApp diferentes conectados simultaneamente
curl http://localhost:5000/api/tenants/empresa1/sessions
```

---

### Exemplo 3: Múltiplas Empresas (Multi-Tenant)

```bash
# Empresa A
curl -X POST http://localhost:5000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa_a","name":"Empresa A","email":"admin@a.com"}'
curl -X POST http://localhost:5000/api/tenants/empresa_a/sessions \
  -d '{"user_id":"gerente_a"}'
curl http://localhost:5000/api/tenants/empresa_a/sessions/gerente_a/qr
# Escanear com número de A...

# Empresa B
curl -X POST http://localhost:5000/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa_b","name":"Empresa B","email":"admin@b.com"}'
curl -X POST http://localhost:5000/api/tenants/empresa_b/sessions \
  -d '{"user_id":"gerente_b"}'
curl http://localhost:5000/api/tenants/empresa_b/sessions/gerente_b/qr
# Escanear com número de B...

# Resultado: 2 grandes empresas com 2 números WhatsApp diferentes
#            gerenciadas pela mesma API!
```

---

## 🔄 Fluxo de Conexão

### Passo a Passo

```
1. Criar Tenant (conta/organização)
        ↓
2. Criar Session (usuário/número)
        ↓
3. Obter QR Code
        ↓
4. Decodificar Base64
        ↓
5. Escanear com WhatsApp (30 segundos)
        ↓
6. Verificar Status (aguardar "connected")
        ↓
7. Enviar Mensagens (quando conectado)
        ↓
8. Manter sessão ativa ou deletar quando não usar
```

### Diagrama Visual

```
┌─────────────────┐
│   Criar Tenant  │ POST /api/tenants
└────────┬────────┘
         │
         ↓
┌─────────────────────┐
│  Criar Sessão       │ POST /api/tenants/{id}/sessions
└────────┬────────────┘
         │
         ↓
┌─────────────────────┐
│  Gerar QR Code      │ GET /api/tenants/{id}/sessions/{userId}/qr
└────────┬────────────┘
         │
         ↓
┌──────────────────────────┐
│  Escanear com WhatsApp   │ 30 segundos
└────────┬─────────────────┘
         │
         ↓
┌──────────────────────────┐
│  Verificar Conexão       │ GET /api/tenants/{id}/sessions/{userId}/status
│  (aguardar "connected")  │
└────────┬─────────────────┘
         │
         ↓
┌──────────────────────────┐
│  Enviar Mensagens        │ POST /api/tenants/{id}/sessions/{userId}/send-message
└──────────────────────────┘
```

---

## 🏢 Multi-Tenant e Multi-Sessão

### O Que É Multi-Tenant?

**Multi-Tenant** significa que a API pode gerenciar múltiplas contas/organizações independentes.

```
API WhatsApp
├── Empresa A (tenant: "empresa_a")
│   ├── Gerente (usuario: "gerente_a") → WhatsApp: (55) 98888-1111
│   └── Vendedor (usuario: "vendedor_a") → WhatsApp: (55) 99999-2222
│
├── Empresa B (tenant: "empresa_b")
│   ├── Suporte (usuario: "suporte_b") → WhatsApp: (55) 97777-3333
│   └── Vendedor (usuario: "vendedor_b") → WhatsApp: (55) 96666-4444
│
└── Consultoria XYZ (tenant: "consultoria_xyz")
    └── Proprietário (usuario: "owner") → WhatsApp: (55) 95555-5555
```

**Total**: 3 tenants, 5 sessões, 5 números WhatsApp diferentes conectados simultaneamente! 🎉

### O Que É Multi-Sessão?

**Multi-Sessão** significa que cada tenant pode ter múltiplos usuários conectados com números diferentes.

```
Empresa A
├── Gerente: (55) 98888-1111
├── Vendedor 1: (55) 99999-2222
├── Vendedor 2: (55) 98888-3333
├── Suporte: (55) 97777-4444
└── Administrativo: (55) 96666-5555
```

**Benefício**: Um único servidor gerenciando 5 números WhatsApp da mesma empresa!

---

## 🐧 Deployment em Linux

### Pré-requisitos

```bash
# Ubuntu/Debian
sudo apt-get install golang-go git gcc

# CentOS/RHEL
sudo yum install golang git gcc
```

### Compilar para Linux

No Windows ou Mac:

```bash
cd api
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o whatsapp-api .
```

### Transferir para VPS

```bash
# De sua máquina local
scp whatsapp-api user@seu-vps.com:/app/
scp -r data/ user@seu-vps.com:/app/  # Se quiser manter dados

# Conectar e executar
ssh user@seu-vps.com
cd /app
chmod +x whatsapp-api
./whatsapp-api
```

### Usar Docker (Recomendado)

Já existe um `Dockerfile` na pasta `api`:

```bash
# Build
docker build -t whatsapp-api .

# Run
docker run -p 5000:5000 -v $(pwd)/data:/app/data whatsapp-api
```

Com Docker Compose:

```bash
docker-compose up -d
```

### Usar Systemd (Iniciar Automaticamente)

Criar `/etc/systemd/system/whatsapp-api.service`:

```ini
[Unit]
Description=WhatsApp REST API
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/app
ExecStart=/app/whatsapp-api
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Ativar:

```bash
sudo systemctl daemon-reload
sudo systemctl enable whatsapp-api
sudo systemctl start whatsapp-api
```

### Usar Nginx como Reverse Proxy

```nginx
server {
    listen 80;
    server_name seu-dominio.com;

    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Reiniciar Nginx:

```bash
sudo systemctl restart nginx
```

Agora a API estará acessível em: `http://seu-dominio.com`

---

## 🔮 Próximos Passos

### Funcionalidades a Implementar

- ✅ Multi-tenant (PRONTO)
- ✅ Multi-sessão (PRONTO)
- ✅ QR Code (PRONTO)
- ⏳ Envio de mensagens (Em desenvolvimento)
- ⏳ Receber mensagens (Webhooks)
- ⏳ Envio de mídia (Imagens, vídeos)
- ⏳ Autenticação JWT/OAuth
- ⏳ Rate limiting
- ⏳ Logging estruturado
- ⏳ Monitoramento com Prometheus
- ⏳ PostgreSQL em vez de SQLite
- ⏳ Redis para cache/sessões

### Integrações Recomendadas

1. **Chatbot**: Integre com seu chatbot favorito
2. **CRM**: Conecte com Salesforce, Pipedrive, etc
3. **E-commerce**: Acople com WooCommerce, Shopify
4. **Sistemas Legados**: Use a API para integrar sistemas antigos
5. **Apps Mobile**: Crie apps que usam esta API como backend

###Sugestões de Scaling

- Use Kubernetes para orquestrações
- Implemente load balancing com HAProxy
- Use PostgreSQL para múltiplas instâncias
- Implemente cache com Redis
- Monitore com ELK Stack ou DataDog

---

## 📞 Suporte

### Problemas Comuns

**P: QR Code não aparece?**  
R: Verifique se a sessão foi criada corretamente e se o servidor está rodando.

**P: "Sessão não conectada" ao enviar mensagem?**  
R: Aguarde o status ficar "connected" após escanear o QR Code.

**P: Erro de porta em uso?**  
R: Mude a porta em `main.go` (atualmente 5000) ou feche a aplicação anterior.

**P: Erro de SQLite no Windows?**  
R: Execute `go run .` em vez de compilar manualmente.

### Comandos Úteis

```bash
# Ver logs do servidor
tail -f server.log

# Matar processo Go
pkill -f "go run"

# Ver portas em uso
netstat -an | grep 5000

# Testar API
curl -v http://localhost:5000/health

# Ver estrutura de dados
ls -la data/sessions/
```

---

## 📚 Referências

- [Whatsmeow GitHub](https://github.com/tulir/whatsmeow)
- [Go Documentation](https://golang.org/doc/)
- [REST API Best Practices](https://restfulapi.net/)
- [Docker Documentation](https://docs.docker.com/)

---

## 📄 Licença

Este projeto usa a biblioteca whatsmeow que está sob licença MPL 2.0.

---

## ✍️ Versão

**v1.0.0** - 22/02/2026

- ✅ Multi-tenant funcionando
- ✅ Multi-sessão funcionando
- ✅ QR Code funcionando
- ✅ Gerenciamento de sessões completo

---

**Desenvolvido com ❤️ para gerenciar WhatsApp em escala**
