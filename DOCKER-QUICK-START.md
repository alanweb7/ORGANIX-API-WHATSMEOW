# 🎯 Comandos Rápidos - Docker Reorganizado

## 📍 Inicio Rápido (Copie & Cole)

### Windows PowerShell

```powershell
# Entrar na pasta Docker
cd e:\apps\CLIENTES\ORGANIX-API-WHATSMEOW\docker

# Build
.\build.ps1 -Version "1.0.0"

# Iniciar (escolha uma)
docker-compose -f docker-compose.yml up -d      # Produção
docker-compose -f docker-compose.dev.yml up -d  # Desenvolvimento

# Ou use Menu Interativo
.\start.cmd

# Verificar
curl http://localhost:5000/health

# Ver logs
docker-compose logs -f

# Parar
docker-compose down
```

### Linux/Mac Bash

```bash
# Entrar na pasta Docker
cd ORGANIX-API-WHATSMEOW/docker

# Build
chmod +x build.sh
./build.sh 1.0.0

# Iniciar (escolha uma)
docker-compose -f docker-compose.yml up -d      # Produção
docker-compose -f docker-compose.dev.yml up -d  # Desenvolvimento

# Ou use Menu Interativo
chmod +x start.sh
./start.sh

# Verificar
curl http://localhost:5000/health

# Ver logs
docker-compose -f docker-compose.yml logs -f

# Parar
docker-compose down
```

## 🛠 Usando Make (Recomendado)

```bash
cd docker

# Ver todos os comandos
make help

# Build
make build VERSION=1.0.0
make build-dev

# Iniciar
make up              # Produção
make up-dev          # Desenvolvimento

# Monitoramento
make logs            # Ver logs
make logs-follow     # Tail logs
make health          # Verificar saúde
make stats           # Ver CPU/Memory

# Parar
make down

# Deploy
make deploy HOST=ubuntu@seu-vps.com VERSION=v1.0.0
```

## 🎪 Menu Interativo

### Windows
```powershell
cd docker
.\start.cmd
# Digite opção (1-6)
```

### Linux/Mac
```bash
cd docker
./start.sh
# Digite opção (1-6)
```

## 🗂 Estrutura (Lembrança)

```
ORGANIX-API-WHATSMEOW/
├── api/               ← Código Go (sem Docker)
│   └── go run .
├── docker/            ← ⭐ NOVO! Tudo Docker aqui
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   ├── Makefile
│   ├── build.sh / build.ps1
│   ├── start.sh / start.cmd
│   ├── deploy.sh
│   └── README.md
├── README.md          ← Leia primeiro
├── DOCKER-REORGANIZADO.md ← Mudanças
└── ... (outras docs)
```

## 📋 Checklist: Está Tudo Funcional?

- [ ] `docker/Dockerfile` existe ✅
- [ ] `docker/docker-compose.yml` existe ✅
- [ ] `docker/Makefile` existe ✅
- [ ] `docker/build.sh` ou `build.ps1` existe ✅
- [ ] `docker/start.sh` ou `start.cmd` existe ✅
- [ ] Código em `api/main.go` existe ✅
- [ ] Nenhuma funcionalidade quebrada ✅

Se tudo acima está ✅, **pronto para usar!**

## 🚀 Testando

### 1️⃣ Rápido (Local)
```bash
cd api
go run .
# Abra outro terminal
curl http://localhost:5000/health
```

### 2️⃣ Com Docker (Recomendado)
```bash
cd docker
make build VERSION=1.0.0
make up
make health
```

### 3️⃣ Desenvolvimento (Hot Reload)
```bash
cd docker
docker-compose -f docker-compose.dev.yml up -d
# Edite arquivos em api/ e veja mudanças ao vivo
docker-compose logs -f
```

---

**Está tudo pronto!** 🎉

Para documentação completa:
- [docker/README.md](docker/README.md) - Docker específico
- [QUICK-REFERENCE.md](QUICK-REFERENCE.md) - Tutti comandos
- [README.md](README.md) - Visão geral
