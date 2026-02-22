# 📋 Checklist de Deployment & Boas Práticas

## ✅ Pré-Deployment

### Validação Local
- [ ] `go run .` funciona localmente (port 5000)
- [ ] Testado endpoints com curl/Postman
- [ ] Verificado logs sem erros
- [ ] Base de dados criada e operacional
- [ ] Todas as dependências instaladas

### Validação Docker Local
```bash
# Windows PowerShell
.\build.ps1 -Version "1.0.0"
docker-compose up -d
curl http://localhost:5000/health

# Linux/Mac
./build.sh 1.0.0
docker-compose up -d
curl http://localhost:5000/health
```

- [ ] Docker build completa com sucesso
- [ ] Imagem criada (~100MB)
- [ ] Container inicia sem erros
- [ ] Health endpoint responde
- [ ] Volumes montados corretamente
- [ ] Logs acessíveis via `docker logs`

### Validação Segurança
- [ ] .env contém variáveis sensitive (não commitado)
- [ ] .dockerignore reduz build context
- [ ] Dockerfile usa usuário não-root
- [ ] Senhas não estão hardcoded no código
- [ ] Secrets em .env ou docker-compose

---

## 🚀 Deploy em VPS (Linux)

### 1️⃣ Preparação do VPS

```bash
# SSH para VPS
ssh user@seu-vps.com

# Criar diretório do projeto
mkdir -p /opt/whatsapp-api
cd /opt/whatsapp-api

# Instalar Docker e Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verificar instalação
docker --version
docker-compose --version
```

**Checklist:**
- [ ] SSH acesso funcionando
- [ ] Docker instalado
- [ ] Docker Compose instalado
- [ ] Diretório `/opt/whatsapp-api` criado
- [ ] Usuário com permissão Docker

### 2️⃣ Deploy Automático

**Opção A: Script Automático (Recomendado)**

```bash
# Do seu PC local
./deploy.sh seu-usuario seu-vps.com v1.0.0

# Script fará:
# ✅ SSH para VPS
# ✅ Clonar/atualizar repo
# ✅ Build imagem Docker
# ✅ Start containers com docker-compose
# ✅ Verificar health
```

**Opção B: Manual**

```bash
# No VPS
cd /opt/whatsapp-api

# Criar .env
cat > .env << 'EOF'
ENVIRONMENT=production
API_PORT=5000
LOG_LEVEL=info
STORAGE_PATH=/data/sessions
DATABASE_TYPE=sqlite
EOF

# Carregar docker-compose.yml do repo
wget https://seu-repo/docker-compose.yml

# Iniciar
docker-compose up -d

# Verificar
docker-compose ps
docker logs whatsapp-api
```

**Checklist:**
- [ ] Arquivo `.env` criado com valores corretos
- [ ] `docker-compose.yml` no VPS
- [ ] Containers iniciados com `docker-compose up -d`
- [ ] Health check passando: `curl http://localhost:5000/health`
- [ ] Logs verificados: `docker logs whatsapp-api`

### 3️⃣ Nginx Reverse Proxy

**Instalar Nginx:**

```bash
sudo apt update
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

**Configurar:**

```bash
sudo nano /etc/nginx/sites-available/whatsapp-api
```

Conteúdo:

```nginx
upstream whatsapp_backend {
    server localhost:5000;
}

