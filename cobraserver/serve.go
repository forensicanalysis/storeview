package cobraserver

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func ServeCommand(staticPath pkger.Dir, commands ...*Command) *cobra.Command {
	router := mux.NewRouter()

	serveCmd := &cobra.Command{Use: "serve", Short: "Start the api server"}

	re := regexp.MustCompile(`^(\w+)(?:\[(.+)\])?$`)

	for _, commandP := range commands {
		command := *commandP
		router.HandleFunc("/api"+command.Route, func(w http.ResponseWriter, r *http.Request) {
			setupResponse(&w, r)
			if (*r).Method == "OPTIONS" {
				return
			}

			var argsString []string
			for key, values := range r.URL.Query() {
				submatches := re.FindAllStringSubmatch(key, -1)
				group := ""
				if submatches[0][2] != "" {
					group = submatches[0][2] + ":"
				}
				for _, value := range values {
					argsString = append(argsString, "--"+submatches[0][1], group+value)
				}
			}
			argsString = append(argsString, serveCmd.Flags().Args()...)
			flagset := &pflag.FlagSet{}
			setupFlags(&command, flagset, "json")
			flagset.Parse(argsString)

			if err := command.Handler(w, r.Body, flagset); err != nil {
				fmt.Fprintf(w, "{'error':'%s'}", err)
				// w.WriteHeader(http.StatusInternalServerError)
			}
		}).Methods(command.Method, http.MethodOptions)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(staticPath)))

	var port int
	serveCmd.Run = func(_ *cobra.Command, _ []string) {
		http.Handle("/", router)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port")
	return serveCmd
}
