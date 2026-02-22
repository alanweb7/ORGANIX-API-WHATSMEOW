# 🚀 Instalação em VPS Linux Ubuntu (SEM DOCKER)

Guia completo para instalar a **WhatsApp API** diretamente no servidor VPS sem usar Docker.

---

## 📋 Pré-requisitos

- **VPS Linux**: Ubuntu 20.04, 22.04 ou Debian 11+
- **Acesso SSH**: Com permissões sudo
- **RAM mínima**: 256MB (recomendado 512MB+)
- **Espaço em disco**: 100MB mínimo
- **Conexão**: SSH disponível

---

## ✅ Instalação Completa - Passo a Passo

### Passo 1: Conectar ao VPS

```bash
ssh usuario@seu-vps.com
# ou com IP direto
ssh root@192.168.1.100
```

### Passo 2: Atualizar o sistema

```bash
# Atualizar repositórios
sudo apt update

# Upgrade de pacotes
sudo apt upgrade -y

# Instalar dependências básicas
sudo apt install -y git curl wget build-essential
```

### Passo 3: Instalar Go 1.21+

#### Opção A: Via repositório (mais fácil - Ubuntu 22.04+)

```bash
# Instalar Go
sudo apt install -y golang-go

# Verificar versão
go version
```

#### Opção B: Download direto (compatível com todas as versões)

```bash
# Baixar Go 1.21
cd /tmp
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Extrair para /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Adicionar ao PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verificar
go version
```

### Passo 4: Preparar diretório da aplicação

```bash
# Criar diretório
sudo mkdir -p /opt/whatsapp-api
cd /opt/whatsapp-api

# Ajustar permissões
sudo chown $USER:$USER /opt/whatsapp-api

# Criar pastas necessárias
mkdir -p data
mkdir -p logs
mkdir -p sessions
```

### Passo 5: Obter o código-fonte

#### Opção A: Via Git (recomendado)

```bash
cd /opt/whatsapp-api

# Se for repositório privado, use SSH
git clone git@github.com:seu-usuario/seu-repo.git .
git pull origin main

# Ou HTTPS
git clone https://github.com/seu-usuario/seu-repo.git .
```

#### Opção B: Copiar arquivos via SCP

Na sua máquina local:

```bash
scp -r e:\apps\CLIENTES\ORGANIX-API-WHATSMEOW\api\* usuario@seu-vps.com:/opt/whatsapp-api/
```

### Passo 6: Configurar variáveis de ambiente

```bash
cd /opt/whatsapp-api

# Copiar arquivo de configuração
cp .env.example .env

# Editar variáveis
nano .env
```

**Exemplo de `.env`:**

```env
# Porta (pode ser qualquer porta > 1024 sem sudo, ou >= 1024 com sudo)
API_PORT=5000

# Diretórios
DB_PATH=/opt/whatsapp-api/data/whatsapp-api.db
SESSIONS_DIR=/opt/whatsapp-api/data/sessions
LOGS_DIR=/opt/whatsapp-api/logs

# WhatsApp (ajustar conforme necessário)
WHATSAPP_TIMEOUT=30
WHATSAPP_MAX_RETRIES=3

# Logging
LOG_LEVEL=info
LOG_FILE=/opt/whatsapp-api/logs/app.log
```

### Passo 7: Compilar a aplicação

```bash
cd /opt/whatsapp-api

# Build para Linux x64 (o padrão)
go build -o api .

# Verificar se foi criado
ls -lh api
file api
```

**Se receber erro sobre CGO:**

```bash
# Algumas dependências precisam de compilação C
sudo apt install -y gcc sqlite3 libsqlite3-dev

# Tentar build novamente
CGO_ENABLED=1 go build -o api .
```

### Passo 8: Testar a aplicação

```bash
cd /opt/whatsapp-api

# Rodar manualmente (primeira vez)
./api

# Esperado (saída no console):
# 2024/02/22 10:30:45 Servidor iniciado em :5000
# 2024/02/22 10:30:45 Banco de dados inicializado
```

**Testar em outro terminal SSH:**

```bash
# Em outro terminal do mesmo VPS
curl http://localhost:5000/health

# Esperado:
# {"status":"ok"}
```

Pressionar **Ctrl+C** no terminal da aplicação para parar.

---

## 🔧 Configurar para executar Automaticamente

### Opção 1: Systemd Service (Recomendado)

Criar arquivo de serviço:

```bash
sudo nano /etc/systemd/system/whatsapp-api.service
```

Conteúdo:

```ini
[Unit]
Description=WhatsApp API Service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=10
RemainAfterExit=yes

# Usuário que executará o serviço
User=API_USER
WorkingDirectory=/opt/whatsapp-api

# Comando a executar
ExecStart=/opt/whatsapp-api/api

# Redirecionamento de output
StandardOutput=journal
StandardError=journal
SyslogIdentifier=whatsapp-api

# Limites de recursos
MemoryMax=512M
CPUQuota=50%

[Install]
WantedBy=multi-user.target
```

