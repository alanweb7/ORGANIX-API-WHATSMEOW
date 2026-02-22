# ⚡ Quick Reference - Comando Rápidos

## 🚀 Iniciar Rápido (Windows)

```powershell
# Terminal PowerShell na pasta do projeto

# 1. Desenvolvimento local (mais rápido)
$env:CGO_ENABLED="1"
go run .\api\main.go

# 2. Com Docker (mais realista)
.\build.ps1 -Version "1.0.0"
docker-compose up -d

# 3. Testar health
curl http://localhost:5000/health

# 4. Ver logs
docker logs whatsapp-api
```

---

## 🐧 Iniciar Rápido (Linux/Mac)

```bash
# Terminal bash na pasta do projeto

# 1. Desenvolvimento local
go run ./api/main.go

# 2. Com Docker
./build.sh 1.0.0
docker-compose up -d

# 3. Testar health
curl http://localhost:5000/health

# 4. Ver logs
docker logs whatsapp-api
```

---

## 📋 Comandos Make (Recomendado)

```bash
# Ajuda
make help

# Build Docker
make build VERSION=1.0.0
make build-dev

# Rodar
make up              # Produção
make up-dev          # Desenvolvimento
make down            # Parar containers

# Logs
make logs
make logs-follow

# Testes
make test
make lint

# Deploy
make deploy HOST=seu-vps.com USER=ubuntu VERSION=1.0.0
make deploy-k8s NAMESPACE=whatsapp

# Limpeza
make clean
make clean-all
```

---

## 🐳 Docker Compose Direto

```bash
# Produção
docker-compose up -d                    # Iniciar
docker-compose down                     # Parar
docker-compose logs -f                  # Logs em tempo real
docker-compose ps                       # Ver containers
docker-compose stats                    # Ver recursos

# Desenvolvimento
docker-compose -f docker-compose.dev.yml up
docker-compose -f docker-compose.dev.yml down

# Rebuild (se mudou código)
docker-compose up -d --build
docker-compose up -d --no-deps --build whatsapp-api
```

---

## 🧪 Testar API

### Health Check

```bash
curl http://localhost:5000/health

# Resposta esperada:
# {"status":"ok","timestamp":"2024-01-15T10:00:00Z"}
```

### Criar Tenant

```bash
curl -X POST http://localhost:5000/tenants \
  -H "Content-Type: application/json" \
  -d '{"name":"Minha Empresa","email":"admin@empresa.com"}'
```

### Listar Tenants

```bash
curl http://localhost:5000/tenants
```

### Criar Sessão

```bash
# Substituir TENANT_ID
curl -X POST http://localhost:5000/tenants/TENANT_ID/sessions \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"5511999999999"}'
```

### Obter QR Code

```bash
# Substituir TENANT_ID e SESSION_ID
curl http://localhost:5000/tenants/TENANT_ID/sessions/SESSION_ID/qrcode
```

### Deletar Session

```bash
curl -X DELETE http://localhost:5000/tenants/TENANT_ID/sessions/SESSION_ID
```

---

## 🔨 Build Manual

### Windows PowerShell

```powershell
# Setup
$env:CGO_ENABLED="1"
$env:GOOS="linux"
$env:GOARCH="amd64"

# Build para Linux
go build -o whatsapp-api-linux ./api

# Build para Windows
$env:GOOS="windows"
go build -o whatsapp-api.exe ./api

# Build com versão
go build -ldflags="-X main.Version=1.0.0 -s -w" -o whatsapp-api ./api
```

### Linux/Mac

```bash
# Build local
go build -o whatsapp-api ./api

# Build para Linux (de Mac/Windows)
GOOS=linux GOARCH=amd64 go build -o whatsapp-api-linux ./api

# Build com versão
go build -ldflags="-X main.Version=1.0.0 -s -w" -o whatsapp-api ./api
```

---

## 🐳 Docker Manual

```bash
# Build imagem
docker build -t whatsapp-api:1.0.0 ./api

# Tag para registry
docker tag whatsapp-api:1.0.0 seu-registry/whatsapp-api:1.0.0

# Push para registry
docker push seu-registry/whatsapp-api:1.0.0

# Rodar container
docker run -d \
  --name whatsapp-api \
  -p 5000:5000 \
  -e ENVIRONMENT=production \
  -v ./data:/data \
  whatsapp-api:1.0.0

# Ver logs
docker logs -f whatsapp-api

# Exec shell no container
docker exec -it whatsapp-api /bin/sh

# Ver recursos
docker stats whatsapp-api

# Parar container
docker stop whatsapp-api
docker rm whatsapp-api
```

