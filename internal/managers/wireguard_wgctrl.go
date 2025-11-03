//go:build wireguard
// +build wireguard

package managers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"vpnclient/src/core"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (cm *ConnectionManager) connectWireGuard(server *core.Server) error {
	if server == nil {
		return fmt.Errorf("nil server")
	}

	cfg := server.Config
	ifaceName, _ := cfg["interface_name"].(string)
	if ifaceName == "" {
		ifaceName = "wg0"
	}
	privKeyStr, _ := cfg["private_key"].(string)
	peerPubKeyStr, _ := cfg["peer_public_key"].(string)
	endpointStr, _ := cfg["endpoint"].(string)
	keepAlive := intFromAnyWG(cfg["persistent_keepalive"], 25)

	if privKeyStr == "" || peerPubKeyStr == "" || endpointStr == "" {
		return fmt.Errorf("missing wireguard config: private_key, peer_public_key, endpoint are required")
	}

	privKey, err := wgtypes.ParseKey(privKeyStr)
	if err != nil {
		return fmt.Errorf("invalid private_key: %w", err)
	}
	peerPubKey, err := wgtypes.ParseKey(peerPubKeyStr)
	if err != nil {
		return fmt.Errorf("invalid peer_public_key: %w", err)
	}

	host, portStr, err := net.SplitHostPort(endpointStr)
	if err != nil {
		return fmt.Errorf("invalid endpoint (want host:port): %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid endpoint port: %s", portStr)
	}
	udpAddr := &net.UDPAddr{IP: net.ParseIP(host), Port: port}
	if udpAddr.IP == nil {
		ips, rErr := net.LookupIP(host)
		if rErr != nil || len(ips) == 0 {
			return fmt.Errorf("failed to resolve endpoint host: %s", host)
		}
		udpAddr = &net.UDPAddr{IP: ips[0], Port: port}
	}

	var allowedIPs []wgtypes.IPNet
	if arr, ok := cfg["allowed_ips"].([]interface{}); ok {
		for _, v := range arr {
			if s, ok := v.(string); ok {
				if ipn, perr := parseCIDRWG(s); perr == nil {
					allowedIPs = append(allowedIPs, ipn)
				}
			}
		}
	}

	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("wgctrl open failed: %w", err)
	}
	defer client.Close()

	if _, err := client.Device(ifaceName); err != nil {
		return fmt.Errorf("wireguard interface %s not found: create it first (wg-quick or OS tool)", ifaceName)
	}

	peer := wgtypes.PeerConfig{
		PublicKey:                   peerPubKey,
		Endpoint:                    udpAddr,
		PersistentKeepaliveInterval: durationPtrSecondsWG(keepAlive),
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  allowedIPs,
	}

	devCfg := wgtypes.Config{
		PrivateKey:   &privKey,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{peer},
	}

	if err := client.ConfigureDevice(ifaceName, devCfg); err != nil {
		return fmt.Errorf("configure wireguard device failed: %w", err)
	}

	return nil
}

func durationPtrSecondsWG(s int) *time.Duration {
	if s <= 0 {
		return nil
	}
	d := time.Duration(s) * time.Second
	return &d
}

func intFromAnyWG(v interface{}, def int) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	case string:
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
	}
	return def
}

func parseCIDRWG(s string) (wgtypes.IPNet, error) {
	_, ipn, err := net.ParseCIDR(s)
	if err != nil {
		return wgtypes.IPNet{}, err
	}
	return wgtypes.IPNet{IP: ipn.IP, Mask: ipn.Mask}, nil
}
