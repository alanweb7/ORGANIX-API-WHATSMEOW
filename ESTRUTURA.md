# Estrutura do Projeto - Visão Geral

## 📦 Estrutura Completa

```
ORGANIX-API-WHATSMEOW/
│
├── 📘 DOCUMENTAÇÃO (Raiz - Limpa)
│   ├── README.md                         # Este arquivo
│   ├── README-RESUMO.md                  # Visão geral executiva
│   ├── MANUAL.md (em api/)               # Manual API completo
│   ├── DOCKER.md                         # Docker detalhado
│   ├── DEPLOYMENT-CHECKLIST.md           # Deploy passo-a-passo
│   ├── ESTRUTURA.md                      # Este arquivo (visão arquitetura)
│   ├── QUICK-REFERENCE.md                # Comandos rápidos
│   └── nginx.conf.template               # Nginx pronto para produção
│
├── 🚀 API (Aplicação Go)
│   └── api/
│       ├── main.go                       # Servidor HTTP
│       ├── models.go                     # Estruturas de dados
│       ├── session.go                    # Gerenciador de sessões
│       ├── handlers.go                   # Handlers HTTP
│       ├── go.mod                        # Dependências Go
│       ├── go.sum                        # Hash das dependências
│       ├── MANUAL.md                     # Referência API completa
│       ├── .env.example                  # Variáveis de exemplo
│       └── data/                         # Dados persistidos (criado automaticamente)
│           ├── sessions/
│           │   ├── empresa1/
│           │   │   └── wa-usuario1.db
│           │   └── empresa2/
│           │       └── wa-usuario2.db
│           └── logs/
│
├── 🐳 DOCKER (Tudo isolado aqui!)
│   └── docker/
│       ├── Dockerfile                    # Multi-stage build otimizado
│       ├── docker-compose.yml            # Composição PRODUÇÃO
│       ├── docker-compose.dev.yml        # Composição DESENVOLVIMENTO (hot reload)
│       ├── .dockerignore                 # Arquivos ignorados no build
│       ├── build.sh                      # Build script (Linux/Mac)
│       ├── build.ps1                     # Build script (Windows)
│       ├── start.sh                      # Menu interativo (Linux/Mac)
│       ├── start.cmd                     # Menu interativo (Windows)
│       ├── deploy.sh                     # Deploy automático para VPS
│       ├── Makefile                      # Automação (make help)
│       ├── README.md                     # Documentação Docker
│       ├── k8s-deployment.yaml           # Manifests Kubernetes
│       └── logs/                         # Logs de containers
│
├── 🔧 CI/CD (GitHub Actions)
│   └── .github/
│       └── workflows/
│           └── build-deploy.yml          # Pipeline automático
│
├── 📚 WHATSAPP (Biblioteca)
│   └── whatsmeow/
│       ├── client.go
│       ├── send.go
│       ├── qrchan.go
│       ├── store/
│       ├── proto/
│       ├── types/
│       ├── util/
│       ├── go.mod
│       └── ... (outros arquivos)
│
└── LICENSE                               # Licença MIT
```

## 🔄 Fluxo de Desenvolvimento

```
┌──────────────────┐
│  Codificação     │
│  local (editor)  │
└────────┬─────────┘
         │
         ↓
┌──────────────────┐
│  go run .        │  ← Teste rápido
│  (sem Docker)    │
└────────┬─────────┘
         │
         ↓
┌─────────────────────────────┐
│ docker-compose.dev.yml up   │  ← Teste em container
│ (com hot reload)            │
└────────┬────────────────────┘
         │
         ↓
┌──────────────────┐
│  Commit & Push   │
│  para Git        │
└────────┬─────────┘
         │
         ↓
┌──────────────────────────┐
│  CI/CD Pipeline          │
│  (GitHub Actions, etc)   │
└────────┬─────────────────┘
         │
         ↓
┌──────────────────────────┐
│  Build Docker image      │
│  Tag & Push registry     │
└────────┬─────────────────┘
         │
         ↓
┌──────────────────────────┐
│  Deploy em VPS/K8s       │
│  docker-compose/kubectl  │
└──────────────────────────┘
```

## 🐳 Fluxo Docker

```
┌────────────────┐
│  Source Code   │
│   (local)      │
└────────┬───────┘
         │
         ↓
┌─────────────────────────┐
│  docker build           │ ← Build image
│  (Dockerfile)           │
└────────┬────────────────┘
         │
         ↓
    ┌────┴────────┬──────────┐
    │             │          │
    ↓             ↓          ↓
┌────────┐  ┌──────────┐  ┌─────────┐
│Builder │  │ Imagem   │  │ Registry│
│(Stage1)│  │ Docker   │  │ (Hub)   │
└────────┘  │ (100MB)  │  └─────────┘
            └────┬─────┘
                 │
                 ↓
        ┌────────────────┐
        │ docker-compose │
        │   ou kubectl   │
        └────────┬───────┘
                 │
                 ↓
        ┌────────────────┐
        │  Containers    │
        │  (Linux VPS)   │
        └────────────────┘
```

