package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"telemetry.local/drone/internal/api"
	"telemetry.local/drone/internal/service"
	"telemetry.local/drone/internal/store"
)

func main() {
	options := api.ParseCLI(os.Args[1:])
	if err := api.ValidateCLI(options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	path := os.Getenv("TELEMETRY_DB")
	if strings.TrimSpace(path) == "" {
		path = "telemetry.db"
	}
	storage, err := store.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer storage.Close()
	application := service.New(storage, "cli-operator")
	switch strings.ToLower(options.Command) {
	case "serve":
		server := api.NewServer(application)
		if err := http.ListenAndServe("127.0.0.1:8080", server.Handler); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "report":
		summary, reportErr := application.BuildReport()
		if reportErr != nil {
			fmt.Fprintln(os.Stderr, reportErr)
			return
		}
		fmt.Println(summary.Total)
	case "search":
		records, searchErr := application.Search(options.Query, "")
		if searchErr != nil {
			fmt.Fprintln(os.Stderr, searchErr)
			return
		}
		for _, record := range records {
			fmt.Printf("%s\t%s\t%s\n", record.ID, record.Status, record.Title)
		}
	case "import":
		data, readErr := os.ReadFile(options.File)
		if readErr != nil {
			fmt.Fprintln(os.Stderr, readErr)
			return
		}
		text, importErr := application.ImportText(string(data), options.Batch)
		if importErr != nil {
			fmt.Fprintln(os.Stderr, importErr)
			return
		}
		fmt.Println(text)
	}
}
