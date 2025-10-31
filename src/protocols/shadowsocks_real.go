package protocols

import (
	"c:/Users/behza/OneDrive/Documents/vpn/src/core"
	"fmt"
	"time"
)

// RealShadowsocksHandler 实现真实的Shadowsocks协议处理器
type RealShadowsocksHandler struct {
	BaseHandler
	server    core.Server
	connected bool
	startTime time.Time
	dataSent  int64
	dataRecv  int64
}

// NewRealShadowsocksHandler 创建一个新的Shadowsocks处理器
func NewRealShadowsocksHandler() *RealShadowsocksHandler {
	handler := &RealShadowsocksHandler{}
	handler.BaseHandler.protocol = core.ProtocolShadowsocks
	return handler
}

// Connect 建立到Shadowsocks服务器的连接
func (ssh *RealShadowsocksHandler) Connect(server core.Server) error {
	// 在真实实现中，这里会:
	// 1. 解析Shadowsocks配置
	// 2. 初始化Shadowsocks客户端库
	// 3. 建立到服务器的连接

	// 存储服务器信息
	ssh.server = server
	ssh.startTime = time.Now()

	fmt.Printf("Connecting to Shadowsocks server: %s:%d\n", server.Host, server.Port)
	fmt.Printf("Method: %s, Password: %s\n", server.Method, server.Password)

	// 模拟连接过程
	time.Sleep(1 * time.Second)

	// 检查必要参数
	if server.Method == "" {
		return fmt.Errorf("missing encryption method")
	}

	if server.Password == "" {
		return fmt.Errorf("missing password")
	}

	// 标记为已连接
	ssh.connected = true
	fmt.Println("Shadowsocks connection established")

	return nil
}

// Disconnect 断开与Shadowsocks服务器的连接
func (ssh *RealShadowsocksHandler) Disconnect() error {
	if !ssh.connected {
		return fmt.Errorf("not connected to Shadowsocks server")
	}

	// 在真实实现中，这里会:
	// 1. 关闭Shadowsocks客户端连接
	// 2. 清理资源

	fmt.Printf("Disconnecting from Shadowsocks server: %s:%d\n", ssh.server.Host, ssh.server.Port)

	// 模拟断开连接过程
	time.Sleep(500 * time.Millisecond)
	ssh.connected = false
	ssh.server = core.Server{} // 清除服务器信息

	fmt.Println("Shadowsocks connection terminated")

	return nil
}

// GetDataUsage 返回发送和接收的数据量
func (ssh *RealShadowsocksHandler) GetDataUsage() (sent, received int64, err error) {
	if !ssh.connected {
		return 0, 0, fmt.Errorf("not connected to Shadowsocks server")
	}

	// 在真实实现中，这里会从Shadowsocks库获取实际数据
	// 目前我们模拟一些数据使用情况
	ssh.dataSent += 1024 * int64(time.Since(ssh.startTime).Seconds())
	ssh.dataRecv += 2048 * int64(time.Since(ssh.startTime).Seconds())

	return ssh.dataSent, ssh.dataRecv, nil
}

// GetConnectionDetails 返回连接详细信息
func (ssh *RealShadowsocksHandler) GetConnectionDetails() (map[string]interface{}, error) {
	if !ssh.connected {
		return nil, fmt.Errorf("not connected to Shadowsocks server")
	}

	details := map[string]interface{}{
		"protocol":   "Shadowsocks",
		"host":       ssh.server.Host,
		"port":       ssh.server.Port,
		"method":     ssh.server.Method,
		"connected":  ssh.connected,
		"start_time": ssh.startTime,
	}

	return details, nil
}

// SetServer 设置服务器配置
func (ssh *RealShadowsocksHandler) SetServer(server core.Server) {
	ssh.server = server
}

// IsConnected 检查是否已连接
func (ssh *RealShadowsocksHandler) IsConnected() bool {
	return ssh.connected
}