---

## ☸ Kubernetes

```bash
# Aplicar manifests
kubectl apply -f k8s-deployment.yaml -n whatsapp

# Verificar deployment
kubectl get deployment -n whatsapp
kubectl get pods -n whatsapp
kubectl get svc -n whatsapp

# Ver logs de um pod
kubectl logs POD_NAME -n whatsapp
kubectl logs -f deployment/whatsapp-api -n whatsapp

# Descrever erro
kubectl describe pod POD_NAME -n whatsapp

# Escalar
kubectl scale deployment whatsapp-api --replicas=5 -n whatsapp

# Port forward (acessar localmente)
kubectl port-forward svc/whatsapp-api 5000:5000 -n whatsapp

# Deletar
kubectl delete deployment whatsapp-api -n whatsapp
```

---

## 🚀 Deploy VPS

### Automático (Recomendado)

```bash
./deploy.sh seu-usuario seu-vps.com v1.0.0

# Ou se tiver Make
make deploy HOST=seu-vps.com USER=seu-usuario VERSION=v1.0.0
```

### Semi-Automático

```bash
# 1. SSH para VPS
ssh seu-usuario@seu-vps.com

# 2. Entrar na pasta
cd /opt/whatsapp-api

# 3. Atualizar código
git pull origin main

# 4. Deploy
docker-compose pull
docker-compose up -d

# 5. Verificar
docker logs whatsapp-api
curl http://localhost:5000/health
```

### Manual Passo-a-Passo

```bash
# 1. Setup VPS (primeira vez)
ssh sua-usuario@seu-vps.com
mkdir -p /opt/whatsapp-api
cd /opt/whatsapp-api

# 2. Clonar repo (ou copiar arquivos)
git clone https://github.com/seu-usuario/organix-api.git .

# 3. Setup variáveis
cp .env.example .env
nano .env  # Editar variáveis

# 4. Iniciar containers
docker-compose up -d

# 5. Verificar status
docker-compose ps
docker logs whatsapp-api

# 6. Testar
curl http://localhost:5000/health

# 7. Setup Nginx (se não feito)
sudo cp nginx.conf.template /etc/nginx/sites-available/whatsapp-api
sudo ln -s /etc/nginx/sites-available/whatsapp-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# 8. Certificado SSL (primeira vez)
sudo certbot certonly --nginx -d seu-dominio.com

# 9. Acessar API
curl https://seu-dominio.com/health
```

---

## 🔍 Troubleshooting

### API não responde

```bash
# Verificar se está rodando
docker ps | grep whatsapp

# Ver logs de erro
docker logs whatsapp-api

# Testar porta localmente
netstat -ano | findstr :5000    # Windows
ss -tlnp | grep :5000           # Linux

# Testar health
curl -v http://localhost:5000/health
```

### Docker build falha

```bash
# Verificar Dockerfile syntaxe
docker build --no-cache -t teste ./api

# Ver layer por layer
docker build --progress=plain ./api

# Limpar tudo e recomeçar
docker system prune -a
docker build -t whatsapp-api:1.0.0 ./api
```

### Container crashes

```bash
# Ver logs detalhados
docker logs whatsapp-api -n 100  # últimas 100 linhas

# Restartar container
docker restart whatsapp-api

# Verificar recursos (OOM?)
docker stats whatsapp-api

# Aumentar memória em docker-compose.yml
#   deploy:
#     resources:
#       limits:
#         memory: 512M
```

### Erro de permissão

```bash
# Verificar ownership
ls -la ./data/

# Corrigir
sudo chown -R 1000:1000 ./data/

# Ou em Docker
docker exec whatsapp-api chown -R appuser:appuser /data
```

### Nginx redireciona errado

```bash
# Editar config
sudo nano /etc/nginx/sites-available/whatsapp-api

# Sintaxe correta?
sudo nginx -t

# Recarregar
sudo systemctl reload nginx

# Ver logs
sudo tail -f /var/log/nginx/error.log
```

---

## 📊 Monitoramento

