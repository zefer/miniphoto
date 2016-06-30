package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path"
)

type Image struct {
	Path string `json:"path"`
	Src  string `json:"src"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}

func listImages(root *string, uri string) ([]byte, error) {
	dir := path.Join(*root, uri)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	images := readImageProperties(files, *root)

	b, err := json.Marshal(images)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func readImageProperties(orig []os.FileInfo, root string) []*Image {
	images := make([]*Image, 0)
	for _, file := range orig {
		w, h, err := getImageDimension(path.Join(root, file.Name()))
		if err == nil {
			images = append(images, &Image{
				Path: file.Name(),
				Src:  path.Join("/photo", file.Name()),
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
