# Certificate Management REST API

## Background and Goals

After deploying an internal OpenVPN service, enterprises often need to provision, renew, revoke, and distribute VPN credentials through their own business systems (such as OA portals, ops platforms, or employee self-service portals), rather than relying on administrators to use the web UI manually.

This feature provides a **token-authenticated REST API** that covers the full certificate lifecycle, enabling enterprise systems to:

- Automatically create VPN accounts and deliver `.ovpn` profiles when employees join
- Revoke old certificates when employees leave or devices are replaced
- Trigger certificate renewal before expiration via business workflows
- Integrate VPN management with existing IAM, HR, or ticketing systems

The API is independent of the web UI: the web interface uses session-based login, while the API uses a static token suitable for server-to-server calls.

## Configuration

### ApiToken

Set in `conf/app.conf`:

```ini
ApiToken = "your-secret-api-token"
```

### Docker Environment Variable Override

When deployed with Docker, the token can be overridden via environment variable (takes precedence over the config file):

```yaml
environment:
  - OPENVPN_API_TOKEN=your-secret-api-token
```

If neither `ApiToken` nor `OPENVPN_API_TOKEN` is set, all API requests return `401`.

## Authentication

Every request must include one of the following headers:

```http
Authorization: Bearer your-secret-api-token
```

or:

```http
X-API-Token: your-secret-api-token
```

Authentication failure response:

```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

HTTP status code: `401`

## General Information

| Item | Value |
|------|-------|
| Base URL | `http://<host>:<port>/api/v1/certificates` |
| Default port | `8080` (see `HttpPort` in `app.conf`) |
| Content-Type | `application/json` (for POST / DELETE bodies) |
| Response format | JSON (except `.ovpn` download) |

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/certificates` | Create certificate |
| `GET` | `/api/v1/certificates/:name` | Download `.ovpn` profile |
| `POST` | `/api/v1/certificates/:name/renew` | Renew certificate |
| `POST` | `/api/v1/certificates/:name/revoke` | Revoke certificate |
| `DELETE` | `/api/v1/certificates/:name` | Delete revoked certificate |

> **About `name`**: `name` must be unique when creating a certificate; duplicate names are rejected. Use `name` as the stable account identifier in your business system (e.g. employee ID). Renew, revoke, and delete only need `name` in the URL path; the server resolves the matching certificate from the PKI index.

---

## Create Certificate

Provision a VPN account for an employee or device.

**Request**

```http
POST /api/v1/certificates
Authorization: Bearer your-secret-api-token
Content-Type: application/json
```

```json
{
  "name": "johndoe",
  "staticip": "10.0.70.50",
  "passphrase": "",
  "expire_days": "825",
  "email": "johndoe@company.com",
  "country": "US",
  "province": "California",
  "city": "San Francisco",
  "org": "MyCompany",
  "org_unit": "IT",
  "tfa_name": "",
  "tfa_issuer": ""
}
```

**Parameters**

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Certificate name; must be unique; map to user ID or employee number in your system |
| `staticip` | No | Static IP; leave empty for dynamic assignment |
| `passphrase` | No | Private key passphrase |
| `expire_days` | No | Validity in days; defaults to EasyRSA config |
| `email` | No | Defaults to EasyRSA config |
| `country` / `province` / `city` / `org` / `org_unit` | No | Certificate DN fields; default to EasyRSA config |
| `tfa_name` / `tfa_issuer` | No | 2FA settings; leave empty for standard certificates |

**Success response** `200`

```json
{
  "status": "success",
  "message": "Certificate \"johndoe\" has been created"
}
```

**curl example**

```bash
curl -X POST http://localhost:8080/api/v1/certificates \
  -H "Authorization: Bearer your-secret-api-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "johndoe", "email": "johndoe@company.com"}'
```

---

## Download .ovpn Profile

Retrieve the client connection configuration for delivery to employees or devices.

**Request**

```http
GET /api/v1/certificates/johndoe
Authorization: Bearer your-secret-api-token
```

**Success response** `200`

Returns `.ovpn` file content:

```http
Content-Type: application/x-openvpn-profile
Content-Disposition: attachment; filename="johndoe.ovpn"
```

**curl example**

```bash
curl -o johndoe.ovpn \
  -H "Authorization: Bearer your-secret-api-token" \
  http://localhost:8080/api/v1/certificates/johndoe
