package cobraserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zserge/webview"
)

func UICommand(name string, width, height int, staticPath pkger.Dir, commands ...*Command) *cobra.Command {
	commandsMap := map[string]*Command{}
	for _, c := range commands {
		commandsMap[c.Name] = c
	}

	return &cobra.Command{Use: "ui", Short: "Start the ui", Run: func(_ *cobra.Command, _ []string) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer ln.Close()
			http.Handle("/", http.FileServer(staticPath))
			log.Fatal(http.Serve(ln, nil))
		}()
		w := webview.New(webview.Settings{
			Width:  width,
			Height: height,
			Title:  name,
			URL:    "http://" + ln.Addr().String(),
			ExternalInvokeCallback: func(w webview.WebView, dataS string) {
				// decode input
				data := strings.NewReader(dataS)
				var args map[string]interface{}
				d := json.NewDecoder(data)
				err := d.Decode(&args)
				if err != nil {
					rpcError(w, err)
					return
				}

				// get command
				cmdName := args["cmd"].(string)       // TODO: error handling
				callback := args["callback"].(string) // TODO: error handling
				delete(args, "cmd")
				delete(args, "callback")

				// parse arguments to the correct types
				var argsString []string
				for name, arg := range args {
					argsString = append(argsString, "--"+name, fmt.Sprint(arg))
				}
				flagset := &pflag.FlagSet{}
				setupFlags(commandsMap[cmdName], flagset, "json")
				flagset.Parse(argsString)

				buf := &bytes.Buffer{}
				err = commandsMap[cmdName].Handler(buf, data, flagset)

				if err != nil {
					rpcError(w, err)
					return
				}
				w.Eval(fmt.Sprintf("window.%s(%s)", callback, buf.String()))
			},
			Resizable: true,
		})
		defer w.Exit()
		w.Run()
	}}
}

func rpcError(w webview.WebView, err error) {
	log.Println("handleRPC error:", err)
	w.Eval(fmt.Sprintf("window.error({'error':'%s'})", err.Error()))
}
