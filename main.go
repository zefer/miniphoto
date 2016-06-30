package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"gopkg.in/airbrake/glog.v1"
	"gopkg.in/airbrake/gobrake.v1"
)

var (
	port      = flag.String("port", ":8080", "listen port")
	photoRoot = flag.String("root", "", "photo root dir")
	appTitle  = flag.String("title", "Photos", "app title/name")

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
	pp := path.Join("templates", "photoswipe.html")

	dirs, err := listDirs(*photoRoot)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	path := r.URL.Path
	title := strings.Replace(path, "/", "", -1)
	if title == "" {
		title = *appTitle
	} else {
		title = *appTitle + ": " + title
	}

	imgJson, err := listImages(*photoRoot, path)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmpl, err := template.ParseFiles(lp, pp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := struct {
		Title  string
		Dirs   []string
		Images template.JS
		Home   bool
	}{
		title,
		dirs,
		template.JS(imgJson),
		path == "/",
	}

	if err := tmpl.ExecuteTemplate(w, "layout", &data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
