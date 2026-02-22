# WhatsApp API - Exemplos de Requisições

## Variáveis
```
BASE_URL: http://localhost:8080
TENANT_ID: empresa1
USER_ID: usuario1
PHONE_NUMBER: 5585988888888
```

## 1. Criar Tenant

```http
POST http://localhost:8080/api/tenants HTTP/1.1
Content-Type: application/json

{
  "id": "empresa1",
  "name": "Empresa 1",
  "email": "admin@empresa1.com"
}
```

## 2. Obter Informações do Tenant

```http
GET http://localhost:8080/api/tenants/empresa1 HTTP/1.1
```

## 3. Criar Sessão

```http
POST http://localhost:8080/api/tenants/empresa1/sessions HTTP/1.1
Content-Type: application/json

{
  "user_id": "usuario1"
}
```

## 4. Obter QR Code (ESSE É O PRINCIPAL PARA CONECTAR)

```http
GET http://localhost:8080/api/tenants/empresa1/sessions/usuario1/qr HTTP/1.1
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

Então copie o valor de `qr_code` e decodifique de Base64 em uma ferramenta online como:
- https://www.base64-image-decoder.com/
- Ou use um decodificador online

Depois escaneie o QR Code com seu WhatsApp.

## 5. Verificar Status da Sessão

```http
GET http://localhost:8080/api/tenants/empresa1/sessions/usuario1/status HTTP/1.1
```

**Resposta quando conectado:**
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

## 6. Listar Sessões do Tenant

```http
GET http://localhost:8080/api/tenants/empresa1/sessions HTTP/1.1
```

## 7. Enviar Mensagem (Após estar CONNECTED)

```http
POST http://localhost:8080/api/tenants/empresa1/sessions/usuario1/send-message HTTP/1.1
Content-Type: application/json

{
  "number": "5585988888888",
  "message": "Olá! Esta é uma mensagem de teste da API"
}
```

## 8. Health Check

```http
GET http://localhost:8080/health HTTP/1.1
```

## 9. Deletar Sessão

```http
DELETE http://localhost:8080/api/tenants/empresa1/sessions/usuario1 HTTP/1.1
```

## 10. Deletar Tenant

```http
DELETE http://localhost:8080/api/tenants/empresa1 HTTP/1.1
```

---

## Fluxo Completo de Teste

1. **Criar Tenant**
   ```bash
   curl -X POST http://localhost:8080/api/tenants \
     -H "Content-Type: application/json" \
     -d '{"id":"empresa1","name":"Empresa 1","email":"admin@e1.com"}'
   ```

2. **Criar Sessão**
   ```bash
   curl -X POST http://localhost:8080/api/tenants/empresa1/sessions \
     -H "Content-Type: application/json" \
     -d '{"user_id":"usuario1"}'
   ```

3. **Obter QR Code**
   ```bash
   curl http://localhost:8080/api/tenants/empresa1/sessions/usuario1/qr
   ```
   - Decodifique o Base64 e escaneie com WhatsApp

4. **Aguardar conexão (30 segundos)**
   ```bash
   curl http://localhost:8080/api/tenants/empresa1/sessions/usuario1/status
   ```

5. **Enviar Mensagem (quando status = connected)**
   ```bash
   curl -X POST http://localhost:8080/api/tenants/empresa1/sessions/usuario1/send-message \
     -H "Content-Type: application/json" \
     -d '{"number":"5585988888888","message":"Teste API!"}'
   ```

---

## Múltiplas Sessões no Mesmo Tenant

```bash
# Sessão do Usuário 1
curl -X POST http://localhost:8080/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"usuario1"}'

# Sessão do Usuário 2
curl -X POST http://localhost:8080/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"usuario2"}'

# Listar todas
curl http://localhost:8080/api/tenants/empresa1/sessions
```

Cada sessão pode conectar com um número WhatsApp diferente!

---

## Múltiplos Tenants

```bash
# Tenant 1
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa1","name":"Empresa 1","email":"admin@e1.com"}'

# Tenant 2
curl -X POST http://localhost:8080/api/tenants \
  -H "Content-Type: application/json" \
  -d '{"id":"empresa2","name":"Empresa 2","email":"admin@e2.com"}'

# Sessão Tenant 1
curl -X POST http://localhost:8080/api/tenants/empresa1/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"vendedor1"}'

# Sessão Tenant 2
curl -X POST http://localhost:8080/api/tenants/empresa2/sessions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"operador1"}'
```

Total: 2 tenants com 1 sessão cada = 2 números WhatsApp conectados simultaneamente!
