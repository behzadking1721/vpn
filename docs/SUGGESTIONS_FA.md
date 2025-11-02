# پیشنهادات و چیزهای کم در پروژه VPN Client

## 🔴 بخش‌های اصلی که کاملاً کم هستند

### 1. پیاده‌سازی واقعی پروتکل‌های VPN ⚠️ **اولویت بالا**

**مشکل:** در مستندات گفته شده که 8 پروتکل پشتیبانی می‌شود، اما پیاده‌سازی واقعی وجود ندارد:
- VMess
- VLESS  
- Shadowsocks
- Trojan
- Reality
- Hysteria2
- TUIC
- SSH

**پیشنهاد:**
- استفاده از کتابخانه‌های Go مثل:
  - `github.com/v2fly/v2ray-core/v5` برای VMess/VLESS
  - `github.com/shadowsocks/go-shadowsocks2` برای Shadowsocks
  - `github.com/p4gefau1t/trojan-go` برای Trojan
  - `github.com/apernet/hysteria` برای Hysteria2
- یا استفاده از wrapper های این کتابخانه‌ها

### 2. REST API Server ⚠️ **اولویت بالا**

**مشکل:** پوشه `internal/api` خالی است. UI نیاز به API دارد اما هیچ endpoint ای پیاده‌سازی نشده.

**پیشنهاد:**
- پیاده‌سازی کامل API server با gorilla/mux یا gin
- Endpoint های مورد نیاز:
  ```
  GET    /api/servers           - لیست سرورها
  POST   /api/servers           - افزودن سرور
  GET    /api/servers/{id}      - جزئیات سرور
  PUT    /api/servers/{id}      - به‌روزرسانی سرور
  DELETE /api/servers/{id}      - حذف سرور
  POST   /api/connect           - اتصال
  POST   /api/disconnect        - قطع اتصال
  GET    /api/status            - وضعیت اتصال
  GET    /api/stats             - آمار اتصال
  GET    /api/subscriptions     - لیست اشتراک‌ها
  POST   /api/subscriptions     - افزودن اشتراک
  ```

### 3. Server Manager ⚠️ **اولویت بالا**

**مشکل:** `ServerManager` در مستندات ذکر شده اما پیاده‌سازی ندارد.

**پیشنهاد:**
- پیاده‌سازی کامل ServerManager با قابلیت‌های:
  - افزودن/حذف/ویرایش سرور
  - ذخیره‌سازی در دیتابیس
  - فیلتر کردن سرورها
  - جستجوی سرور
  - مرتب‌سازی بر اساس ping

### 4. Database/Persistence Layer ⚠️ **اولویت بالا**

**مشکل:** پوشه `internal/database` خالی است. هیچ سیستم ذخیره‌سازی داده وجود ندارد.

**پیشنهاد:**
- استفاده از SQLite برای ذخیره‌سازی محلی (برای desktop)
- یا استفاده از JSON file (ساده‌تر اما قابل اعتماد)
- ساختار پیشنهادی:
  ```
  servers.db
  ├── servers (id, name, host, port, protocol, config_json, created_at, updated_at)
  ├── subscriptions (id, name, url, auto_update, last_update, created_at)
  └── settings (key, value)
  ```

### 5. Subscription Parser ⚠️ **اولویت متوسط-بالا**

**مشکل:** فقط UI برای import subscription وجود دارد، اما parsing واقعی نیست.

