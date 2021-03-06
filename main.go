package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
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
	glog.Infof("Starting photos on port %s.", *port)

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

	http.HandleFunc("/", handleMain)

	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	dirs, err := listDirs(*photoRoot)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// TODO: 404 when the path doesn't exist.

	p := r.URL.Path
	title := strings.Replace(p, "/", "", -1)
	if title == "" {
		title = *appTitle
	}

	images, err := listImages(*photoRoot, p)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	credits, err := ioutil.ReadFile(path.Join(*photoRoot, p, "credits.html"))
	if err != nil {
		credits = nil
	}

	imageJson, err := json.Marshal(images)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tmpl := template.New("layout")
	tmpl, err = tmpl.Parse(layoutTemplate)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	_ = template.New("photoswipe")
	_, err = tmpl.Parse(photoswipeTemplate)
	if err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := struct {
		AppName   string
		Title     string
		Credits   template.HTML
		Dirs      []string
		Images    []*Image
		ImageJson template.JS
		Home      bool
	}{
		*appTitle,
		title,
		template.HTML(string(credits)),
		dirs,
		images,
		template.JS(imageJson),
		p == "/",
	}

	if err := tmpl.ExecuteTemplate(w, "layout", &data); err != nil {
		glog.Error(err)
		http.Error(w, http.StatusText(500), 500)
	}
}
