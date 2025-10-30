package managers

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// SubscriptionParser پردازش‌گر لینک‌های اشتراک
type SubscriptionParser struct{}

// NewSubscriptionParser ایجاد یک پردازش‌گر جدید
func NewSubscriptionParser() *SubscriptionParser {
	return &SubscriptionParser{}
}

// ParseSubscriptionLink تجزیه لینک اشتراک و استخراج سرورها
func (sp *SubscriptionParser) ParseSubscriptionLink(subLink string) ([]core.Server, error) {
	// حذف پیشوند اگر وجود دارد
	link := strings.TrimSpace(subLink)
	
	// پشتیبانی از فرمت‌های مختلف
	if strings.HasPrefix(link, "vmess://") {
		return sp.parseVMessLink(link)
	} else if strings.HasPrefix(link, "ss://") {
		return sp.parseShadowsocksLink(link)
	} else if strings.HasPrefix(link, "trojan://") {
		return sp.parseTrojanLink(link)
	} else {
		// فرض کنید لینک یک لینک اشتراک Base64 است
		return sp.parseBase64Subscription(link)
	}
}

// parseVMessLink تجزیه لینک VMess
func (sp *SubscriptionParser) parseVMessLink(link string) ([]core.Server, error) {
	// حذف پیشوند
	base64Str := strings.TrimPrefix(link, "vmess://")
	
	// رمزگشایی Base64
	jsonStr, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("خطا در رمزگشایی لینک VMess: %v", err)
	}
	
	// تجزیه JSON
	var vmessConfig map[string]interface{}
	err = json.Unmarshal(jsonStr, &vmessConfig)
	if err != nil {
		return nil, fmt.Errorf("خطا در تجزیه JSON VMess: %v", err)
	}
	
	// ایجاد سرور از پیکربندی
	server := core.Server{
		ID:         utils.GenerateID(),
		Name:       getStringValue(vmessConfig, "ps", "VMess Server"),
		Host:       getStringValue(vmessConfig, "add", ""),
		Port:       getIntValue(vmessConfig, "port", 0),
		Protocol:   core.ProtocolVMess,
		Encryption: getStringValue(vmessConfig, "scy", "auto"),
		Password:   getStringValue(vmessConfig, "id", ""),
		TLS:        getStringValue(vmessConfig, "tls", "") == "tls",
		SNI:        getStringValue(vmessConfig, "sni", ""),
	}
	
	// تنظیم نام سرور اگر خالی باشد
	if server.Name == "" {
		server.Name = fmt.Sprintf("VMess %s:%d", server.Host, server.Port)
	}
	
	// فعال کردن سرور به طور پیش‌فرض
	server.Enabled = true
	
	return []core.Server{server}, nil
}

// parseShadowsocksLink تجزیه لینک Shadowsocks
func (sp *SubscriptionParser) parseShadowsocksLink(link string) ([]core.Server, error) {
	// حذف پیشوند
	base64Part := strings.TrimPrefix(link, "ss://")
	
	// تقسیم به قسمت‌های اطلاعات و سرور
	parts := strings.Split(base64Part, "#")
	
	// رمزگشایی قسمت اطلاعات
	decodedInfo, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		// امتحان کدگشایی URL-Safe
		decodedInfo, err = base64.URLEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, fmt.Errorf("خطا در رمزگشایی لینک Shadowsocks: %v", err)
		}
	}
	
	// تجزیه اطلاعات
	infoStr := string(decodedInfo)
	atIndex := strings.Index(infoStr, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("فرمت لینک Shadowsocks نامعتبر است")
	}
	
	authInfo := infoStr[:atIndex]
	hostInfo := infoStr[atIndex+1:]
	
	// تجزیه اطلاعات احراز هویت
	colonIndex := strings.Index(authInfo, ":")
	if colonIndex == -1 {
		return nil, fmt.Errorf("فرمت اطلاعات احراز هویت Shadowsocks نامعتبر است")
	}
	
	method := authInfo[:colonIndex]
	password := authInfo[colonIndex+1:]
	
	// تجزیه اطلاعات سرور
	hostParts := strings.Split(hostInfo, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("فرمت آدرس سرور Shadowsocks نامعتبر است")
	}
	
	host := hostParts[0]
	var port int
	// تبدیل پورت به عدد صحیح
	fmt.Sscanf(hostParts[1], "%d", &port)
	
	// استخراج نام سرور از هش‌تگ
	serverName := "Shadowsocks Server"
	if len(parts) > 1 {
		decodedName, err := base64.URLEncoding.DecodeString(parts[1])
		if err == nil {
			serverName = string(decodedName)
		} else {
			serverName = parts[1]
		}
	}
	
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     serverName,
		Host:     host,
		Port:     port,
		Protocol: core.ProtocolShadowsocks,
		Method:   method,
		Password: password,
		Enabled:  true,
	}
	
	return []core.Server{server}, nil
}