**پیشنهاد:**
- پیاده‌سازی parser برای:
  - لینک‌های subscription (vmess://, ss://, vless://, ...)
  - Base64 encoded subscription links (مثل v2rayNG format)
  - JSON subscription format
- کتابخانه پیشنهادی: `github.com/v2fly/v2ray-core` برای parse کردن vmess/vless links

### 6. QR Code Parser ⚠️ **اولویت متوسط**

**مشکل:** UI برای QR scanner وجود دارد اما parsing واقعی نیست.

**پیشنهاد:**
- استفاده از کتابخانه‌های QR code:
  - `github.com/skip2/go-qrcode` برای generate
  - `github.com/makiuchi-d/gozxing` برای scan
- Parse کردن محتوای QR به server configuration

---

## 🟡 بخش‌های ناقص که نیاز به تکمیل دارند

### 7. Config Manager ⚠️ **اولویت متوسط**

**مشکل:** فقط یک فایل JSON ساده وجود دارد، مدیریت پیشرفته نیست.

**پیشنهاد:**
- پیاده‌سازی ConfigManager با:
  - Load/Save configuration
  - Validation
  - Default values
  - Environment variable support

### 8. Connection Manager - تکمیل ⚠️ **اولویت بالا**

**مشکل:** فقط یک skeleton ساده وجود دارد:
- اتصال واقعی به VPN انجام نمی‌شود
- Data usage tracking نیست
- Connection statistics نیست
- IPv6 support نیست

**پیشنهاد:**
- اضافه کردن:
  - Real connection handling با protocol handlers
  - Data usage tracking (bytes sent/received)
  - Connection time tracking
  - Speed monitoring
  - Kill switch support

### 9. Ping/Tester برای سرورها ⚠️ **اولویت متوسط**

**مشکل:** Ping measurement فقط در UI شبیه‌سازی شده.

**پیشنهاد:**
- پیاده‌سازی واقعی ping با `github.com/go-ping/ping`
- یا استفاده از TCP connection test
- Latency measurement
- Speed test functionality

### 10. Alert System ⚠️ **اولویت پایین-متوسط**

**مشکل:** فقط models وجود دارد، پیاده‌سازی نیست.

**پیشنهاد:**
- پیاده‌سازی AlertManager با:
  - Rule evaluation
  - Notification system
  - Alert history

---

## 🟢 بهبودهای پیشنهادی (Nice to Have)

### 11. Logging System

**پیشنهاد:**
- استفاده از `github.com/sirupsen/logrus` یا `go.uber.org/zap`
- Log rotation
- Different log levels
- File and console logging

### 12. Testing

**پیشنهاد:**
- افزایش test coverage از 72% به بالای 90%
- Integration tests برای API
- Mock tests برای protocol handlers

### 13. Documentation

**پیشنهاد:**
- API documentation با Swagger/OpenAPI
- Code comments به فارسی/انگلیسی
- User manual

### 14. Security Features

**پیشنهاد:**
- Kill switch (قطع خودکار اینترنت در صورت قطع VPN)
- DNS leak protection
- Firewall rules management
- Split tunneling

### 15. Advanced Features

**پیشنهاد:**
- Auto-connect on startup
- Fastest server auto-selection
- Server health monitoring
- Traffic statistics charts
- Export/Import configuration

---

## 📋 اولویت‌بندی پیشنهادی برای پیاده‌سازی

### Phase 1 (ضروری - 2-3 هفته)
1. ✅ Database/Persistence Layer
2. ✅ Server Manager
3. ✅ REST API Server (حداقل endpoints اصلی)
4. ✅ Connection Manager تکمیل (با real connection)

### Phase 2 (مهم - 3-4 هفته)
5. ✅ پیاده‌سازی حداقل 2-3 پروتکل اصلی (Shadowsocks, VMess, VLESS)
6. ✅ Subscription Parser
7. ✅ Config Manager
8. ✅ Ping/Tester

### Phase 3 (تکمیل - 2-3 هفته)
9. ✅ QR Code Parser
10. ✅ Data usage tracking
11. ✅ Logging System
12. ✅ بهبود تست‌ها

### Phase 4 (پیشرفته - 4-6 هفته)
13. ✅ باقی پروتکل‌ها
14. ✅ Security features (Kill switch, DNS leak protection)
15. ✅ Advanced statistics
16. ✅ UI improvements

---

## 🛠️ پیشنهادات فنی

### ساختار پیشنهادی برای اضافه کردن:

```
internal/
├── api/
│   ├── server.go          # API server setup
│   ├── handlers.go        # HTTP handlers
│   ├── middleware.go      # CORS, logging, etc.
│   └── routes.go          # Route definitions
├── database/
│   ├── db.go              # Database connection
│   ├── server_store.go    # Server CRUD operations
│   └── subscription_store.go
├── managers/
│   ├── server_manager.go  # NEW - Server management logic
│   ├── subscription_manager.go  # NEW
│   └── config_manager.go  # NEW
├── protocols/
│   ├── shadowsocks.go     # NEW - Real implementation
│   ├── vmess.go           # NEW
│   ├── vless.go           # NEW
│   └── ...
└── utils/
    ├── subscription_parser.go  # NEW
    ├── qr_parser.go           # NEW
    └── ping.go                # NEW
```

### Dependencies پیشنهادی برای go.mod:

```go
require (
    github.com/gorilla/mux v1.8.1
    github.com/v2fly/v2ray-core/v5 v5.15.0
    github.com/shadowsocks/go-shadowsocks2 v0.1.5
    github.com/go-ping/ping v1.1.0
    github.com/mattn/go-sqlite3 v1.14.19
    github.com/sirupsen/logrus v1.9.3
    github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
)
```

---

## 💡 نکات مهم

1. **شروع با MVP**: ابتدا یک پروتکل (مثلاً Shadowsocks) را کاملاً پیاده‌سازی کنید و بعد بقیه را اضافه کنید.

2. **API First**: قبل از UI، API را کامل کنید تا frontend بتواند به آن متصل شود.

3. **Database**: SQLite بهترین انتخاب برای desktop application است (نیاز به سرور جداگانه ندارد).

4. **Testing**: با هر feature جدید، تست‌های مربوطه را هم بنویسید.

5. **Documentation**: کد را کامنت کنید و API documentation بنویسید.

