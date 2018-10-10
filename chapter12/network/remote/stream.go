package remote

import (
	"image"
	"io"
	"log"
	"net/http"
)

func ReadStream(url string) io.ReadCloser {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		log.Fatal("Error reading URL", err)
	}

	return res.Body
}

func RemoteImage(stream io.ReadCloser) image.Image {
	defer stream.Close()

	m, _, err := image.Decode(stream)
	if err != nil {
		log.Fatal("Error reading image", err)
	}

	return m
}