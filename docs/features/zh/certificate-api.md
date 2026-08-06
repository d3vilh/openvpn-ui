# 证书管理 REST API

## 背景与目标

企业部署内部 OpenVPN 后，往往需要将 VPN 账号的开通、续期、吊销与配置下发，纳入自有业务系统（如 OA、运维平台、员工门户）统一编排，而不是让管理员手工登录 Web 界面操作。

本特性提供一套 **Token 鉴权的 REST API**，覆盖证书全生命周期管理，使企业系统可以：

- 员工入职时自动创建 VPN 账号并下发 `.ovpn` 配置
- 员工离职或设备更换时自动吊销旧证书
- 证书到期前通过业务系统触发续期
- 将 VPN 管理流程与现有 IAM / HR / 工单系统对接

API 与 Web 管理界面相互独立：Web 仍使用 Session 登录；API 使用固定 Token，适合服务端对服务端调用。

## 配置

### ApiToken

在 `conf/app.conf` 中配置：

```ini
ApiToken = "your-secret-api-token"
```

### Docker 环境变量覆盖

Docker 部署时可通过环境变量覆盖（优先级高于配置文件）：

```yaml
environment:
  - OPENVPN_API_TOKEN=your-secret-api-token
```

未配置 `ApiToken` 且未设置 `OPENVPN_API_TOKEN` 时，所有 API 请求将返回 `401`。

## 鉴权

每个请求须携带以下任一 Header：

```http
Authorization: Bearer your-secret-api-token
```

或：

```http
X-API-Token: your-secret-api-token
```

鉴权失败响应：

```json
{
  "status": "error",
  "message": "Unauthorized"
}
```

HTTP 状态码：`401`

## 基础信息

| 项目 | 值 |
|------|-----|
| Base URL | `http://<host>:<port>/api/v1/certificates` |
| 默认端口 | `8080`（以 `app.conf` 中 `HttpPort` 为准） |
| Content-Type | `application/json`（POST / DELETE 请求体） |
| 响应格式 | JSON（下载 `.ovpn` 除外） |

## 接口一览

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/certificates` | 创建证书 |
| `GET` | `/api/v1/certificates/:name` | 下载 `.ovpn` 配置 |
| `POST` | `/api/v1/certificates/:name/renew` | 续期证书 |
| `POST` | `/api/v1/certificates/:name/revoke` | 吊销证书 |
| `DELETE` | `/api/v1/certificates/:name` | 删除已吊销的证书 |

> **关于 `name`**：创建时 `name` 必须唯一，已存在同名证书则无法再次创建。业务系统可用 `name` 作为稳定账号标识（如工号）。续期、吊销、删除只需在 URL 路径中传入 `name`，服务端会自动从 PKI 索引解析对应证书。

---

## 创建证书

为员工或设备开通 VPN 账号。

**请求**

```http
POST /api/v1/certificates
Authorization: Bearer your-secret-api-token
Content-Type: application/json
```

```json
{
  "name": "zhangsan",
  "staticip": "10.0.70.50",
  "passphrase": "",
  "expire_days": "825",
  "email": "zhangsan@company.com",
  "country": "CN",
  "province": "Beijing",
  "city": "Beijing",
  "org": "MyCompany",
  "org_unit": "IT",
  "tfa_name": "",
  "tfa_issuer": ""
}
```

**参数说明**

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 是 | 证书名称，必须唯一，建议与业务系统用户 ID 或工号对应 |
| `staticip` | 否 | 静态 IP，留空则动态分配 |
| `passphrase` | 否 | 私钥密码 |
| `expire_days` | 否 | 有效期（天），默认取 EasyRSA 配置 |
| `email` | 否 | 默认取 EasyRSA 配置 |
| `country` / `province` / `city` / `org` / `org_unit` | 否 | 证书 DN 字段，默认取 EasyRSA 配置 |
| `tfa_name` / `tfa_issuer` | 否 | 2FA 相关，留空表示普通证书 |

**成功响应** `200`

```json
{
  "status": "success",
  "message": "Certificate \"zhangsan\" has been created"
}
```

**curl 示例**

```bash
curl -X POST http://localhost:8080/api/v1/certificates \
  -H "Authorization: Bearer your-secret-api-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "zhangsan", "email": "zhangsan@company.com"}'
```

---

## 下载 .ovpn 配置

获取客户端连接配置文件，可下发给员工或设备。

**请求**

```http
GET /api/v1/certificates/zhangsan
Authorization: Bearer your-secret-api-token
```

**成功响应** `200`

返回 `.ovpn` 文件内容：

```http
Content-Type: application/x-openvpn-profile
Content-Disposition: attachment; filename="zhangsan.ovpn"
```

**curl 示例**

```bash
curl -o zhangsan.ovpn \
  -H "Authorization: Bearer your-secret-api-token" \
  http://localhost:8080/api/v1/certificates/zhangsan
