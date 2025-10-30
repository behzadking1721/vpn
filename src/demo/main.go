package main

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"c:/Users/behza/OneDrive/Documents/vpn/src/managers"
	"c:/Users/behza/OneDrive/Documents/vpn/src/protocols"
	"c:/Users/behza/OneDrive/Documents/vpn/src/utils"
	"fmt"
	"time"
)

func main() {
	fmt.Println("دموی کلاینت VPN چندپروتکل")
	fmt.Println("========================")

	// ایجاد مدیر داده
	dataManager := managers.NewDataManager("./data/usage.json")

	// ایجاد مدیر سرور
	serverManager := managers.NewServerManager()

	// ایجاد سرورهای نمونه
	vmessServer := core.Server{
		ID:         utils.GenerateID(),
		Name:       "سرور VMess نمونه",
		Host:       "vmess.example.com",
		Port:       443,
		Protocol:   core.ProtocolVMess,
		Encryption: "auto",
		Password:   "test-uuid",
		TLS:        true,
		Remark:     "سرور VMess برای تست",
		Enabled:    true,
	}

	shadowsocksServer := core.Server{
		ID:       utils.GenerateID(),
		Name:     "سرور Shadowsocks نمونه",
		Host:     "ss.example.com",
		Port:     8388,
		Protocol: core.ProtocolShadowsocks,
		Method:   "aes-256-gcm",
		Password: "test-password",
		Remark:   "سرور Shadowsocks برای تست",
		Enabled:  true,
	}

	// افزودن سرورها به مدیر
	serverManager.AddServer(vmessServer)
	serverManager.AddServer(shadowsocksServer)

	// نمایش لیست سرورها
	fmt.Println("\nسرورهای پیکربندی شده:")
	servers := serverManager.GetAllServers()
	for _, server := range servers {
		fmt.Printf("- %s (%s) - %s:%d\n", server.Name, server.Protocol, server.Host, server.Port)
	}

	// تست پردازش لینک اشتراک
	fmt.Println("\nدر حال پردازش لینک اشتراک...")
	parser := managers.NewSubscriptionParser()

	// لینک VMess نمونه (رمزشده با Base64)
	vmessJSON := `{"v":"2","ps":"سرور تست VMess","add":"test.vmess.com","port":"443","id":"test-uuid","aid":"0","scy":"auto","net":"tcp","type":"none","host":"","path":"","tls":"tls"}`
	vmessLink := "vmess://" + utils.Base64Encode([]byte(vmessJSON))

	parsedServers, err := parser.ParseSubscriptionLink(vmessLink)
	if err != nil {
		fmt.Printf("خطا در پردازش لینک: %v\n", err)
	} else {
		fmt.Printf("پردازش موفقیت‌آمیز. %d سرور استخراج شد:\n", len(parsedServers))
		for _, server := range parsedServers {
			fmt.Printf("- %s (%s)\n", server.Name, server.Host)
		}
	}

	// تست پروتکل‌ها
	fmt.Println("\nدر حال تست پروتکل‌ها...")
	protocolFactory := protocols.NewProtocolFactory()

	for _, protocolType := range []core.ProtocolType{
		core.ProtocolVMess,
		core.ProtocolShadowsocks,
	} {
		fmt.Printf("\nتست پروتکل %s:\n", protocolType)
		
		handler, err := protocolFactory.CreateHandler(protocolType)
		if err != nil {
			fmt.Printf("  خطا در ایجاد handler: %v\n", err)
			continue
		}

		// انتخاب سرور مناسب
		var testServer core.Server
		switch protocolType {
		case core.ProtocolVMess:
			testServer = vmessServer
		case core.ProtocolShadowsocks:
			testServer = shadowsocksServer
		}

		// اتصال
		fmt.Printf("  در حال اتصال به %s...\n", testServer.Name)
		err = handler.Connect(testServer)
		if err != nil {
			fmt.Printf("  خطا در اتصال: %v\n", err)
			continue
		}

		fmt.Printf("  اتصال موفقیت‌آمیز\n")

		// شبیه‌سازی انتقال داده
		time.Sleep(2 * time.Second)
		
		// دریافت آمار مصرف
		sent, received, err := handler.GetDataUsage()
		if err != nil {
			fmt.Printf("  خطا در دریافت آمار: %v\n", err)
		} else {
			fmt.Printf("  داده ارسالی: %s\n", utils.FormatBytes(sent))
			fmt.Printf("  داده دریافتی: %s\n", utils.FormatBytes(received))
			
			// ثبت داده در مدیر داده
			dataManager.RecordDataUsage(testServer.ID, sent, received)
		}

		// جزئیات اتصال
		details, err := handler.GetConnectionDetails()
		if err != nil {
			fmt.Printf("  خطا در دریافت جزئیات: %v\n", err)
		} else {
			fmt.Printf("  جزئیات اتصال: %v\n", details)
		}

		// قطع اتصال
		fmt.Printf("  در حال قطع اتصال...\n")
		err = handler.Disconnect()
		if err != nil {
			fmt.Printf("  خطا در قطع اتصال: %v\n", err)
		} else {
			fmt.Printf("  اتصال قطع شد\n")
		}
	}

	// نمایش آمار مصرفی
	fmt.Println("\nآمار مصرفی:")
	allData := dataManager.GetAllData()
	for serverID, data := range allData {
		fmt.Printf("سرور %s: ارسال %s، دریافت %s\n", 
			serverID[:8], 
			utils.FormatBytes(data.TotalSent), 
			utils.FormatBytes(data.TotalRecv))
	}

	fmt.Println("\nدمو به پایان رسید!")
}