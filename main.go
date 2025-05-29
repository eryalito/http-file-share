package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/eryalito/http-file-share/internal/listener"
	"github.com/eryalito/http-file-share/internal/services"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	httpFileServer, err := listener.NewHttpFileServer()
	if err != nil {
		log.Fatalf("Failed to start HTTP file server: %v", err)
	}
	log.Printf("HTTP file server is running on port %d", httpFileServer.Port())
	log.Printf("%s", httpFileServer.Addresses())

	app := application.New(application.Options{
		Name:        "http-file-share",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&services.DialogsService{}),
			application.NewService(&services.HttpFileServerService{
				HttpFileServer: httpFileServer,
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "HTTP File Share",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
		Width:            480,
		Height:           600,
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
