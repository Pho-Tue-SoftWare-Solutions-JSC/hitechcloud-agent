# HiTechCloud Agent

HiTechCloud Agent là dịch vụ backend viết bằng Go, chạy trực tiếp trên máy chủ để cung cấp API quản trị hạ tầng cho hệ sinh thái HiTechCloud. Dự án tập trung vào quản lý website, ứng dụng, container, backup, công cụ hệ thống và các tính năng AI/runtime mở rộng.

## Tính năng chính

- Quản lý website, Nginx, SSL, ACME, DNS account và hosting account.
- Quản lý Docker, container, image, network, volume, compose và template triển khai.
- Cài đặt, đồng bộ, nâng cấp và quản lý ứng dụng từ app store.
- Backup, restore, upload recovery và kết nối nhiều nhà cung cấp cloud storage.
- Quản lý database, cronjob, process, dashboard, log và monitoring.
- Toolbox hệ thống: firewall, SSH, FTP, Fail2Ban, ClamAV, device/system tools.
- Tích hợp AI tools như Ollama, MCP server, GPU/runtime và AI agent/provider.
- Hỗ trợ WebSocket cho các tác vụ realtime như terminal/process hoặc luồng theo dõi trạng thái.

## Kiến trúc tổng quan

Agent này thường được cài trên từng máy chủ Linux và hoạt động như một lớp điều khiển cục bộ:

- Cung cấp API qua các nhóm route `/api/v1` và `/api/v2`.
- Lưu trữ dữ liệu cục bộ bằng SQLite cho settings, task, monitor, alert và các metadata khác.
- Tích hợp sâu với tài nguyên hệ thống như Docker socket, filesystem, firewall, cron và network.
- Hỗ trợ chạy qua TCP port hoặc Unix socket trong chế độ master.

## Cấu trúc chính

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

## Yêu cầu hệ thống

Khuyến nghị môi trường chạy:

- Linux server.
- Go 1.25.7 hoặc phiên bản tương thích với `go.mod`.
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

## Cấu hình

Cấu hình mặc định được embed từ file:

- `cmd/server/conf/app.yaml`

Ví dụ cấu hình mặc định hiện tại:

- `base.install_dir: /opt`
- `base.mode: dev`
- `log.level: debug`

Trong môi trường dev, agent có thể ưu tiên file cấu hình ngoài tại:

- `/opt/HiTechCloud/conf/app.yaml`

Port server được lấy từ `global.CONF.Base.Port`. Nếu không cấu hình, hệ thống dùng mặc định:

- `8443`

## Build và chạy

### Chạy trực tiếp

```bash
go run ./cmd/server
```

### Build binary

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

## Cơ chế lắng nghe

Agent hỗ trợ 2 chế độ chính:

- **TCP server**: lắng nghe tại `0.0.0.0:<port>`.
- **Master mode**: lắng nghe qua Unix socket `/etc/HiTechCloud/agent.sock`.

Nếu trong settings đã có `ServerCrt` và `ServerKey`, agent sẽ ưu tiên chạy HTTPS với TLS >= 1.2. Nếu không có, hệ thống sẽ fallback sang HTTP.

## Bảo mật và xác thực

- Public endpoint mặc định:
  - `GET /api/v1/health`
- Các route private dưới `/api/v1` và `/api/v2` dùng middleware API key.
- Header xác thực hiện tại:
  - `X-API-Key`
- API key được hash bằng SHA-256 và so sánh theo constant-time.
- Hỗ trợ allowlist IP thông qua setting `ApiKeyAllowedIPs`.

## Các nhóm API chính

### Website / Hosting

Quản lý domain, SSL, ACME, rewrite, redirect, proxy, auth, CORS, anti-leech, PHP, composer và các thành phần hosting liên quan.

### Container / Docker

Quản lý container, image, volume, network, compose, template, log, stats, exec và thao tác file trong container.

### App Store / Ứng dụng

Đồng bộ danh sách app, tìm kiếm, cài đặt, cập nhật, quản lý ứng dụng đã cài và ignore upgrade.

### Backup / Restore

Quản lý cấu hình backup, token, bucket, backup thủ công/định kỳ, restore và tải tệp recovery.

### System Toolbox

Cung cấp các công cụ thao tác với thiết bị, firewall, FTP, Fail2Ban, ClamAV và các tác vụ hệ thống khác.

### AI / Runtime

Hỗ trợ các thành phần AI như Ollama, MCP server, GPU monitoring/runtime, AI account/provider/channel/plugin/security/skill.

### Monitoring / Alert / Cron

Theo dõi trạng thái hệ thống, dashboard, tiến trình, cảnh báo và tác vụ định kỳ.

## Database cục bộ

Dự án sử dụng SQLite cho nhiều miền dữ liệu khác nhau, bao gồm:

- `agent.db`
- `task.db`
- `monitor.db`
- `gpu_monitor.db`
- `alert.db`
- `core.db` (tùy trường hợp)

## Logging

Hệ thống sử dụng Logrus kết hợp ghi log file và stdout. Thành phần logging được khởi tạo sớm trong quá trình start để hỗ trợ theo dõi toàn bộ vòng đời ứng dụng.

## Đa ngôn ngữ

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

## Cron và tác vụ nền

Hệ thống có các tác vụ định kỳ như:

- đồng bộ thời gian NTP
- job website hằng ngày
- job SSL hằng ngày
- đồng bộ app store
- refresh backup token
- khôi phục cronjob đã bật từ database

## Ghi chú triển khai

- Dự án tích hợp mạnh với hệ điều hành và dịch vụ hệ thống, vì vậy nên chạy trong môi trường server thật hoặc môi trường test có đủ quyền cần thiết.
- Nếu chạy production, nên đặt agent sau reverse proxy hoặc cấu hình TLS đầy đủ.
- Một số tính năng phụ thuộc trực tiếp vào phần mềm ngoài như Docker, Nginx, SSH, firewall tools hoặc antivirus service.

## License

Dự án được phát hành theo giấy phép `GPL-3.0`. Xem thêm tại file `LICENSE`.

## Module và thư mục liên quan

- `cmd/server/main.go`
- `server/server.go`
- `init/router/router.go`
- `middleware/apikey.go`
- `router/`
- `app/api/v2/`
- `utils/`
- `cron/`
- `init/`

Nếu cần, có thể bổ sung thêm các phần như hướng dẫn cài đặt production, ví dụ cấu hình `app.yaml`, sơ đồ kiến trúc hoặc tài liệu API chi tiết cho từng nhóm route.