**Ativar o serviço:**

```bash
# Recarregar systemd
sudo systemctl daemon-reload

# Ativar para iniciar no boot
sudo systemctl enable whatsapp-api.service

# Iniciar serviço
sudo systemctl start whatsapp-api.service

# Verificar status
sudo systemctl status whatsapp-api.service

# Ver logs
sudo journalctl -u whatsapp-api.service -f
```

### Opção 2: Supervisor (Alternativa)

```bash
# Instalar supervisor
sudo apt install -y supervisor

# Criar arquivo de configuração
sudo nano /etc/supervisor/conf.d/whatsapp-api.conf
```

Conteúdo:

```ini
[program:whatsapp-api]
directory=/opt/whatsapp-api
command=/opt/whatsapp-api/api
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/opt/whatsapp-api/logs/supervisord.log
environment=PATH=/usr/local/go/bin:/usr/bin
```

**Ativar:**

```bash
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start whatsapp-api

# Ver status
sudo supervisorctl status whatsapp-api
```

### Opção 3: Cron Job (Boot)

Adicionar ao crontab:

```bash
sudo crontab -e
```

Adicionar linha:

```cron
@reboot cd /opt/whatsapp-api && nohup ./api > logs/app.log 2>&1 &
```

---

## 🌐 Acessar a API Remotamente

### Opção 1: Acesso direto pela porta

A API estará acessível em:

```
http://seu-ip-vps:5000/health
http://seu-ip-vps:5000/api/tenant/{tenantId}/sessions
```

### Opção 2: Nginx Reverse Proxy

Instalar Nginx:

```bash
sudo apt install -y nginx
```

Criar arquivo de configuração:

```bash
sudo nano /etc/nginx/sites-available/whatsapp-api
```

Conteúdo:

```nginx
server {
    listen 80;
    server_name seu-dominio.com api.seu-dominio.com;

    location / {
        proxy_pass http://127.0.0.1:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts para conexões longas
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

**Ativar:**

```bash
# Criar link simbólico
sudo ln -s /etc/nginx/sites-available/whatsapp-api /etc/nginx/sites-enabled/

# Testar configuração
sudo nginx -t

# Reiniciar Nginx
sudo systemctl restart nginx

# Ver status
sudo systemctl status nginx
```

Acesso via domínio:

```
http://seu-dominio.com/health
http://seu-dominio.com/api/tenant/{tenantId}/sessions
```

### Opção 3: SSL/HTTPS com Let's Encrypt

```bash
# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obter certificado
sudo certbot --nginx -d seu-dominio.com -d www.seu-dominio.com

# Renovação automática já vem ativada
sudo systemctl enable certbot.timer
```

Acesso seguro:

```
https://seu-dominio.com/health
```

---

## 📊 Monitoramento e Manutenção

### Ver status do serviço

```bash
# Via systemd
sudo systemctl status whatsapp-api.service

# Via supervisor
sudo supervisorctl status whatsapp-api
```

### Ver logs em tempo real

```bash
# Via systemd (últimas 100 linhas + follow)
sudo journalctl -u whatsapp-api.service -f --lines=100

# Via arquivo de log
tail -f /opt/whatsapp-api/logs/app.log

# Ver logs históricos
journalctl -u whatsapp-api.service --since "2024-02-22"
```

### Reiniciar a aplicação

```bash
# Parar
sudo systemctl stop whatsapp-api.service

# Iniciar
sudo systemctl start whatsapp-api.service

# Reiniciar
sudo systemctl restart whatsapp-api.service

# Recarregar (sem interrupção)
sudo systemctl reload whatsapp-api.service
```

### Ver uso de recursos

```bash
# Verificar processo
ps aux | grep api

# Ver recursos em tempo real
top -p $(pidof api)

# Usar htop (mais amigável)
sudo apt install -y htop
htop -p $(pidof api)
```

### Fazer backup de dados

```bash
# Fazer backup do banco de dados
tar -czf ~/whatsapp-api-backup-$(date +%Y%m%d_%H%M%S).tar.gz \
  /opt/whatsapp-api/data/

