# 🚀 Instalação em VPS Linux Ubuntu

Guia completo para instalar a **WhatsApp API** em um servidor VPS com Ubuntu/Debian.

---

## 📋 Pré-requisitos

- **VPS Linux**: Ubuntu 20.04, 22.04 ou Debian 11+
- **Acesso SSH**: Com permissões sudo
- **RAM mínima**: 512MB (recomendado 1GB+)
- **Espaço em disco**: 200MB mínimo

---

## ✅ Opção 1: Deploy Automatizado (Recomendado)

Use o script pronto para fazer tudo automaticamente.

### Passo 1: Preparar a máquina local

No seu **computador/máquina local** (Windows/Mac/Linux):

```bash
cd e:\apps\CLIENTES\ORGANIX-API-WHATSMEOW\docker
bash deploy.sh usuario@seu-vps.com v1.0.0
```

### Exemplos reais:

```bash
# Com IP direto
bash deploy.sh root@192.168.1.100 v1.0.0

# Com domínio
bash deploy.sh ubuntu@api.meuservidor.com.br v1.0.0

# Com porta SSH customizada
# (edite o script para: ssh -p 2222 -o ConnectTimeout=10)
```

### O que o script faz automaticamente:

✅ Conecta via SSH  
✅ Cria diretório `/opt/whatsapp-api`  
✅ Faz git pull (se repositório)  
✅ Faz pull da imagem Docker  
✅ Inicia containers  
✅ Valida Health Check  

---

## 📝 Opção 2: Instalação Manual Completa

Para maior controle e entendimento do processo.

### Passo 1: Conectar ao VPS

```bash
ssh usuario@seu-vps.com
# ou
ssh -i /caminho/para/chave.pem ubuntu@seu-vps.com
```

### Passo 2: Verificar e atualizar o sistema

```bash
# Atualizar repositórios
sudo apt update

# Upgrade de pacotes (opcional)
sudo apt upgrade -y

# Verificar versão Ubuntu
lsb_release -a
```

### Passo 3: Instalar Docker e Docker Compose

#### Para Ubuntu/Debian:

```bash
# Instalar Docker
sudo apt install -y docker.io docker-compose git curl

# Adicionar usuário ao grupo docker (sem sudo)
sudo usermod -aG docker $USER

# Aplicar grupo novo (necessário logout/login ou:)
newgrp docker

# Verificar instalação
docker --version
docker-compose --version
```

#### Para CentOS/RHEL (alternativa):

```bash
sudo yum install -y docker docker-compose git curl
sudo systemctl start docker
sudo systemctl enable docker
```

### Passo 4: Preparar diretório de aplicação

```bash
# Criar diretório
sudo mkdir -p /opt/whatsapp-api
cd /opt/whatsapp-api

# Ajustar permissões (se necessário)
sudo chown $USER:$USER /opt/whatsapp-api

# Criar volume Docker para persistência
docker volume create whatsapp-data
```

### Passo 5: Obter o código

#### Opção A: Via Git (se tiver repositório)

```bash
cd /opt/whatsapp-api
git clone https://seu-repo-privado.git .
# ou se já estiver em um repo
git pull origin main
```

#### Opção B: Copiar arquivos manualmente

Via SCP (da sua máquina local):

```bash
# No seu computador:
scp -r e:\apps\CLIENTES\ORGANIX-API-WHATSMEOW\* usuario@seu-vps.com:/opt/whatsapp-api/
```

### Passo 6: Configurar variáveis de ambiente

```bash
cd /opt/whatsapp-api

# Copiar arquivo de exemplo
cp api/.env.example .env

# Editar o arquivo
nano .env
# (Use Ctrl+X para salvar, Y para confirmar)
```

**Variáveis importantes:**

```env
# Porta da API
API_PORT=5000

# Base de dados
DB_PATH=/opt/whatsapp-api/data/whatsapp-api.db

# Sessions
SESSIONS_DIR=/opt/whatsapp-api/data/sessions

# WhatsApp
WHATSAPP_TIMEOUT=30

# Log
LOG_LEVEL=info
```

### Passo 7: Iniciar containers

```bash
cd /opt/whatsapp-api

# Criar arquivo docker-compose (se não tiver)
# ou use:
docker-compose -f docker/docker-compose.yml up -d

# Verificar status
docker-compose ps

# Ver logs
docker-compose logs -f whatsapp-api
```

### Passo 8: Testar a API

```bash
# No VPS (teste local)
curl http://localhost:5000/health

# Esperado:
# {"status":"ok"}

# De fora do VPS (via IP ou domínio)
curl http://seu-vps.com:5000/health
curl http://192.168.1.100:5000/health
```

---

## 🛠️ Opção 3: Usando Makefile (Mais simples)

Se copiar os arquivos do `/docker`:

```bash
cd /opt/whatsapp-api/docker

# Ver todos os comandos disponíveis
make help

# Build da imagem
make build VERSION=1.0.0

# Iniciar containers
make up

# Parar containers
make down

# Ver logs
make logs

# Ver status
make ps
```