server {
    listen 80;
    server_name seu-dominio.com www.seu-dominio.com;

    # Redirecionar HTTP para HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name seu-dominio.com www.seu-dominio.com;

    ssl_certificate /etc/letsencrypt/live/seu-dominio.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/seu-dominio.com/privkey.pem;

    # SSL Moderno
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    client_max_body_size 100M;

    location / {
        proxy_pass http://whatsapp_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check (sem log)
    location /health {
        access_log off;
        proxy_pass http://whatsapp_backend;
    }
}
```

Habilitar:

```bash
sudo ln -s /etc/nginx/sites-available/whatsapp-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

**Certificado SSL:**

```bash
sudo apt install certbot python3-certbot-nginx -y
sudo certbot certonly --nginx -d seu-dominio.com -d www.seu-dominio.com
```

**Checklist:**
- [ ] Nginx instalado
- [ ] Config criada e testada
- [ ] Re-enabled e recarregado
- [ ] Certificado SSL obtido
- [ ] HTTPS funcionando
- [ ] Redirecionamento HTTP→HTTPS

---

## ☸ Deploy em Kubernetes

### Pré-requisito
- [ ] Kubernetes cluster pronto (EKS, GKE, AKS ou K3s)
- [ ] `kubectl` configurado
- [ ] Namespace criado

### Deploy

```bash
# Criar namespace
kubectl create namespace whatsapp

# Aplicar manifests
kubectl apply -f k8s-deployment.yaml -n whatsapp

# Verificar
kubectl get all -n whatsapp
kubectl logs -f deployment/whatsapp-api -n whatsapp

# Expor serviço
kubectl port-forward svc/whatsapp-api 5000:5000 -n whatsapp
```

**Checklist:**
- [ ] Cluster Kubernetes operacional
- [ ] Namespace whatsapp criado
- [ ] ConfigMap aplicado
- [ ] PersistentVolume criado
- [ ] Deployment rodando (3 replicas)
- [ ] Service exposto
- [ ] HPA configurado (auto-scaling)

### Expandir Replicas

```bash
# Manual
kubectl scale deployment whatsapp-api --replicas=5 -n whatsapp

# Automático (já configurado no manifest)
# HPA escala automaticamente entre 2-10 podas
# baseado em CPU (70%) e Memória (80%)

# Monitorar scaling
kubectl get hpa -n whatsapp -w
```

---

## 📊 Monitoramento & Logging

### Docker (Local)

```bash
# Ver logs em tempo real
docker logs -f whatsapp-api

# Ver containers
docker ps
docker stats

# Ver volumes
docker volume ls
```

### Kubernetes

```bash
# Logs de um pod
kubectl logs -f pod/whatsapp-api-xxxxx -n whatsapp

# Logs de toda deployment
kubectl logs -f deployment/whatsapp-api -n whatsapp

# Descrever erro
kubectl describe pod whatsapp-api-xxxxx -n whatsapp

# Ver eventos
kubectl get events -n whatsapp

# Dashboard (se disponível)
kubectl proxy
# Acesse: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/
```

### Métricas de Health

```bash
# API health
curl http://seu-dominio.com/health

# Resposta esperada:
# {
#   "status": "ok",
#   "timestamp": "2024-01-15T10:30:00Z"
# }

# Server metrics (implementar em v1.1)
# GET /metrics (Prometheus format)
```

---

## 🔍 Troubleshooting

### Container não inicia

```bash
# Ver erro
docker logs whatsapp-api

# Comuns:
# 1. Porta já em uso
#    → Mudar porta em .env ou usar netstat para liberar

# 2. Permissão de arquivo
#    → sudo chown -R 1000:1000 ./data

# 3. Arquivo corrupto
#    → Deletar data/ e recomeçar
```

### API lenta

```bash
# Verificar recursos
docker stats whatsapp-api

# Se memória alta:
# 1. Aumentar limite em docker-compose.yml
# 2. Limpar logs antigos
# 3. Escalar com múltiplas instâncias

# Se CPU alta:
# 1. Implementar cache
# 2. Otimizar queries
# 3. Escalar horizontalmente
```

### Imagem muito grande

```bash
# Verificar tamanho
docker images whatsapp-api

# Otimizações:
# 1. Multi-stage build (já implementado)
# 2. -ldflags="-s -w" (strip binário)
# 3. Alpine instead of Ubuntu

# Esperado: ~100MB (otimizado)
```

---

## 🔐 Segurança em Produção

### Firewall

```bash
# Abrir apenas portas necessárias
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
sudo ufw status
```

### Backup

```bash
# Backup diário do banco de dados
0 2 * * * docker exec whatsapp-api tar -czf /backup/db-$(date +\%Y\%m\%d).tar.gz /data

# Backup mensal para cloud
0 3 1 * * aws s3 sync /backup s3://seu-bucket/whatsapp-api/
```

### Atualizações

```bash
# Verificar atualizações
docker pull seu-registry/whatsapp-api:latest

# Deploy com zero-downtime
docker-compose up -d --no-deps --build whatsapp-api
```

---

## 📈 Performance Tuning

### Docker Compose

```yaml
services:
  whatsapp-api:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 512M
        reservations:
          cpus: '1'
          memory: 256M
```

### Kernel Parameters (VPS)

```bash
# Aumentar file descriptors
echo "fs.file-max = 2097152" | sudo tee -a /etc/sysctl.conf
echo "net.core.somaxconn = 32768" | sudo tee -a /etc/sysctl.conf
echo "net.ipv4.tcp_max_syn_backlog = 16384" | sudo tee -a /etc/sysctl.conf

sudo sysctl -p
```

### Nginx Tuning

```nginx
# Em /etc/nginx/nginx.conf

user www-data;
worker_processes auto;  # Auto-detect CPU cores
worker_connections 2048;

keepalive_timeout 65;
gzip on;
gzip_types text/plain application/json;
```

---

## 📝 Checklist Final

Antes de considerar "pronto para produção":

**Código:**
- [ ] Sem console.log ou print statements
- [ ] Tratamento de erros em todos endpoints
- [ ] Validação de input em requests
- [ ] Rate limiting implementado
- [ ] Testes unitários passando

**Docker:**
- [ ] Imagem otimizada (<150MB)
- [ ] Health checks funcionando
- [ ] Environment variables documentadas
- [ ] Dockerfile bem comentado
- [ ] .dockerignore completo

**Segurança:**
- [ ] SSL/TLS habilitado
- [ ] Firewall configurado
- [ ] Secrets não expostos
- [ ] CORS restrito
- [ ] API key/JWT implementado

**Operações:**
- [ ] Backup strategy definida
- [ ] Monitoramento ativo
- [ ] Logging centralizado
- [ ] Runbooks para operações
- [ ] Alertas configurados

**Documentação:**
- [ ] README atualizado
- [ ] API documentation completa
- [ ] Troubleshooting guide
- [ ] Deployment procedure
- [ ] Rollback procedure

**Testes:**
- [ ] Load test (> 1000 req/min)
- [ ] Failover test
- [ ] Recovery test
- [ ] Segurança test (OWASP top 10)
- [ ] Performance test (latência <500ms)

---

## 🎯 Próximas Fases

### v1.1 - Observability
- Prometheus metrics
- Datadog/NewRelic integration
- Distributed tracing
- Custom dashboards

### v1.2 - Resiliência
- Circuit breaker pattern
- Retry with exponential backoff
- Bulkhead isolation
- Graceful degradation

### v1.3 - Scale
- Redis cache
- PostgreSQL (multi-instance)
- Message queue (RabbitMQ/Kafka)
- Read replicas

---

## 📞 Support & Escalation

Se tiver problemas:

1. **Local Issues**: `go run .` debug
2. **Docker Issues**: `docker logs` análise
3. **Deployment Issues**: Check `.env` e permissões
4. **Performance Issues**: Monitor recursos (`docker stats`)
5. **Security Issues**: Revisar `Dockerfile` e `nginx.conf`

Para maiores informações: Veja [MANUAL.md](MANUAL.md) e [DOCKER.md](DOCKER.md)
