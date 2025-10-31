package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"encoding/base64"
	"testing"
)

func TestSubscriptionParser(t *testing.T) {
	parser := NewSubscriptionParser()

	// تست لینک VMess
	t.Run("ParseVMessLink", func(t *testing.T) {
		// لینک VMess نمونه (رمزشده با Base64)
		vmessJSON := `{"v":"2","ps":"Test Server","add":"example.com","port":"443","id":"test-uuid","aid":"0","scy":"auto","net":"tcp","type":"none","host":"","path":"","tls":"tls"}`
		vmessLink := "vmess://" + base64.StdEncoding.EncodeToString([]byte(vmessJSON))

		servers, err := parser.ParseSubscriptionLink(vmessLink)
		if err != nil {
			t.Errorf("خطا در تجزیه لینک VMess: %v", err)
			return
		}

		if len(servers) != 1 {
			t.Errorf("انتظار 1 سرور داشتیم، دریافت کردیم: %d", len(servers))
			return
		}

		server := servers[0]
		if server.Protocol != core.ProtocolVMess {
			t.Errorf("انتظار پروتکل VMess داشتیم، دریافت کردیم: %s", server.Protocol)
		}

		if server.Host != "example.com" {
			t.Errorf("انتظار هاست example.com داشتیم، دریافت کردیم: %s", server.Host)
		}

		if server.Name != "Test Server" {
			t.Errorf("انتظار نام Test Server داشتیم، دریافت کردیم: %s", server.Name)
		}
	})

	// تست لینک Shadowsocks
	t.Run("ParseShadowsocksLink", func(t *testing.T) {
		// لینک Shadowsocks نمونه
		ssInfo := "aes-256-gcm:password@example.com:8388"
		encodedInfo := base64.StdEncoding.EncodeToString([]byte(ssInfo))
		ssLink := "ss://" + encodedInfo + "#Test%20SS%20Server"

		servers, err := parser.ParseSubscriptionLink(ssLink)
		if err != nil {
			t.Errorf("خطا در تجزیه لینک Shadowsocks: %v", err)
			return
		}

		if len(servers) != 1 {
			t.Errorf("انتظار 1 سرور داشتیم، دریافت کردیم: %d", len(servers))
			return
		}

		server := servers[0]
		if server.Protocol != core.ProtocolShadowsocks {
			t.Errorf("انتظار پروتکل Shadowsocks داشتیم، دریافت کردیم: %s", server.Protocol)
		}

		if server.Host != "example.com" {
			t.Errorf("انتظار هاست example.com داشتیم، دریافت کردیم: %s", server.Host)
		}

		if server.Method != "aes-256-gcm" {
			t.Errorf("انتظار روش aes-256-gcm داشتیم، دریافت کردیم: %s", server.Method)
		}
	})

	// تست لینک Trojan
	t.Run("ParseTrojanLink", func(t *testing.T) {
		// لینک Trojan نمونه
		trojanLink := "trojan://password@example.com:443#Test%20Trojan%20Server"

		servers, err := parser.ParseSubscriptionLink(trojanLink)
		if err != nil {
			t.Errorf("خطا در تجزیه لینک Trojan: %v", err)
			return
		}

		if len(servers) != 1 {
			t.Errorf("انتظار 1 سرور داشتیم، دریافت کردیم: %d", len(servers))
			return
		}

		server := servers[0]
		if server.Protocol != core.ProtocolTrojan {
			t.Errorf("انتظار پروتکل Trojan داشتیم، دریافت کردیم: %s", server.Protocol)
		}

		if server.Host != "example.com" {
			t.Errorf("انتظار هاست example.com داشتیم، دریافت کردیم: %s", server.Host)
		}

		if server.Password != "password" {
			t.Errorf("انتظار رمز عبور password داشتیم، دریافت کردیم: %s", server.Password)
		}
	})

	// تست لینک اشتراک Base64
	t.Run("ParseBase64Subscription", func(t *testing.T) {
		// ایجاد محتوای اشتراک با چند لینک
		subContent := "vmess://" + base64.StdEncoding.EncodeToString([]byte(`{"v":"2","ps":"VMess Server","add":"vmess.example.com","port":"443","id":"test-uuid","aid":"0","scy":"auto","net":"tcp","type":"none","host":"","path":"","tls":"tls"}`)) + "\n" +
			"ss://" + base64.StdEncoding.EncodeToString([]byte("aes-256-gcm:password@ss.example.com:8388")) + "\n" +
			"trojan://password@trojan.example.com:443"

		encodedSub := base64.StdEncoding.EncodeToString([]byte(subContent))

		servers, err := parser.ParseSubscriptionLink(encodedSub)
		if err != nil {
			t.Errorf("خطا در تجزیه لینک اشتراک: %v", err)
			return
		}

		if len(servers) != 3 {
			t.Errorf("انتظار 3 سرور داشتیم، دریافت کردیم: %d", len(servers))
			return
		}

		// بررسی نوع پروتکل‌ها
		protocols := make(map[core.ProtocolType]bool)
		for _, server := range servers {
			protocols[server.Protocol] = true
		}

		if !protocols[core.ProtocolVMess] {
			t.Error("سرور VMess یافت نشد")
		}

		if !protocols[core.ProtocolShadowsocks] {
			t.Error("سرور Shadowsocks یافت نشد")
		}

		if !protocols[core.ProtocolTrojan] {
			t.Error("سرور Trojan یافت نشد")
		}
	})
}