## 📋 Opções de Execução

```
┌─────────────────────────────────────────────────────────────┐
│                    OPÇÕES DE EXECUÇÃO                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1️⃣  SEM DOCKER (Desenvolvimento Local)                     │
│    ├─ go run .                                            │
│    ├─ .\api.exe (Windows compilado)                       │
│    └─ ./whatsapp-api (Linux compilado)                    │
│                                                             │
│ 2️⃣  COM DOCKER (Desenvolvimento)                          │
│    ├─ docker-compose -f docker-compose.dev.yml up        │
│    ├─ make up-dev                                        │
│    └─ ./build.sh && docker run whatsapp-api:dev         │
│                                                             │
│ 3️⃣  COM DOCKER (Produção)                                 │
│    ├─ docker-compose up -d                               │
│    ├─ make up                                            │
│    └─ ./deploy.sh user@vps.com v1.0.0                   │
│                                                             │
│ 4️⃣  KUBERNETES (Escala em Produção)                       │
│    ├─ kubectl apply -f k8s-deployment.yaml              │
│    ├─ kubectl get pods -n whatsapp                      │
│    └─ kubectl logs deployment/whatsapp-api              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 Casos de Uso

```
┌──────────────────────────────────────────────────────────┐
│         QUAL OPÇÃO USAR PARA CADA CASO                  │
├──────────────────────────────────────────────────────────┤
│                                                          │
│ 📍 Desenvolvimento Local                                │
│    → go run . (mais rápido)                            │
│    → ou docker-compose.dev.yml (mais realista)        │
│                                                          │
│ 🧪 Testes/QA                                            │
│    → docker-compose.dev.yml                            │
│    → ou staging com docker-compose                     │
│                                                          │
│ 🚀 Produção em VPS                                      │
│    → docker-compose up -d (simples, até ~100 req/s)   │
│    → ou Docker Swarm (> 100 req/s)                     │
│    → ou Kubernetes (escala automática, alta disp.)     │
│                                                          │
│ 📈 Escala Empresarial                                   │
│    → Kubernetes com auto-scaling                       │
│    → PersistentVolume para dados                       │
│    → LoadBalancer com Nginx/Traefik                    │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

## 🔐 Segurança em Camadas

```
┌──────────────────────────────────────────┐
│  CAMADA 1: Dockerfile                   │
│  ✅ Usuário não-root                     │
│  ✅ Base image mínima (Alpine)           │
│  ✅ Binário estático                     │
│  ✅ Sem ferramentas desnecessárias       │
└──────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────┐
│  CAMADA 2: Docker Compose                │
│  ✅ Network isolation                    │
│  ✅ Resource limits                      │
│  ✅ Restart policies                     │
│  ✅ Health checks                        │
└──────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────┐
│  CAMADA 3: Kubernetes (Opcional)         │
│  ✅ Pod security policies                │
│  ✅ Network policies                     │
│  ✅ Role-based access control            │
│  ✅ Secrets management                   │
└──────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────┐
│  CAMADA 4: VPS/Cloud                     │
│  ✅ Firewall rules                       │
│  ✅ Reverse proxy (Nginx)                │
│  ✅ SSL/TLS (Let's Encrypt)              │
│  ✅ DDoS protection                      │
└──────────────────────────────────────────┘
```

## 📊 Performance

```
┌────────────────────────────────────────┐
│  MÉTRICA          │ SEM DOCKER │ COM   │
├───────────────────┼────────────┼───────┤
│ Tempo boot        │ <1s        │ ~2s   │
│ Mem (idle)        │ ~50MB      │ ~80MB │
│ CPU (idle)        │ ~5%        │ ~2%   │
│ Latência API      │ ~1ms       │ ~2ms  │
│ Throughput        │ ~1000 req/s│ ~900  │
│ Portabilidade     │ ❌         │ ✅    │
│ Deploy fácil      │ ❌         │ ✅    │
│ Reproducibilidade│ ❌         │ ✅    │
└────────────────────────────────────────┘
```

## 🎓 Recomendações

```
┌─────────────────────────────────────────────────┐
│  FASE DO PROJETO      │  RECOMENDAÇÃO            │
├──────────────────────┼──────────────────────────┤
│ MVP / Prototipagem    │ go run . (local)         │
│ Testes Iniciais       │ docker-compose.dev.yml   │
│ Staging               │ docker-compose.yml       │
│ Produção Beta         │ Docker Compose (1 VPS)   │
│ Produção (Escala)     │ Kubernetes / Cloud       │
│ Enterprise            │ K8s + observability      │
└─────────────────────────────────────────────────┘
```

## 📚 Próximo Passo

Leia [MANUAL.md](MANUAL.md) para:
- Setup completo
- Ejemplos de uso
- Deployment em produção
- Troubleshooting

Ou [DOCKER.md](DOCKER.md) para:
- Docker setup detalhado
- CI/CD integration
- Kubernetes
- Performance tuning
