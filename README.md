# 🚀 WhatsApp API Multi-Tenant & Multi-Sessão

API REST profissional para automação do WhatsApp com suporte a múltiplas empresas (tenants) e múltiplas sessões por empresa.

## 📚 Documentação Completa

> **👉 Comece aqui**: Leia primeiro [README-RESUMO.md](README-RESUMO.md) para uma visão geral

### 📖 Documentação por Tipo

| Documento | Para | Leia | 
|-----------|------|------|
| **[README-RESUMO.md](README-RESUMO.md)** | Visão geral | Todos |
| **[MANUAL.md](api/MANUAL.md)** | Referência API completa | Developers |
| **[docker/README.md](docker/README.md)** | Docker organizado | Todos |
| **[DOCKER.md](DOCKER.md)** | Docker detalhado | DevOps |
| **[DEPLOYMENT-CHECKLIST.md](DEPLOYMENT-CHECKLIST.md)** | Deploy em produção | DevOps/SysAdmins |
| **[ESTRUTURA.md](ESTRUTURA.md)** | Arquitetura projeto | Todos |
| **[QUICK-REFERENCE.md](QUICK-REFERENCE.md)** | Comandos rápidos | Developers |

---

## 🎯 Começar Rápido (3 passos)

### 1️⃣ Desenvolvimento Local

```bash
# Windows PowerShell
cd api
$env:CGO_ENABLED="1"
go run .

# Linux/Mac
cd api
go run .
```

Teste: `curl http://localhost:5000/health`

### 2️⃣ Com Docker (Recomendado)

```bash
cd docker

# Windows
.\build.ps1 -Version "1.0.0"
docker-compose -f docker-compose.yml up -d

# Linux/Mac
chmod +x build.sh start.sh
./build.sh 1.0.0
docker-compose -f docker-compose.yml up -d
```

Teste: `curl http://localhost:5000/health`

### 3️⃣ Usando Make (Mais fácil)

```bash
cd docker
make help          # Ver todos os comandos
make build VERSION=1.0.0
make up
make health
```

---

## 📁 Estrutura do Projeto

```
ORGANIX-API-WHATSMEOW/
│
├── 📘 DOCUMENTAÇÃO (Raiz)
│   ├── README-RESUMO.md          👈 COMECE AQUI
│   ├── MANUAL.md (em api/)       API completa
│   ├── DOCKER.md                 Docker detalhado
│   ├── DEPLOYMENT-CHECKLIST.md   Deploy passo-a-passo
│   ├── ESTRUTURA.md              Arquitetura
│   ├── QUICK-REFERENCE.md        Comandos rápidos
│   └── nginx.conf.template       Nginx pronto
│
├── 🚀 API (Aplicação Go)
│   └── api/
│       ├── main.go               Servidor HTTP
│       ├── models.go             Estruturas de dados
│       ├── session.go            Gerenciador
│       ├── handlers.go           Endpoints HTTP
│       ├── go.mod                Dependências
│       └── MANUAL.md             Referência API
│
├── 🐳 DOCKER (Tudo isolado aqui!)
│   └── docker/
│       ├── Dockerfile            Imagem multi-stage
│       ├── docker-compose.yml    Produção
│       ├── docker-compose.dev.yml Dev com hot reload
│       ├── .dockerignore         Excludes
│       ├── build.sh              Build Linux/Mac
│       ├── build.ps1             Build Windows
│       ├── deploy.sh             Deploy VPS
│       ├── start.sh              Menu Linux/Mac
│       ├── start.cmd             Menu Windows
│       ├── Makefile              Automation
│       └── README.md             Docker docs
│
├── 🔧 CI/CD (GitHub Actions)
│   └── .github/workflows/
│       └── build-deploy.yml      Automação completa
│
├── 📚 WHATSAPP (Biblioteca)
│   └── whatsmeow/                Código whatsmeow
│
└── LICENSE
```

---

## 🎨 Arquitetura

```
┌─────────────────────────────────┐
│  Cliente (Web/Mobile/CLI)       │
└────────────┬────────────────────┘
             │ HTTP/HTTPS
             ↓
┌─────────────────────────────────┐
│  Nginx Reverse Proxy            │
│  (SSL, Rate Limit, Load Bal)    │
└────────────┬────────────────────┘
             │ HTTP
             ↓
┌─────────────────────────────────┐
│  WhatsApp API (Docker)          │
│  ┌───────────────────────────┐  │
│  │ SessionManager            │  │
│  ├─ Tenant 1                 │  │
│  │  ├─ Session 1 → WhatsApp  │  │
│  │  └─ Session 2 → WhatsApp  │  │
│  ├─ Tenant 2                 │  │
│  │  └─ Session 1 → WhatsApp  │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

---

## 🚀 Opções de Execução

### Local (Rápido, sem Docker)
```bash
cd api && go run .
# Ideal para: Desenvolvimento, debugging
```

### Docker Dev (Hot Reload)
```bash
cd docker && docker-compose -f docker-compose.dev.yml up
# Ideal para: Testes em container com mudanças ao vivo
```

### Docker Prod (Otimizado)
```bash
cd docker && docker-compose -f docker-compose.yml up -d
# Ideal para: Staging, produção
```

### Kubernetes (Escala)
```bash
kubectl apply -f docker/k8s-deployment.yaml
# Ideal para: Enterprise, alta disponibilidade
```

---

## 🎯 Endpoints Principais

```
POST   /tenants                           # Criar empresa
GET    /tenants                           # Listar empresas
DELETE /tenants/{id}                      # Deletar empresa