```

**Certificate not found** `404`

```json
{
  "status": "error",
  "message": "Certificate not found"
}
```

---

## Renew Certificate

**Request**

```http
POST /api/v1/certificates/johndoe/renew
Authorization: Bearer your-secret-api-token
```

No request body is required. The server looks up the active certificate by `name` and renews it.

**curl example**

```bash
curl -X POST http://localhost:8080/api/v1/certificates/johndoe/renew \
  -H "Authorization: Bearer your-secret-api-token"
```

**Success response** `200`

```json
{
  "status": "success",
  "message": "Certificate \"johndoe\" has been renewed"
}
```

Renewal creates a new certificate with the same name. The old certificate remains in the list until revoked after the new profile is delivered.

---

## Revoke Certificate

Use when an employee leaves or a device is decommissioned, preventing further VPN connections.

**Request**

```http
POST /api/v1/certificates/johndoe/revoke
Authorization: Bearer your-secret-api-token
```

No request body is required. The server looks up the active certificate by `name` and revokes it.

**curl example**

```bash
curl -X POST http://localhost:8080/api/v1/certificates/johndoe/revoke \
  -H "Authorization: Bearer your-secret-api-token"
```

**Success response** `200`

```json
{
  "status": "success",
  "message": "Certificate \"johndoe\" has been revoked"
}
```

> Restart the OpenVPN service after revocation to immediately disconnect active clients (same behavior as the web UI).

---

## Delete Certificate

Permanently remove a **revoked** certificate record.

**Request**

```http
DELETE /api/v1/certificates/johndoe
Authorization: Bearer your-secret-api-token
```

No request body is required. The server looks up the revoked certificate by `name` and deletes it.

**curl example**

```bash
curl -X DELETE http://localhost:8080/api/v1/certificates/johndoe \
  -H "Authorization: Bearer your-secret-api-token"
```

**Success response** `200`

```json
{
  "status": "success",
  "message": "Certificate \"johndoe\" has been deleted"
}
```

---

## Error Responses

Except for the download endpoint, errors are returned as JSON:

```json
{
  "status": "error",
  "message": "error description"
}
```

| HTTP status | Meaning |
|-------------|---------|
| `400` | Invalid parameters or operation failure |
| `401` | Missing or invalid token |
| `404` | Certificate not found (download endpoint) |

---

## Typical Enterprise Integration Flows

### Employee Onboarding

```
Business system ──POST /certificates──▶ Create certificate
                ◀── 200 success ────────
                ──GET /certificates/:name──▶ Download .ovpn
                ◀── .ovpn file ───────────
                ── Store in object storage / email / portal download
```

### Employee Offboarding

```
Business system ──POST /certificates/:name/revoke──▶ Revoke certificate
                ──DELETE /certificates/:name───────▶ Clean up record (optional)
                ── Notify ops to restart OpenVPN (if immediate disconnect needed)
```

### Certificate Renewal

```
Business system ──POST /certificates/:name/renew──▶ Renew
                ──GET /certificates/:name─────────▶ Fetch new .ovpn
                ──POST /certificates/:name/revoke──▶ Revoke old cert (optional)
```

We recommend maintaining a VPN account ledger in your business system with at least:

| Field | Description |
|-------|-------------|
| `user_id` | User identifier in your system |
| `cert_name` | Corresponds to API `name` |
| `status` | active / revoked |

---

## API Verification Tool

A built-in verification program is included for smoke and integration testing:

```bash
# Smoke tests (auth, routing, validation)
go run ./cmd/verify-api

# Custom base URL and token
go run ./cmd/verify-api \
  -base http://localhost:8080 \
  -token your-secret-api-token

# Full lifecycle tests (requires complete PKI environment)
go run ./cmd/verify-api -integration
```

Build as a standalone tool:

```bash
go build -o verify-api ./cmd/verify-api
./verify-api
```

---

## Related Code

| File | Description |
|------|-------------|
| `controllers/api-certificates.go` | API controller |
| `controllers/api-base.go` | Token authentication |
| `routers/router.go` | Route registration |
| `cmd/verify-api/main.go` | API verification tool |

[中文版本](../zh/certificate-api.md)