---

## 🌐 Configurar Acesso Remoto

### Opção 1: Com IP direto

A API estará acessível em:
```
http://seu-ip-vps:5000/health
http://seu-ip-vps:5000/api/tenant/{tenantId}/sessions
```

### Opção 2: Com domínio + Nginx Reverse Proxy

```bash
# Instalar Nginx
sudo apt install -y nginx

# Criar arquivo de configuração
sudo nano /etc/nginx/sites-available/whatsapp-api
```

Conteúdo:

```nginx
server {
    listen 80;
    server_name seu-dominio.com api.seu-dominio.com;

    location / {
        proxy_pass http://localhost:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Ativar:

```bash
sudo ln -s /etc/nginx/sites-available/whatsapp-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

Acesso via domínio:
```
http://seu-dominio.com/health
http://seu-dominio.com/api/tenant/{tenantId}/sessions
```

### Opção 3: SSL/HTTPS com Certbot

```bash
# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obter certificado (usando Nginx)
sudo certbot --nginx -d seu-dominio.com

# Auto-renovação já vem ativada
```

Pronto! Acesso seguro:
```
https://seu-dominio.com/health
```

---

## 📊 Monitoramento e Manutenção

### Ver logs da API

```bash
# Logs em tempo real
docker-compose logs -f whatsapp-api

# Últimas 100 linhas
docker-compose logs --tail=100 whatsapp-api

# Filtrar por data/hora
docker-compose logs --since 2024-01-20 whatsapp-api
```

### Reiniciar containers

```bash
# Reiniciar apenas API
docker-compose restart whatsapp-api

# Reiniciar todos
docker-compose restart

# Atualizar imagem e reiniciar
docker-compose pull
docker-compose up -d
```

### Verificar recursos (CPU, memória)

```bash
# Status dos containers
docker-compose ps

# Uso de recursos
docker stats

# Detalhes da imagem
docker inspect whatsapp-api
```

### Fazer backup de dados

```bash
# Copiar dados para máquina local
scp -r usuario@seu-vps.com:/opt/whatsapp-api/data ./backup-$(date +%Y%m%d)

# Ou usando rsync (mais eficiente)
rsync -avz usuario@seu-vps.com:/opt/whatsapp-api/data/ ./data-backup/
```

---

## ⚠️ Troubleshooting

### Erro: "Connection refused"

```bash
# Verificar se containers estão rodando
docker-compose ps

# Se não estão, iniciar
docker-compose up -d

# Ver erros
docker-compose logs whatsapp-api
```

### Erro: "Permission denied"

```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER
newgrp docker

# Ou executar com sudo
sudo docker-compose up -d
```

### Porta 5000 já em uso

```bash
# Ver quem está usando
sudo netstat -tulpn | grep 5000

# Alterar porta no docker-compose.yml
# De: "5000:5000"
# Para: "8080:5000" (externa:interna)
```

### Erro: "Your kernel does not support cgroup memory limit"

Geralmente em VPS compartilhadas. Contato com provedor é necessário ou ignorar o aviso.

### Nenhuma conexão com banco de dados

```bash
# Verificar permissões da pasta data
ls -la /opt/whatsapp-api/data/

# Corrigir se necessário
chmod 755 /opt/whatsapp-api/data
chmod 644 /opt/whatsapp-api/data/*.db
```

---

## 🔒 Segurança

### Firewall (UFW)

```bash
# Ativar firewall
sudo ufw enable

# Permitir SSH
sudo ufw allow 22

# Permitir HTTP
sudo ufw allow 80

# Permitir HTTPS
sudo ufw allow 443

# Bloquear porta 5000 (usar só via reverse proxy)
sudo ufw deny 5000

# Ver status
sudo ufw status
```

### Proteger banco de dados

```bash
# Usar volume Docker com permissões restritas
chmod 700 /opt/whatsapp-api/data

# Backups criptografados
gpg --symmetric /opt/whatsapp-api/data/whatsapp-api.db
```

### Limpar imagens antigas

```bash
# Ver imagens
docker images

# Remover imagens não utilizadas
docker image prune -a
```

---

## 📈 Próximas Etapas

1. **Registre o domínio** (se usar domínio)
2. **Configure DNS** apontando para o VPS
3. **Obtenha certificado SSL** (letsencrypt via certbot)
4. **Configure alertas** (ex: monitoramento básico)
5. **Faça backup automático** (cron job)
6. **Documente credenciais** (senha, chaves SSH)

---

## 📞 Help & Suporte

Se encontrar problemas:

1. **Verificar logs**: `docker-compose logs whatsapp-api`
2. **Testar conectividade**: `curl http://localhost:5000/health`
3. **Verificar recursos**: `docker stats`
4. **Consultar documentação**: Veja `README.md` e `docker/README.md`

---

**Última atualização**: Fevereiro 2026  
**Versão**: v1.0.0