POST   /tenants/{tid}/sessions            # Criar sessão WhatsApp
GET    /tenants/{tid}/sessions            # Listar sessões
DELETE /tenants/{tid}/sessions/{sid}      # Deletar sessão

GET    /tenants/{tid}/sessions/{sid}/qrcode  # QR Code
POST   /tenants/{tid}/sessions/{sid}/messages # Enviar msg

GET    /health                            # Status API
```

Exemplos: Veja [api/MANUAL.md](api/MANUAL.md)

---

## 🐳 Docker: Tudo Separado

Toda configuração Docker está na pasta `docker/` para manter a raiz limpa:

```
docker/
├── Dockerfile               ← Imagem otimizada
├── docker-compose.yml       ← Produção
├── docker-compose.dev.yml   ← Desenvolvimento
├── Makefile                 ← make help
├── build.sh / build.ps1     ← Build scripts
├── deploy.sh                ← Deploy VPS
├── start.sh / start.cmd     ← Menu interativo
└── README.md                ← Docker docs
```

**Para usar:**
```bash
cd docker
make help                    # Ver todas as opções
make up                      # Iniciar
make logs-follow             # Ver logs
make down                    # Parar
```

---

## 🔒 Segurança

✅ **Implementado:**
- HTTPS/TLS 1.2+
- Non-root user em Docker
- Rate limiting
- HSTS headers
- Secrets em .env

⚠️ **Recomendado para Produção:**
- JWT authentication
- WAF (Web Application Firewall)
- DDoS protection
- Audit logging

---

## 📊 Performance

| Métrica | Local | Docker |
|---------|-------|--------|
| Boot | <1s | ~2s |
| Memory (idle) | ~50MB | ~80MB |
| Imagem | N/A | ~100MB |
| Latência | ~1ms | ~2ms |

---

## 🆘 Precisa de Ajuda?

### Problemas Comuns

**"Porta 5000 em uso"**
```bash
# Windows
netstat -ano | findstr :5000

# Linux
lsof -i :5000

# Usar outra porta em docker-compose.yml
```

**"CGO_ENABLED required"**
```powershell
$env:CGO_ENABLED="1"
go run ./api
```

**"Docker build lento"**
```bash
DOCKER_BUILDKIT=1 docker build .
```

### Documentação Relevante

- 🐛 Erro ao rodar → [QUICK-REFERENCE.md#troubleshooting](QUICK-REFERENCE.md#troubleshooting)
- 🚀 Deploy → [DEPLOYMENT-CHECKLIST.md](DEPLOYMENT-CHECKLIST.md)
- 📦 Docker → [docker/README.md](docker/README.md)
- 💻 API → [api/MANUAL.md](api/MANUAL.md)

---

## 🔄 Workflow Típico

```bash
# 1. Clonar
git clone <repo>
cd ORGANIX-API-WHATSMEOW

# 2. Testar local
cd api
go run .

# 3. Testar em Docker
cd ../docker
make build VERSION=1.0.0
make up-dev

# 4. Editar código, changes vão ser recarregadas

# 5. Deploy
make deploy HOST=user@seu-vps VERSION=v1.0.0
```

---

## 📈 Roadmap

| Versão | Status | Features |
|--------|--------|----------|
| **v1.0** | ✅ | MVP (criar tenants, sessions, QR code) |
| **v1.1** | 🔄 | Envio real de mensagens, webhooks |
| **v1.2** | 📅 | JWT auth, rate limiting, metrics |
| **v1.3** | 📅 | Redis cache, PostgreSQL, observability |
| **v2.0** | 📅 | GraphQL, WebSocket, templates, analytics |

---

## 🔗 Stack Técnico

- **Linguagem**: Go 1.21+
- **Framework**: Gorilla Mux
- **Database**: SQLite
- **WhatsApp**: go-whatsmeow
- **Container**: Docker + Compose
- **Orquestração**: Kubernetes (opcional)
- **Reverse Proxy**: Nginx
- **SSL/TLS**: Let's Encrypt
- **CI/CD**: GitHub Actions

---

## 📄 Licença

MIT License - Veja [LICENSE](LICENSE)

---

## 👥 Contribuições

Contribuições são bem-vindas! Por favor:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/minha-feature`)
3. Commit suas mudanças (`git commit -m 'feat: descrição'`)
4. Push para a branch (`git push origin feature/minha-feature`)
5. Abra um Pull Request

---

## 📞 Suporte

### Documentação
📘 Veja [README-RESUMO.md](README-RESUMO.md) para visão geral
📖 Veja [api/MANUAL.md](api/MANUAL.md) para referência API completa
🐳 Veja [docker/README.md](docker/README.md) para Docker

### Comunidade
💬 Abra uma issue para dúvidas
🐛 Reporte bugs com detalhes
📝 Crie discussões para feature requests

---

**Status**: ✅ Pronto para Uso | 🚀 Pronto para Deploy

**Last Updated**: Fevereiro 2026