```bash
# CPU/Memória
docker stats whatsapp-api
docker stats --no-stream

# Logs em tempo real
docker logs -f whatsapp-api

# Últimas N linhas
docker logs whatsapp-api -n 50

# Buscar erro específico
docker logs whatsapp-api 2>&1 | grep "error"

# Salvar logs para arquivo
docker logs whatsapp-api > logs.txt 2>&1

# Com timestamp
docker logs whatsapp-api -t

# Ver eventos Docker
docker events --filter 'container=whatsapp-api'
```

---

## 🧹 Limpeza

```bash
# Parar tudo
docker-compose down

# Remover containers
docker rm whatsapp-api

# Remover imagem
docker rmi whatsapp-api:1.0.0

# Remover volumes (cuidado: deleta dados!)
docker volume rm whatsapp-api_data

# Limpar tudo não usado
docker system prune -a

# Limpar apenas volumes
docker volume prune

# Limpeza profunda (recomeçar do zero)
docker system prune -a --volumes
```

---

## 📝 Git & CI/CD

```bash
# Clone & setup
git clone <repo-url>
cd organix-api
git checkout develop

# Criar feature branch
git checkout -b feature/minha-feature

# Commit & push
git add .
git commit -m "feat: descrição da feature"
git push origin feature/minha-feature

# Criar PR (GitHub)
# 1. Push de novo
# 2. Ir para GitHub
# 3. Clicar "Create Pull Request"
# 4. GitHub Actions roda automático (lint, test, build)

# Merge (após aprovação)
git checkout main
git pull
git merge feature/minha-feature

# Tag & Release
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0
# GitHub Actions: build + deploy + release notes automático

# Atualizar local
git fetch --all --tags
git pull origin main
```

---

## ⚙️ Environment Variables

Editar `.env`:

```bash
# Basic
ENVIRONMENT=production
API_PORT=5000
LOG_LEVEL=info

# Storage
STORAGE_PATH=/data/sessions
DATABASE_TYPE=sqlite

# Security
AUTH_ENABLED=false
JWT_SECRET=sua-chave-secreta-aqui

# Rate Limiting
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=60

# Timezone
TZ=America/Sao_Paulo
```

---

## 🆘 Comandos de Emergência

```bash
# Parar containers ainda em execução
docker-compose kill

# Forçar remover container
docker rm -f whatsapp-api

# Resetar volumes (CUIDADO: deleta dados!)
docker-compose down -v

# Recriar do zero
docker-compose down
docker-compose up --build -d

# Ver o que vai ser deletado
docker system df

# Parar Docker completamente
sudo systemctl stop docker
sudo systemctl start docker
```

---

## 📞 Referências Rápidas

- **API em funcionamento**: http://localhost:5000
- **Health Check**: http://localhost:5000/health
- **Documentação**: [MANUAL.md](MANUAL.md)
- **Docker Detalhado**: [DOCKER.md](DOCKER.md)
- **Deployment**: [DEPLOYMENT-CHECKLIST.md](DEPLOYMENT-CHECKLIST.md)
- **Estrutura**: [ESTRUTURA.md](ESTRUTURA.md)

---

## 🎯 Fluxo Típico de Desenvolvimento

```bash
# 1. Atualizar código local
git pull origin develop

# 2. Criar feature
git checkout -b feature/nova-funcionalidade

# 3. Testar local
go run ./api/main.go
# Abrir outro terminal
curl http://localhost:5000/health

# 4. Fazer commit
git add .
git commit -m "feat: minha nova feature"
git push origin feature/nova-funcionalidade

# 5. GitHub Actions roda testes e build

# 6. PR → Code Review → Merge

# 7. Deploy automático (se merge para develop)

# 8. Pro produção:
git tag -a v1.1.0 -m "Release 1.1.0"
git push origin v1.1.0
# GitHub Actions: build + deploy + release
```

---

## 💡 Pro Tips

✅ **Sempre testar local antes de push**
```bash
go test ./...
go lint ./...
```

✅ **Use docker-compose.dev.yml para desenvolvimento**
```bash
docker-compose -f docker-compose.dev.yml up
# Muda arquivos = rebuild automático
```

✅ **Salve logs em arquivo para debug**
```bash
docker logs whatsapp-api > debug.log 2>&1
grep "error" debug.log
```

✅ **Use make para comandos comuns**
```bash
make help  # Ver todos os comandos
```

✅ **Sempre use .env para variáveis sensitivas**
```bash
# Nunca commit de .env!
# Usar .env.example como template
```

---

**Last Updated:** 2024-01-15
**Status:** ✅ Pronto para Uso
