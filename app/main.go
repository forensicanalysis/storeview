package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
	"path/filepath"
)

// Constants
const htmlAbout = `Welcome on <b>Astilectron</b> demo!<br>
This is using the bootstrap and the bundler.`

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	debug = flag.Bool("d", true, "enables the debug mode")
	w     *astilectron.Window
	app   *astilectron.Astilectron
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())
	l.SetFlags(log.LstdFlags | log.Lshortfile)

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(options(l)); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}

func options(l *log.Logger) bootstrap.Options {
	return bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            "elementary",
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  *debug,
		Logger: l,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							l.Printf("About modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(a *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			app = a
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "open.html",
			MessageHandler: open,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(170),
				Width:           astikit.IntPtr(480),
				TitleBarStyle:   astilectron.TitleBarStyleHiddenInset,
				Resizable:       astikit.BoolPtr(false),
				Minimizable:     astikit.BoolPtr(false),
			},
		}},
	}
}

func cl(e astilectron.Event) (deleteListener bool) {
	windowCount--
	if windowCount == 0 {
		w.Show()
	}
	return true
}

var windowCount = 0

func storeWindow(store string) error {
	url := filepath.Join(app.Paths().DataDirectory(), "resources", "app", "index.html")
	window, err := app.NewWindow(url, &astilectron.WindowOptions{
		Title:         astikit.StrPtr(filepath.Base(store)),
		Center:        astikit.BoolPtr(true),
		Height:        astikit.IntPtr(800),
		Width:         astikit.IntPtr(1024),
		TitleBarStyle: astilectron.TitleBarStyleHiddenInset,
	})
	if err != nil {
		return err
	}
	window.OnMessage(handleMessages(window, handleStoreMessages(store), nil))
	window.On(astilectron.EventNameWindowEventClosed, cl)
	w.Hide()
	windowCount++

	return window.Create()
}

func handleMessages(w *astilectron.Window, messageHandler bootstrap.MessageHandler, l astikit.SeverityLogger) astilectron.ListenerMessage {
	return func(m *astilectron.EventMessage) (v interface{}) {
		// Unmarshal message
		var i bootstrap.MessageIn
		var err error
		if err = m.Unmarshal(&i); err != nil {
			l.Error(fmt.Errorf("unmarshaling message %+v failed: %w", *m, err))
			return
		}

		// Handle message
		var p interface{}
		if p, err = messageHandler(w, i); err != nil {
			l.Error(fmt.Errorf("handling message %+v failed: %w", i, err))
		}

		// Return message
		if p != nil {
			o := &bootstrap.MessageOut{Name: i.Name + ".callback", Payload: p}
			if err != nil {
				o.Name = "error"
			}
			v = o
		}
		return
	}
}
