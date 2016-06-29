package main

import (
	"flag"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/zefer/bass-photo/handlers"
	"gopkg.in/airbrake/glog.v1"
	"gopkg.in/airbrake/gobrake.v1"
)

var (
	port      = flag.String("port", ":8080", "listen port")
	photoRoot = flag.String("root", "", "photo root dir")

	abProjectID = flag.Int64("abprojectid", 0, "Airbrake project ID")
	abApiKey    = flag.String("abapikey", "", "Airbrake API key")
	abEnv       = flag.String("abenv", "development", "Airbrake environment name")
)

func main() {
	flag.Parse()
	glog.Infof("Starting API for Photos on port %s.", *port)

	if *abProjectID > int64(0) && *abApiKey != "" {
		airbrake := gobrake.NewNotifier(*abProjectID, *abApiKey)
		airbrake.SetContext("environment", *abEnv)
		glog.Gobrake = airbrake
	}

	http.Handle("/list", handlers.ListHandler(photoRoot))

	// http.HandleFunc("/websocket", websocket.Serve)
	// http.Handle("/next", handlers.NextHandler(client))
	// http.Handle("/previous", handlers.PreviousHandler(client))
	// http.Handle("/play", handlers.PlayHandler(client))
	// http.Handle("/pause", handlers.PauseHandler(client))
	// http.Handle("/randomOn", handlers.RandomOnHandler(client))
	// http.Handle("/randomOff", handlers.RandomOffHandler(client))
	// http.Handle("/files", handlers.FileListHandler(client))
	// http.Handle("/playlist", handlers.PlayListHandler(client))
	// http.Handle("/library/updated", handlers.LibraryUpdateHandler(client))

	// The front-end assets are served from a go-bindata file.
	http.Handle("/", http.FileServer(&assetfs.AssetFS{
		Asset: Asset, AssetDir: AssetDir, Prefix: "",
	}))

	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}
