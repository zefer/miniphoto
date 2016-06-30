package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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

	// http.Handle("/list", handlers.ListHandler(photoRoot))

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

	// http.Handle("/assets", http.FileServer(&assetfs.AssetFS{
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(&assetfs.AssetFS{
		Asset: Asset, AssetDir: AssetDir, Prefix: "",
	})))

	fs := http.FileServer(http.Dir(*photoRoot))
	http.Handle("/photo/", http.StripPrefix("/photo/", fs))

	http.HandleFunc("/", serveTemplate)

	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	np := path.Join("templates", "nav.html")
	pp := path.Join("templates", "photoswipe.html")
	// fp := path.Join("templates", r.URL.Path)

	imgJson, err := listImages(photoRoot, "")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmpl, err := template.ParseFiles(lp, np, pp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := struct {
		Title  string
		Dirs   []string
		Images template.JS
	}{
		"Banana",
		[]string{},
		template.JS(imgJson),
	}

	if err := tmpl.ExecuteTemplate(w, "layout", &data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
