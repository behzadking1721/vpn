package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"time"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/shadowsocks"
)

// LibShadowsocksHandler پیاده‌سازی Shadowsocks با استفاده از کتابخانه shadowsocks-go
type LibShadowsocksHandler struct {
	BaseHandler
	server     core.Server
	connected  bool
	startTime  time.Time
	client     *shadowsocks.Client
}

// NewLibShadowsocksHandler ایجاد یک Shadowsocks handler جدید
func NewLibShadowsocksHandler() *LibShadowsocksHandler {
	handler := &LibShadowsocksHandler{}
	handler.BaseHandler.protocol = core.ProtocolShadowsocks
	return handler
}

// Connect برقراری اتصال به سرور Shadowsocks
func (lsh *LibShadowsocksHandler) Connect(server core.Server) error {
	// ذخیره اطلاعات سرور
	lsh.server = server
	lsh.startTime = time.Now()
	
	fmt.Printf("در حال اتصال به سرور Shadowsocks: %s:%d\n", server.Host, server.Port)
	fmt.Printf("روش رمزگذاری: %s, رمز عبور: %s\n", server.Method, server.Password)
	
	// اعتبارسنجی پارامترهای ضروری
	if server.Method == "" {
		return fmt.Errorf("روش رمزگذاری ضروری است")
	}
	
	if server.Password == "" {
		return fmt.Errorf("رمز عبور ضروری است")
	}
	
	// ایجاد cipher
	cipher, err := core.PickCipher(server.Method, []byte(server.Password))
	if err != nil {
		return fmt.Errorf("خطا در ایجاد cipher: %v", err)
	}
	
	// ایجاد آدرس سرور
	host := fmt.Sprintf("%s:%d", server.Host, server.Port)
	
	// ایجاد کلاینت Shadowsocks
	lsh.client, err = shadowsocks.NewClient(host, cipher)
	if err != nil {
		return fmt.Errorf("خطا در ایجاد کلاینت Shadowsocks: %v", err)
	}
	
	// شروع اتصال
	err = lsh.client.Connect()
	if err != nil {
		return fmt.Errorf("خطا در برقراری اتصال: %v", err)
	}
	
	// علامت‌گذاری به عنوان متصل
	lsh.connected = true
	fmt.Println("اتصال Shadowsocks برقرار شد")
	
	return nil
}

// Disconnect قطع اتصال از سرور Shadowsocks
func (lsh *LibShadowsocksHandler) Disconnect() error {
	if !lsh.connected {
		return fmt.Errorf("به سرور Shadowsocks متصل نیستید")
	}
	
	fmt.Printf("در حال قطع اتصال از سرور Shadowsocks: %s:%d\n", lsh.server.Host, lsh.server.Port)
	
	// بستن کلاینت Shadowsocks
	if lsh.client != nil {
		lsh.client.Close()
		lsh.client = nil
	}
	
	// علامت‌گذاری به عنوان قطع شده
	lsh.connected = false
	lsh.server = core.Server{} // پاک کردن اطلاعات سرور
	
	fmt.Println("اتصال Shadowsocks قطع شد")
	
	return nil
}

// GetDataUsage دریافت میزان داده ارسالی و دریافتی
func (lsh *LibShadowsocksHandler) GetDataUsage() (sent, received int64, err error) {
	if !lsh.connected {
		return 0, 0, fmt.Errorf("به سرور Shadowsocks متصل نیستید")
	}
	
	// دریافت آمار از کلاینت Shadowsocks
	if lsh.client != nil {
		sent = lsh.client.GetBytesSent()
		received = lsh.client.GetBytesReceived()
	} else {
		// شبیه‌سازی مصرف داده
		duration := time.Since(lsh.startTime).Seconds()
		sent = int64(1024 * duration)     // 1 KB/s
		received = int64(5120 * duration) // 5 KB/s
	}
	
	return sent, received, nil
}

// GetConnectionDetails دریافت جزئیات اتصال
func (lsh *LibShadowsocksHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !lsh.connected {
		return nil, fmt.Errorf("به سرور Shadowsocks متصل نیستید")
	}
	
	details := map[string]interface{}{
		"protocol":   "Shadowsocks",
		"host":       lsh.server.Host,
		"port":       lsh.server.Port,
		"method":     lsh.server.Method,
		"connected":  lsh.connected,
		"start_time": lsh.startTime,
	}
	
	// اگر کلاینت موجود باشد، آمار را نیز اضافه کن
	if lsh.client != nil {
		details["bytes_sent"] = lsh.client.GetBytesSent()
		details["bytes_received"] = lsh.client.GetBytesReceived()
	}
	
	return details, nil
}