// parseTrojanLink تجزیه لینک Trojan
func (sp *SubscriptionParser) parseTrojanLink(link string) ([]core.Server, error) {
	// حذف پیشوند
	linkContent := strings.TrimPrefix(link, "trojan://")
	
	// پیدا کردن اسلش برای جدا کردن پارامترها
	slashIndex := strings.Index(linkContent, "/")
	if slashIndex == -1 {
		slashIndex = len(linkContent)
	}
	
	// استخراج اطلاعات اصلی
	mainPart := linkContent[:slashIndex]
	
	// پیدا کردن @ برای جدا کردن رمز عبور و سرور
	atIndex := strings.Index(mainPart, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("فرمت لینک Trojan نامعتبر است")
	}
	
	password := mainPart[:atIndex]
	serverInfo := mainPart[atIndex+1:]
	
	// تجزیه اطلاعات سرور
	colonIndex := strings.LastIndex(serverInfo, ":")
	if colonIndex == -1 {
		return nil, fmt.Errorf("فرمت آدرس سرور Trojan نامعتبر است")
	}
	
	host := serverInfo[:colonIndex]
	var port int
	fmt.Sscanf(serverInfo[colonIndex+1:], "%d", &port)
	
	// استخراج نام سرور از پارامترها
	serverName := "Trojan Server"
	if slashIndex < len(linkContent) {
		paramsPart := linkContent[slashIndex+1:]
		if strings.Contains(paramsPart, "#") {
			namePart := strings.Split(paramsPart, "#")[1]
			decodedName, err := base64.URLEncoding.DecodeString(namePart)
			if err == nil {
				serverName = string(decodedName)
			} else {
				serverName = namePart
			}
		}
	}
	
	server := core.Server{
		ID:       utils.GenerateID(),
		Name:     serverName,
		Host:     host,
		Port:     port,
		Protocol: core.ProtocolTrojan,
		Password: password,
		Enabled:  true,
	}
	
	return []core.Server{server}, nil
}

// parseBase64Subscription تجزیه لینک اشتراک Base64
func (sp *SubscriptionParser) parseBase64Subscription(base64Content string) ([]core.Server, error) {
	// رمزگشایی محتوای Base64
	decodedContent, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, fmt.Errorf("خطا در رمزگشایی محتوای اشتراک: %v", err)
	}
	
	// تقسیم به خطوط
	links := strings.Split(string(decodedContent), "\n")
	
	var servers []core.Server
	
	// پردازش هر لینک
	for _, link := range links {
		link = strings.TrimSpace(link)
		if link == "" {
			continue
		}
		
		// تجزیه لینک
		linkServers, err := sp.ParseSubscriptionLink(link)
		if err != nil {
			// ادامه دادن با لینک‌های بعدی در صورت خطا
			continue
		}
		
		servers = append(servers, linkServers...)
	}
	
	return servers, nil
}

// getStringValue دریافت مقدار رشته‌ای از map با مقدار پیش‌فرض
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getIntValue دریافت مقدار عددی از map با مقدار پیش‌فرض
func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if str, ok := val.(string); ok {
			var result int
			fmt.Sscanf(str, "%d", &result)
			return result
		}
	}
	return defaultValue
}