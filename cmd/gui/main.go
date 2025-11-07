package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"

	// Fix import paths
	"vpnclient/internal/database"
	"vpnclient/internal/managers"
	"vpnclient/src/core"
)

func main() {
	a := app.NewWithID("vpnclient.gui")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("VPN Client")
	w.Resize(fyne.NewSize(800, 600))

	// Init data dir and stores
	dataDir := filepath.Join(".", "data")
	if env := os.Getenv("VPN_DATA_DIR"); env != "" {
		dataDir = env
	}
	store, err := database.NewDB(dataDir)
	if err != nil {
		dialog := widget.NewLabel(fmt.Sprintf("Failed to init DB: %v", err))
		w.SetContent(container.NewCenter(dialog))
		w.ShowAndRun()
		return
	}
	defer store.Close()

	serverMgr := managers.NewServerManager(store)
	connMgr := managers.NewConnectionManager()

	// Bindings
	statusStr := binding.NewString()
	_ = statusStr.Set(connMgr.GetStatusString())
	uploadStr := binding.NewString()
	downloadStr := binding.NewString()

	// Server list
	servers, _ := serverMgr.GetAllServers()
	serverNames := []string{}
	for _, s := range servers {
		serverNames = append(serverNames, fmt.Sprintf("%s [%s]", s.Name, s.Protocol))
	}
	listData := binding.NewStringList()
	_ = listData.Set(serverNames)

	selectedIndex := -1
	var getSelectedServer = func() *core.Server {
		if selectedIndex < 0 || selectedIndex >= len(servers) {
			return nil
		}
		s, _ := serverMgr.GetServer(servers[selectedIndex].ID)
		return s
	}

	list := widget.NewListWithData(listData,
		func() fyne.CanvasObject { return widget.NewLabel("item") },
		func(di binding.DataItem, co fyne.CanvasObject) {
			s, _ := di.(binding.String).Get()
			co.(*widget.Label).SetText(s)
		},
	)
	list.OnSelected = func(id widget.ListItemID) { selectedIndex = int(id) }

	// Controls
	connectBtn := widget.NewButtonWithIcon("Connect", theme.ConfirmIcon(), func() {
		var server *core.Server
		if selectedIndex >= 0 {
			server = getSelectedServer()
		} else {
			// If none selected, try fastest enabled
			server, _ = serverMgr.GetFastestServer()
		}

		if server != nil {
			_ = doConnect(connMgr, server, statusStr)
		}
	})

	disconnectBtn := widget.NewButtonWithIcon("Disconnect", theme.CancelIcon(), func() {
		if err := connMgr.Disconnect(); err == nil {
			_ = statusStr.Set(connMgr.GetStatusString())
		}
	})

	// Subscription section
	subscriptionEntry := widget.NewEntry()
	subscriptionEntry.SetPlaceHolder("Enter subscription URL (Clash, Surfboard, etc.)")

	importBtn := widget.NewButton("Import", func() {
		url := subscriptionEntry.Text
		if url != "" {
			// Create a new subscription
			sub := &core.Subscription{
				Name:       "Imported Subscription",
				URL:        url,
				AutoUpdate: true,
			}

			// Save subscription
			if err := serverMgr.AddSubscription(sub); err != nil {
				// In a real app, show an error dialog
				fmt.Printf("Error adding subscription: %v\n", err)
				return
			}

			// Update server list
			servers, _ = serverMgr.GetAllServers()
			serverNames = []string{}
			for _, s := range servers {
				serverNames = append(serverNames, fmt.Sprintf("%s [%s]", s.Name, s.Protocol))
			}
			_ = listData.Set(serverNames)

			// Clear entry
			subscriptionEntry.SetText("")
		}
	})

	// Status indicators
	statusLabel := widget.NewLabelWithData(statusStr)
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	uploadLabel := widget.NewLabelWithData(uploadStr)
	downloadLabel := widget.NewLabelWithData(downloadStr)

	uploadIcon := widget.NewIcon(theme.UploadIcon())
	downloadIcon := widget.NewIcon(theme.DownloadIcon())

	// Layout
	controls := container.NewHBox(connectBtn, disconnectBtn)

	stats := container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, uploadIcon, nil, uploadLabel),
		container.NewBorder(nil, nil, downloadIcon, nil, downloadLabel),
	)

	statusContainer := container.NewVBox(
		widget.NewLabel("Connection Status:"),
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("Traffic Stats:"),
		stats,
	)

	subscriptionContainer := container.NewVBox(
		widget.NewLabel("Import Subscription:"),
		subscriptionEntry,
		importBtn,
		widget.NewSeparator(),
	)

	content := container.NewBorder(
		container.NewVBox(controls, subscriptionContainer),
		statusContainer,
		nil,
		nil,
		list)
	w.SetContent(content)

	// Poller for status+stats
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = statusStr.Set(connMgr.GetStatusString())
			sent, recv := connMgr.GetDataUsage()
			_ = uploadStr.Set(humanBytes(sent))
			_ = downloadStr.Set(humanBytes(recv))
		}
	}()

	w.ShowAndRun()
}

func doConnect(cm *managers.ConnectionManager, s *core.Server, statusStr binding.String) error {
	if err := cm.Connect(s); err != nil {
		_ = statusStr.Set(cm.GetStatusString())
		return err
	}
	_ = statusStr.Set(cm.GetStatusString())
	return nil
}

func humanBytes(v int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)
	switch {
	case v >= gb:
		return fmt.Sprintf("%.2f GB", float64(v)/float64(gb))
	case v >= mb:
		return fmt.Sprintf("%.2f MB", float64(v)/float64(mb))
	case v >= kb:
		return fmt.Sprintf("%.2f KB", float64(v)/float64(kb))
	default:
		return fmt.Sprintf("%d B", v)
	}
}
