package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"gopkg.in/airbrake/glog.v1"
)

type Image struct {
	Path string `json:"path"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}

func readImageProperties(orig []os.FileInfo, root string) []*Image {
	images := make([]*Image, 0)
	for _, file := range orig {
		w, h, err := getImageDimension(path.Join(root, file.Name()))
		if err == nil {
			images = append(images, &Image{
				Path: file.Name(),
				W:    w,
				H:    h,
			})
		}
	}
	return images
}

func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		return 0, 0, err
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Printf("%s %v\n", imagePath, err)
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

func ListHandler(root *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		uri := r.FormValue("uri")
		if uri == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dir := path.Join(*root, uri)

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		images := readImageProperties(files, *root)

		b, err := json.Marshal(images)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	})
}