```

**证书不存在** `404`

```json
{
  "status": "error",
  "message": "Certificate not found"
}
```

---

## 续期证书

**请求**

```http
POST /api/v1/certificates/zhangsan/renew
Authorization: Bearer your-secret-api-token
```

无需请求体，服务端根据 `name` 查找当前有效证书并续期。

**curl 示例**

```bash
curl -X POST http://localhost:8080/api/v1/certificates/zhangsan/renew \
  -H "Authorization: Bearer your-secret-api-token"
```

**成功响应** `200`

```json
{
  "status": "success",
  "message": "Certificate \"zhangsan\" has been renewed"
}
```

续期后会生成同名新证书，旧证书仍在列表中，待新配置下发后可再次调用吊销接口。

---

## 吊销证书

用于员工离职、设备报废等场景，使证书无法再连接 VPN。

**请求**

```http
POST /api/v1/certificates/zhangsan/revoke
Authorization: Bearer your-secret-api-token
```

无需请求体，服务端根据 `name` 查找当前有效证书并吊销。

**curl 示例**

```bash
curl -X POST http://localhost:8080/api/v1/certificates/zhangsan/revoke \
  -H "Authorization: Bearer your-secret-api-token"
```

**成功响应** `200`

```json
{
  "status": "success",
  "message": "Certificate \"zhangsan\" has been revoked"
}
```

> 吊销后需重启 OpenVPN 服务才能立即断开已连接客户端（与 Web 界面行为一致）。

---

## 删除证书

永久移除**已吊销**的证书记录。

**请求**

```http
DELETE /api/v1/certificates/zhangsan
Authorization: Bearer your-secret-api-token
```

无需请求体，服务端根据 `name` 查找已吊销的证书记录并删除。

**curl 示例**

```bash
curl -X DELETE http://localhost:8080/api/v1/certificates/zhangsan \
  -H "Authorization: Bearer your-secret-api-token"
```

**成功响应** `200`

```json
{
  "status": "success",
  "message": "Certificate \"zhangsan\" has been deleted"
}
```

---

## 错误响应

除下载接口返回文件外，错误统一为 JSON：

```json
{
  "status": "error",
  "message": "错误描述"
}
```

| HTTP 状态码 | 含义 |
|-------------|------|
| `400` | 参数错误或业务操作失败 |
| `401` | Token 缺失或无效 |
| `404` | 证书不存在（下载接口） |

---

## 企业集成典型流程

### 员工入职

```
业务系统 ──POST /certificates──▶ 创建证书
         ◀── 200 success ────────
         ──GET /certificates/:name──▶ 下载 .ovpn
         ◀── .ovpn 文件 ───────────
         ── 存入对象存储 / 邮件下发 / 门户下载
```

### 员工离职

```
业务系统 ──POST /certificates/:name/revoke──▶ 吊销证书
         ──DELETE /certificates/:name───────▶ 清理记录（可选）
         ── 通知运维重启 OpenVPN（如需立即断连）
```

### 证书续期

```
业务系统 ──POST /certificates/:name/renew──▶ 续期
         ──GET /certificates/:name─────────▶ 获取新 .ovpn
         ──POST /certificates/:name/revoke──▶ 吊销旧证（可选）
```

建议业务系统维护一张 VPN 账号台账，至少记录：

| 字段 | 说明 |
|------|------|
| `user_id` | 业务系统用户标识 |
| `cert_name` | 对应 API 中的 `name` |
| `status` | active / revoked |

---

## API 验证工具

项目内置验证程序，可用于冒烟测试或集成测试：

```bash
# 冒烟测试（鉴权、路由、参数校验）
go run ./cmd/verify-api

# 指定地址和 Token
go run ./cmd/verify-api \
  -base http://localhost:8080 \
  -token your-secret-api-token

# 完整生命周期测试（需完整 PKI 环境）
go run ./cmd/verify-api -integration
```

编译为独立工具：

```bash
go build -o verify-api ./cmd/verify-api
./verify-api
```

---

## 相关代码

| 文件 | 说明 |
|------|------|
| `controllers/api-certificates.go` | API 控制器 |
| `controllers/api-base.go` | Token 鉴权 |
| `routers/router.go` | 路由注册 |
| `cmd/verify-api/main.go` | API 验证程序 |

[English version](../en/certificate-api.md)
