# راهنمای ادغام کتابخانه‌های پروتکل

این سند نحوه ادغام کتابخانه‌های واقعی پروتکل در برنامه کلاینت VPN را توضیح می‌دهد.

## نمای کلی

کلاینت VPN از یک معماری ماژولار استفاده می‌کند که امکان ادغام آسان کتابخانه‌های مختلف پروتکل را فراهم می‌کند. هر پروتکل به عنوان یک handler پیاده‌سازی می‌شود که با رابط [ProtocolHandler](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go#L9-L16) مطابقت دارد.

## افزودن یک پروتکل جدید

برای افزودن پشتیبانی از یک پروتکل جدید، مراحل زیر را دنبال کنید:

### 1. به‌روزرسانی انواع اصلی

پروتکل جدید را به شمارش [ProtocolType](file:///c:/Users/behza/OneDrive/Documents/vpn/src/core/types.go#L9-L16) در [src/core/types.go](file:///c:/Users/behza/OneDrive/Documents/vpn/src/core/types.go) اضافه کنید:

```go
const (
    // ... پروتکل‌های موجود ...
    ProtocolNewProtocol ProtocolType = "newprotocol"
)
```

### 2. به‌روزرسانی Protocol Factory

پروتکل جدید را به [ProtocolFactory](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go#L30-L32) در [src/protocols/protocol_handler.go](file:///c:/Users/behza/OneDrive/Documents/vpn/src/protocols/protocol_handler.go) اضافه کنید:

```go
func (pf *ProtocolFactory) CreateHandler(protocolType core.ProtocolType) (ProtocolHandler, error) {
    switch protocolType {
    // ... موارد موجود ...
    case core.ProtocolNewProtocol:
        return NewNewProtocolHandler(), nil
    default:
        return nil, errors.New("نوع پروتکل پشتیبانی نمی‌شود")
    }
}
```

### 3. ایجاد Protocol Handler

فایل جدید `src/protocols/newprotocol_handler.go` با پیاده‌سازی ایجاد کنید:

```go
package protocols

import (
    "c:/Users/behza/OneDrive/Documents/vpn/src/core"
    // کتابخانه واقعی پروتکل را وارد کنید
)

type NewProtocolHandler struct {
    BaseHandler
    server core.Server
    // فیلدهای خاص پروتکل را اضافه کنید
}

func NewNewProtocolHandler() *NewProtocolHandler {
    handler := &NewProtocolHandler{}
    handler.BaseHandler.protocol = core.ProtocolNewProtocol
    return handler
}

func (nph *NewProtocolHandler) Connect(server core.Server) error {
    // پیاده‌سازی در اینجا
    return nil
}

func (nph *NewProtocolHandler) Disconnect() error {
    // پیاده‌سازی در اینجا
    return nil
}

func (nph *NewProtocolHandler) GetDataUsage() (sent, received int64, err error) {
    // پیاده‌سازی در اینجا
    return 0, 0, nil
}

func (nph *NewProtocolHandler) GetConnectionDetails() (map[string]interface{}, error) {
    // پیاده‌سازی در اینجا
    return nil, nil
}
```

## نمونه‌های ادغام پروتکل

### ادغام Shadowsocks

ادغام Shadowsocks نحوه استفاده از یک کتابخانه خارجی را نشان می‌دهد:

1. افزودن وابستگی به [go.mod](file:///c:/Users/behza/OneDrive/Documents/vpn/go.mod):
```go
require (
    github.com/shadowsocks/go-shadowsocks2 v0.0.0-20230516033142-213602970b87
)
```

2. وارد کردن و استفاده در handler:
```go
import (
    "github.com/shadowsocks/go-shadowsocks2/core"
    "github.com/shadowsocks/go-shadowsocks2/shadowsocks"
)
```

3. مقداردهی اولیه رمزنگار و کلاینت:
```go
cipher, err := core.PickCipher(method, []byte(password))
if err != nil {
    return err
}

client, err := shadowsocks.NewClient(host, cipher)
if err != nil {
    return err
}
```

### ادغام V2Ray/XRay

برای پروتکل‌های VMess/VLESS:

1. افزودن وابستگی:
```go
require (
    github.com/v2ray/v2ray-core v0.0.0-20230601040104-11eba6094346
    // یا برای XRay:
    // github.com/xtls/xray-core v0.0.0-20230601040104-11eba6094346
)
```

2. وارد کردن و استفاده:
```go
import (
    "github.com/v2ray/v2ray-core/app/proxyman"
    "github.com/v2ray/v2ray-core/common/protocol"
    "github.com/v2ray/v2ray-core/common/serial"
    "github.com/v2ray/v2ray-core/proxy/vmess"
    "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
)
```

3. پیکربندی و شروع:
```go
config := &core.Config{
    // جزئیات پیکربندی
}

server, err := core.New(config)
if err != nil {
    return err
}

err = server.Start()
if err != nil {
    return err
}
```

## بهترین شیوه‌ها

### مدیریت خطا

همیشه پیام‌های خطای معنادار ارائه دهید:
```go
if server.Method == "" {
    return fmt.Errorf("روش رمزگذاری برای Shadowsocks ضروری است")
}
```

### مدیریت منابع

اطمینان از پاک‌سازی صحیح منابع:
```go
func (handler *ProtocolHandler) Disconnect() error {
    if handler.client != nil {
        handler.client.Close()
        handler.client = nil
    }
    return nil
}
```

### ردیابی داده

پیاده‌سازی ردیابی استفاده از داده در صورت امکان:
```go
func (handler *ProtocolHandler) GetDataUsage() (sent, received int64, err error) {
    if handler.client == nil {
        return 0, 0, fmt.Errorf("متصل نیست")
    }
    
    sent = handler.client.GetBytesSent()
    received = handler.client.GetBytesReceived()
    return sent, received, nil
}
```

## آزمایش Handlerهای پروتکل

تست‌هایی برای هر handler پروتکل در `protocols_test.go` ایجاد کنید:

```go
func TestNewProtocolHandler(t *testing.T) {
    handler := NewNewProtocolHandler()
    
    // تست حالت اولیه
    if handler.IsConnected() {
        t.Error("انتظار می‌رفت handler در ابتدا قطع باشد")
    }
    
    // تست اتصال
    server := core.Server{
        // پیکربندی سرور
    }
    
    err := handler.Connect(server)
    if err != nil {
        t.Errorf("خطا در اتصال: %v", err)
    }
    
    // تست قطع اتصال
    err = handler.Disconnect()
    if err != nil {
        t.Errorf("خطا در قطع اتصال: %v", err)
    }
}
```

## ملاحظات امنیتی

1. اعتبارسنجی تمام پارامترهای ورودی
2. مدیریت امن داده‌های حساس (رمزها، کلیدها)
3. استفاده از پیش‌فرض‌های امن برای روش‌های رمزگذاری
4. پیاده‌سازی مدیریت خطا بدون نشت اطلاعات حساس

## بهینه‌سازی عملکرد

1. استفاده مجدد از اتصالات در صورت امکان
2. پیاده‌سازی استخر اتصال
3. استفاده از عملیات I/O با بافر
4. به حداقل رساندن تخصیص‌های حافظه

## عیب‌یابی

مسائل رایج و راه‌حل‌ها:

1. **تضاد وابستگی**: از ماژول‌های Go برای مدیریت نسخه‌ها استفاده کنید
2. **خطاهای کامپایل**: مستندات کتابخانه را برای تنظیمات CGO مورد نیاز بررسی کنید
3. **خطاهای اتصال**: پیکربندی سرور و دسترسی شبکه را تأیید کنید
4. **مسائل عملکرد**: برنامه را پروفایل کنید تا گلوگاه‌ها را شناسایی کنید

## مراحل بعدی

1. پیاده‌سازی کتابخانه‌های پروتکل واقعی برای تمام پروتکل‌های پشتیبانی شده
2. افزودن تست‌های جامع واحد
3. ایجاد تست‌های ادغام با سرورهای واقعی
4. مستندسازی گزینه‌های پیکربندی خاص پروتکل
5. پیاده‌سازی مدیریت خطا خاص پروتکل