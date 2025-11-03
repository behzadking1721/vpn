//go:build wireguard_stub_only
// +build wireguard_stub_only

package managers

import (
    "fmt"
    "vpnclient/src/core"
)

// connectWireGuard is a stub when WireGuard dependencies are not available.
func (cm *ConnectionManager) connectWireGuard(server *core.Server) error {
    return fmt.Errorf("wireguard not available in this build: ensure WireGuard driver/interface is installed and wgctrl dependencies are integrated")
}


