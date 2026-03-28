# HiTechCloud Agent

> [Tiếng Việt](#tiếng-việt) | [English](#english)

---

## Tiếng Việt

### Giới thiệu

HiTechCloud Agent là dịch vụ backend viết bằng Go, chạy trực tiếp trên máy chủ để cung cấp API quản trị hạ tầng cho hệ sinh thái HiTechCloud. Dự án tập trung vào quản lý website, ứng dụng, container, backup, công cụ hệ thống và các tính năng AI/runtime mở rộng.

### Tính năng chính

- Quản lý website, Nginx, SSL, ACME, DNS account và hosting account.
- Quản lý Docker, container, image, network, volume, compose và template triển khai.
- Cài đặt, đồng bộ, nâng cấp và quản lý ứng dụng từ app store.
- Backup, restore, upload recovery và kết nối nhiều nhà cung cấp cloud storage.
- Quản lý database, cronjob, process, dashboard, log và monitoring.
- Toolbox hệ thống: firewall, SSH, FTP, Fail2Ban, ClamAV, device/system tools.
- Tích hợp AI tools như Ollama, MCP server, GPU/runtime và AI agent/provider.
- Hỗ trợ WebSocket cho các tác vụ realtime như terminal/process hoặc luồng theo dõi trạng thái.

### Kiến trúc tổng quan

Agent này thường được cài trên từng máy chủ Linux và hoạt động như một lớp điều khiển cục bộ:

- Cung cấp API qua các nhóm route `/api/v1` và `/api/v2`.
- Lưu trữ dữ liệu cục bộ bằng SQLite cho settings, task, monitor, alert và các metadata khác.
- Tích hợp sâu với tài nguyên hệ thống như Docker socket, filesystem, firewall, cron và network.
- Hỗ trợ chạy qua TCP port hoặc Unix socket trong chế độ master.

### Cấu trúc chính

```text
cmd/server          Điểm khởi động CLI/server
app/api/v2          API handlers
router/             Khai báo nhóm route theo domain
service/            Nghiệp vụ chính
repo/               Truy cập dữ liệu
model/              Khai báo model
init/               Luồng khởi tạo hệ thống
utils/              Thư viện tiện ích và tích hợp hệ thống
cron/               Tác vụ định kỳ
middleware/         Xử lý auth, operation metadata, certificate
```

### Yêu cầu hệ thống

- Linux server.
- Go `1.25.7` hoặc phiên bản tương thích với `go.mod`.
- Quyền truy cập tới các tài nguyên hệ thống cần thiết nếu dùng đầy đủ tính năng:
  - Docker socket, mặc định `unix:///var/run/docker.sock`
  - các thư mục hệ thống dưới `/opt`, `/etc`, `/run`, `/usr/local/bin`
  - firewall/iptables/ufw
- Các dịch vụ/phần mềm liên quan tùy theo module sử dụng:
  - Docker / Docker Compose
  - Nginx hoặc OpenResty
  - SSH
  - Fail2Ban
  - ClamAV
  - FTP service

### Cấu hình

Cấu hình mặc định được embed từ `cmd/server/conf/app.yaml`.

Giá trị mặc định hiện tại:

- `base.install_dir: /opt`
- `base.mode: dev`
- `log.level: debug`

Trong môi trường dev, agent có thể ưu tiên file cấu hình ngoài tại `/opt/HiTechCloud/conf/app.yaml`.

Port server được lấy từ `global.CONF.Base.Port`. Nếu không cấu hình, hệ thống dùng mặc định `8443`.

### Build và chạy

#### Chạy trực tiếp

```bash
go run ./cmd/server
```

#### Build binary

```bash
go build -o HiTechCloud-agent ./cmd/server
```

Sau khi khởi động, server sẽ:

1. nạp cấu hình
2. khởi tạo thư mục làm việc
3. khởi tạo logger
4. mở database SQLite cục bộ
5. chạy migration
6. khởi tạo i18n, cache, app metadata, validator
7. chạy cronjob
8. khởi tạo hook, firewall và các thành phần phụ trợ
9. mở HTTP hoặc HTTPS server

### Cơ chế lắng nghe

- **TCP server**: lắng nghe tại `0.0.0.0:<port>`.
- **Master mode**: lắng nghe qua Unix socket `/etc/HiTechCloud/agent.sock`.

Nếu trong settings đã có `ServerCrt` và `ServerKey`, agent sẽ ưu tiên chạy HTTPS với TLS >= 1.2. Nếu không có, hệ thống sẽ fallback sang HTTP.

### Bảo mật và xác thực

- Public endpoint mặc định: `GET /api/v1/health`
- Các route private dưới `/api/v1` và `/api/v2` dùng middleware API key.
- Header xác thực hiện tại: `X-API-Key`
- API key được hash bằng SHA-256 và so sánh theo constant-time.
- Hỗ trợ allowlist IP thông qua setting `ApiKeyAllowedIPs`.

### Các nhóm API chính

#### Website / Hosting

Quản lý domain, SSL, ACME, rewrite, redirect, proxy, auth, CORS, anti-leech, PHP, composer và các thành phần hosting liên quan.

#### Container / Docker

Quản lý container, image, volume, network, compose, template, log, stats, exec và thao tác file trong container.

#### App Store / Ứng dụng

Đồng bộ danh sách app, tìm kiếm, cài đặt, cập nhật, quản lý ứng dụng đã cài và ignore upgrade.

#### Backup / Restore

Quản lý cấu hình backup, token, bucket, backup thủ công/định kỳ, restore và tải tệp recovery.

#### System Toolbox

Cung cấp các công cụ thao tác với thiết bị, firewall, FTP, Fail2Ban, ClamAV và các tác vụ hệ thống khác.

#### AI / Runtime

Hỗ trợ các thành phần AI như Ollama, MCP server, GPU monitoring/runtime, AI account/provider/channel/plugin/security/skill.

#### Monitoring / Alert / Cron

Theo dõi trạng thái hệ thống, dashboard, tiến trình, cảnh báo và tác vụ định kỳ.

### Database cục bộ

Dự án sử dụng SQLite cho nhiều miền dữ liệu khác nhau, bao gồm:

- `agent.db`
- `task.db`
- `monitor.db`
- `gpu_monitor.db`
- `alert.db`
- `core.db` (tùy trường hợp)

### Logging

Hệ thống sử dụng Logrus kết hợp ghi log file và stdout. Thành phần logging được khởi tạo sớm trong quá trình start để hỗ trợ theo dõi toàn bộ vòng đời ứng dụng.

### Đa ngôn ngữ

Dự án tích hợp i18n và hiện có các gói ngôn ngữ embed như:

- `en`
- `zh`
- `zh-Hant`
- `pt-BR`
- `ja`
- `ru`
- `ms`
- `ko`
- `tr`
- `es-ES`

### Cron và tác vụ nền

Hệ thống có các tác vụ định kỳ như:

- đồng bộ thời gian NTP
- job website hằng ngày
- job SSL hằng ngày
- đồng bộ app store
- refresh backup token
- khôi phục cronjob đã bật từ database

### Ghi chú triển khai

- Dự án tích hợp mạnh với hệ điều hành và dịch vụ hệ thống, vì vậy nên chạy trong môi trường server thật hoặc môi trường test có đủ quyền cần thiết.
- Nếu chạy production, nên đặt agent sau reverse proxy hoặc cấu hình TLS đầy đủ.
- Một số tính năng phụ thuộc trực tiếp vào phần mềm ngoài như Docker, Nginx, SSH, firewall tools hoặc antivirus service.

### License

Dự án được phát hành theo giấy phép `GPL-3.0`. Xem thêm tại file `LICENSE`.

### Module và thư mục liên quan

- `cmd/server/main.go`
- `server/server.go`
- `init/router/router.go`
- `middleware/apikey.go`
- `router/`
- `app/api/v2/`
- `utils/`
- `cron/`
- `init/`

---

## English

### Overview

HiTechCloud Agent is a Go-based backend service that runs directly on a server and exposes infrastructure management APIs for the HiTechCloud ecosystem. It focuses on website management, application lifecycle, containers, backups, system utilities, and extended AI/runtime capabilities.

### Key Features

- Manage websites, Nginx, SSL, ACME, DNS accounts, and hosting accounts.
- Manage Docker, containers, images, networks, volumes, compose stacks, and deployment templates.
- Install, sync, upgrade, and manage applications from the app store.
- Support backup, restore, recovery upload, and multiple cloud storage providers.
- Handle databases, cronjobs, processes, dashboards, logs, and monitoring.
- Provide system toolbox features such as firewall, SSH, FTP, Fail2Ban, ClamAV, and device/system tools.
- Integrate AI capabilities such as Ollama, MCP server, GPU/runtime, and AI agent/provider components.
- Support WebSocket-based realtime operations such as terminal/process streaming and status tracking.

### Architecture Summary

The agent is typically installed on each Linux server and acts as a local control layer:

- Exposes APIs under `/api/v1` and `/api/v2`.
- Stores local metadata in SQLite for settings, tasks, monitoring, alerts, and related runtime data.
- Integrates deeply with system resources such as the Docker socket, filesystem, firewall, cron, and networking.
- Supports both TCP listening and Unix socket mode for master deployments.

### Main Structure

```text
cmd/server          CLI/server entrypoint
app/api/v2          API handlers
router/             Domain-based route groups
service/            Business logic
repo/               Data access layer
model/              Data models
init/               System bootstrap flow
utils/              Utilities and system integrations
cron/               Scheduled jobs
middleware/         Auth, operation metadata, certificate handling
```

### System Requirements

- Linux server.
- Go `1.25.7` or a compatible version matching `go.mod`.
- Access to required system resources when using full functionality:
  - Docker socket, default `unix:///var/run/docker.sock`
  - system paths under `/opt`, `/etc`, `/run`, `/usr/local/bin`
  - firewall/iptables/ufw
- Related software depending on enabled modules:
  - Docker / Docker Compose
  - Nginx or OpenResty
  - SSH
  - Fail2Ban
  - ClamAV
  - FTP service

### Configuration

Default configuration is embedded from `cmd/server/conf/app.yaml`.

Current default values include:

- `base.install_dir: /opt`
- `base.mode: dev`
- `log.level: debug`

In development mode, the agent may prefer an external config file at `/opt/HiTechCloud/conf/app.yaml`.

The server port is read from `global.CONF.Base.Port`. If not configured, the default port is `8443`.

### Build and Run

#### Run directly

```bash
go run ./cmd/server
```

#### Build binary

```bash
go build -o HiTechCloud-agent ./cmd/server
```

On startup, the server will:

1. load configuration
2. initialize working directories
3. initialize logging
4. open local SQLite databases
5. run migrations
6. initialize i18n, cache, app metadata, and validators
7. start cron jobs
8. initialize hooks, firewall, and supporting components
9. start the HTTP or HTTPS server

### Listening Modes

- **TCP server**: listens on `0.0.0.0:<port>`
- **Master mode**: listens through Unix socket `/etc/HiTechCloud/agent.sock`

If `ServerCrt` and `ServerKey` are available in settings, the agent prefers HTTPS with TLS >= 1.2. Otherwise, it falls back to HTTP.

### Security and Authentication

- Default public endpoint: `GET /api/v1/health`
- Private routes under `/api/v1` and `/api/v2` are protected by API key middleware.
- Current authentication header: `X-API-Key`
- API keys are hashed with SHA-256 and compared in constant time.
- Optional IP allowlist is supported via `ApiKeyAllowedIPs`.

### Main API Domains

#### Website / Hosting

Manage domains, SSL, ACME, rewrites, redirects, proxies, auth, CORS, anti-leech rules, PHP, Composer, and related hosting features.

#### Container / Docker

Manage containers, images, volumes, networks, compose stacks, templates, logs, stats, exec, and file operations inside containers.

#### App Store / Applications

Sync app catalogs, search, install, update, manage installed apps, and ignore upgrades.

#### Backup / Restore

Manage backup settings, tokens, buckets, manual/scheduled backups, restore flows, and recovery file downloads.

#### System Toolbox

Provide utilities for device management, firewall, FTP, Fail2Ban, ClamAV, and other system-level operations.

#### AI / Runtime

Support AI-related components such as Ollama, MCP server, GPU monitoring/runtime, and AI account/provider/channel/plugin/security/skill management.

#### Monitoring / Alert / Cron

Track system state, dashboards, processes, alerts, and scheduled jobs.

### Local Databases

The project uses multiple SQLite databases, including:

- `agent.db`
- `task.db`
- `monitor.db`
- `gpu_monitor.db`
- `alert.db`
- `core.db` (optional depending on deployment)

### Logging

The system uses Logrus with file and stdout logging. Logging is initialized early in the startup flow to capture the full application lifecycle.

### Internationalization

The project includes embedded i18n bundles such as:

- `en`
- `zh`
- `zh-Hant`
- `pt-BR`
- `ja`
- `ru`
- `ms`
- `ko`
- `tr`
- `es-ES`

### Cron and Background Jobs

Scheduled tasks include:

- NTP time synchronization
- daily website jobs
- daily SSL jobs
- app store synchronization
- backup token refresh
- restoring enabled cronjobs from the database

### Deployment Notes

- The project integrates tightly with the operating system and system services, so it should run on a real server or a sufficiently privileged test environment.
- For production use, it is recommended to place the agent behind a reverse proxy or configure TLS properly.
- Some features depend directly on external software such as Docker, Nginx, SSH, firewall tools, or antivirus services.

### License

This project is released under the `GPL-3.0` license. See `LICENSE` for details.

### Relevant Modules and Directories

- `cmd/server/main.go`
- `server/server.go`
- `init/router/router.go`
- `middleware/apikey.go`
- `router/`
- `app/api/v2/`
- `utils/`
- `cron/`
- `init/`

If needed, this README can be extended with production installation steps, an `app.yaml` example, architecture diagrams, or detailed API usage examples.