# Copiar para máquina local
scp usuario@seu-vps.com:~/*.tar.gz ./backups/

# Restaurar backup
tar -xzf whatsapp-api-backup-20240222_103045.tar.gz -C /
```

---

## 🔄 Atualizar a aplicação

### Se usando Git:

```bash
cd /opt/whatsapp-api

# Parar serviço
sudo systemctl stop whatsapp-api.service

# Atualizar código
git pull origin main

# Recompilar
go build -o api .

# Iniciar novamente
sudo systemctl start whatsapp-api.service

# Verificar
sudo systemctl status whatsapp-api.service
```

### Se usando cópia de arquivos:

```bash
# Na máquina local
scp -r e:\apps\CLIENTES\ORGANIX-API-WHATSMEOW\api\* usuario@seu-vps.com:/opt/whatsapp-api/

# No VPS
ssh usuario@seu-vps.com
cd /opt/whatsapp-api
sudo systemctl stop whatsapp-api.service
go build -o api .
sudo systemctl start whatsapp-api.service
```

---

## ⚠️ Troubleshooting

### Erro: "permission denied" ao executar

```bash
# Adicionar permissão de execução
chmod +x /opt/whatsapp-api/api

# Ou com sudo
sudo chmod +x /opt/whatsapp-api/api
```

### Erro: "go: command not found"

```bash
# Verificar install de Go
go version

# Se não encontra, recarregar PATH
source ~/.bashrc

# Ou adicionar permissão ao usuário do serviço
sudo su - seu-usuario
source ~/.bashrc
```

### Erro: "cannot connect to localhost:5000"

```bash
# Verificar se a porta está aberta
sudo netstat -tulpn | grep 5000

# Se não mostrar, a aplicação não está rodando
sudo systemctl status whatsapp-api.service

# Ver logs de erro
sudo journalctl -u whatsapp-api.service -n 50
```

### Erro: "database is locked"

SQLite não está bem configurado. Tente:

```bash
# Parar a aplicação
sudo systemctl stop whatsapp-api.service

# Aguardar 5 segundos
sleep 5

# Iniciar novamente
sudo systemctl start whatsapp-api.service

# Verificar locks
lsof /opt/whatsapp-api/data/whatsapp-api.db
```

### A aplicação consome muita memória

Limitar no arquivo de serviço:

```bash
sudo nano /etc/systemd/system/whatsapp-api.service

# Adicionar ou alterar:
# MemoryMax=256M
# MemoryLimit=256M

sudo systemctl daemon-reload
sudo systemctl restart whatsapp-api.service
```

### Porta já está em uso

```bash
# Ver qual processo está usando a porta
sudo lsof -i :5000

# Matar processo específico
sudo kill -9 PID

# Ou usar outra porta
nano /opt/whatsapp-api/.env
# Alterar API_PORT=8080 (por exemplo)
```

---

## 🔒 Segurança

### Firewall (UFW)

```bash
# Ativar firewall
sudo ufw enable

# Permitir SSH
sudo ufw allow 22

# Permitir HTTP (se usar Nginx)
sudo ufw allow 80

# Permitir HTTPS
sudo ufw allow 443

# Bloquear porta 5000 (acessível só via Nginx)
sudo ufw deny 5000

# Ver regras
sudo ufw status
```

### Restringir acesso à API

```bash
# No Nginx, aceitar só IPs específicos
sudo nano /etc/nginx/sites-available/whatsapp-api

# Adicionar:
# allow 192.168.1.0/24;
# allow YOUR_IP;
# deny all;

sudo nginx -t
sudo systemctl reload nginx
```

### Proteger banco de dados

```bash
# Ajustar permissões
chmod 600 /opt/whatsapp-api/data/*.db
chmod 700 /opt/whatsapp-api/data/

# Ou mais restritivo
sudo chown root:root /opt/whatsapp-api/data/
sudo chmod 700 /opt/whatsapp-api/data/
```

### Monitorar acessos

```bash
# Logs de Nginx
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log

# Filtrar por erro
grep "error" /var/log/nginx/error.log

# Logs de acesso
grep "GET /health" /var/log/nginx/access.log
```

---

## 📈 Próximas Etapas

1. **Monitorar aplicação**: Use htop ou ferramentas de APM
2. **Configurar alertas**: CPU/Memória acima de limiar
3. **Backup automático**: Cron job para backup diário
4. **Logs centralizados**: ELK, Splunk ou similar
5. **Load balancer**: Se múltiplos servidores
6. **CDN**: Cloudflare ou similar para cache/proteção

---

## ✨ Comparativo: Docker vs Sem Docker

| Aspecto | Sem Docker | Com Docker |
|---------|-----------|-----------|
| **Overhead** | Mínimo | Maior (imagem) |
| **Memória** | ~30-50MB | ~100-150MB |
| **Startup** | < 1 segundo | 2-5 segundos |
| **Atualização** | Git pull + rebuild | Pull imagem |
| **Isolamento** | Compartilha sistema | Isolado |
| **Escaling** | Difícil | Fácil (orchestração) |
| **Compatibilidade** | Versão Go pode variar | Consistente |

---

## 📞 Help & Suporte

Se encontrar problemas:

1. **Ver logs**: `sudo journalctl -u whatsapp-api.service -f`
2. **Testar conectividade**: `curl http://localhost:5000/health`
3. **Verificar recursos**: `top -p $(pidof api)`
4. **Consultar docs**: Veja `README.md` e `MANUAL.md`

---

**Última atualização**: Fevereiro 2026  
**Versão**: v1.0.0  
**Testado em**: Ubuntu 20.04, 22.04 | Debian 11
