# راهنمای توسعه برنامه VPN Client

این راهنما اطلاعات لازم برای توسعه و گسترش برنامه VPN Client را فراهم می‌کند.

## ساختار پروژه

```
vpn/
├── src/                    # کد منبع اصلی
│   ├── api/                # API REST
│   ├── cli/                # رابط خط فرمان
│   ├── core/               # مدل‌ها و رابط‌های هسته
│   ├── managers/           # منطق تجاری
│   ├── protocols/          # پیاده‌سازی پروتکل‌ها
│   ├── utils/              # ابزارهای کمکی
│   └── main.go             # نقطه ورود اصلی
├── docs/                   # مستندات
├── scripts/                # اسکریپت‌های ساخت
├── assets/                 # منابع گرافیکی
├── packages/               # بسته‌های ساخته شده
├── go.mod                  # وابستگی‌های Go
└── go.sum                  # چک‌سام وابستگی‌ها
```

## افزودن پروتکل جدید

برای افزودن یک پروتکل جدید:

1. ایجاد یک فایل جدید در `src/protocols/` با نام `[protocol]_handler.go`
2. پیاده‌سازی رابط `ProtocolHandler`:
   ```go
   type ProtocolHandler interface {
       Connect(config *core.ServerConfig) error
       Disconnect() error
       IsConnected() bool
       GetStats() *core.ConnectionStats
   }
   ```
3. ثبت پروتکل در `ProtocolManager`

## ساخت برنامه

### ساخت برای تمام سکوها

```bash
chmod +x scripts/build-all.sh
./scripts/build-all.sh
```

### ساخت برای سکوی خاص

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/vpn-client ./src

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/vpn-client.exe ./src

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/vpn-client ./src
```

## تست

اجرای تمام تست‌ها:
```bash
go test ./...
```

اجرای تست‌های خاص:
```bash
go test ./src/protocols/ -v
```

## کدنویسی

### استانداردهای کد

- استفاده از camelCase برای نام‌گذاری متغیرها
- نوشتن کامنت برای توابع عمومی
- نگه داشتن توابع کوچک و تک‌منظوره

### ابزارهای کد

برای فرمت‌دهی کد:
```bash
go fmt ./...
```

برای بررسی خطاها:
```bash
go vet ./...
```

## مدیریت وابستگی‌ها

### افزودن وابستگی جدید

```bash
go get github.com/example/package
```

### به‌روزرسانی وابستگی‌ها

```bash
go get -u ./...
```

### پاک‌سازی وابستگی‌های استفاده نشده

```bash
go mod tidy
```

## ایجاد نسخه جدید

1. به‌روزرسانی نسخه در `versioninfo.json`
2. اضافه کردن تغییرات به `CHANGELOG.md`
3. ایجاد تگ Git:
   ```bash
   git tag -a v1.x.x -m "Release version 1.x.x"
   git push origin v1.x.x
   ```
4. اجرای اسکریپت انتشار:
   ```bash
   ./scripts/release.sh v1.x.x
   ```

## بهینه‌سازی

### پروفایل کردن

برای پروفایل کردن برنامه:
```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

### کاهش اندازه باینری

استفاده از فلگ‌های زیر در زمان ساخت:
```bash
go build -ldflags="-s -w" -o vpn-client ./src
```

## امنیت

### بررسی آسیب‌پذیری‌ها

```bash
govulncheck ./...
```

### بهترین شیوه‌ها

1. اعتبارسنجی تمام ورودی‌ها
2. استفاده از TLS برای ارتباطات شبکه
3. ذخیره امن کلمات عبور (با رمزنگاری)
4. به‌روز نگه داشتن وابستگی‌ها

## مستندسازی

### به‌روزرسانی مستندات

تمام مستندات در پوشه `docs/` قرار دارند. هنگام افزودن/تغییر عملکردها، مستندات مربوطه را نیز به‌روزرسانی کنید.

### ایجاد مستندات API

برای ایجاد مستندات API از Swagger استفاده می‌شود. فایل‌های Swagger در `docs/api/` قرار دارند.