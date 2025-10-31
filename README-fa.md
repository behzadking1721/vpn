# کلاینت چندپروتکلی VPN

یک برنامه VPN چندسکویی با رابط کاربری ساده و کاربرپسند مشابه Hiddify، با پشتیبانی از پروتکل‌های مختلف و لینک‌های اشتراک.

## ویژگی‌ها

- پشتیبانی چندسکویی (Android, iOS, Windows, macOS, Linux)
- پشتیبانی از پروتکل‌های متعدد:
  - VMess
  - VLESS
  - Trojan
  - Reality
  - Hysteria2
  - TUIC
  - SSH
  - Shadowsocks
- قابلیت مدیریت سرورها
- مدیریت تنظیمات کاربر و اتصال
- انتخاب خودکار سریع‌ترین سرور (LowestPing)
- نمایش مصرف داده
- وارد کردن لینک اشتراک با پیوند عمیق
- وارد کردن با کد QR
- پشتیبانی از IPv6
- تجربه بدون تبلیغات
- متن‌باز و ایمن

## پیش‌نیازها

قبل از اینکه بتوانید این پروژه را بسازید و اجرا کنید، باید وابستگی‌های زیر را نصب کنید:

### Go (Golang)

- **Windows**: دانلود از [golang.org](https://golang.org/dl/)
- **macOS**: `brew install go` یا دانلود از [golang.org](https://golang.org/dl/)
- **Linux**: `sudo apt install golang` (Ubuntu/Debian) یا معادل برای توزیع شما

## راه‌اندازی و تست

### 1. دریافت کد منبع

```bash
git clone <آدرس مخزن Git>
cd vpn
```

### 2. نصب وابستگی‌ها

```bash
go mod tidy
```

### 3. ساخت برنامه

برای Windows:
```bash
go build -o vpn-client.exe ./src
```

برای Linux/macOS:
```bash
go build -o vpn-client ./src
```

### 4. اجرای برنامه

```bash
# Windows
.\vpn-client.exe --version

# Linux/macOS
./vpn-client --version
```

## استفاده

برنامه از حالت‌های مختلفی پشتیبانی می‌کند:

### حالت رابط خط فرمان (CLI)
```bash
go run src/main.go --cli
```

### حالت سرور API
```bash
go run src/main.go --api
```

سپس مرورگر خود را باز کنید و به http://localhost:8080 مراجعه کنید.

### راهنمای استفاده

```bash
go run src/main.go --help
```

## توسعه

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

## ساختار پروژه

```
vpn/
├── src/
│   ├── api/           # API REST برای ادغام رابط کاربری
│   ├── cli/           # رابط خط فرمان
│   ├── core/          # مدل‌های داده و رابط‌های هسته
│   ├── managers/      # منطق تجاری مدیران
│   ├── protocols/     # پیاده‌سازی‌های خاص پروتکل
│   ├── utils/         # توابع کمکی
│   └── main.go        # نقطه ورود برنامه
├── docs/              # مستندات
├── assets/            # تصاویر و دیگر دارایی‌ها
└── scripts/           # اسکریپت‌های ساخت و بسته‌بندی
```

## مجوز

این پروژه تحت مجوز MIT منتشر شده است - برای جزئیات به فایل [LICENSE](LICENSE) مراجعه کنید.