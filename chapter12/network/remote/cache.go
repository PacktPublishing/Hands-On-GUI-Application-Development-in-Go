package remote

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func cacheFileName(u string) string {
	id, _ := url.Parse(u)
	file := filepath.Base(id.Path)
	return path.Join("/tmp/", fmt.Sprintf("%s:%s", id.Hostname(), file))
}

func CacheStream(url string) io.ReadCloser {
	cacheFile := cacheFileName(url)
	if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		fmt.Println("Found cached file at", cacheFile)
		file, _ := os.Open(cacheFile)
		return file
	}

	fmt.Println("No cache found, downloading")
	stream := ReadStream(url)
	writer, _ := os.Create(cacheFile)
	io.Copy(writer, stream)
	stream.Close()
	writer.Close()

	fmt.Println("Saved to", cacheFile)
	stream, _ = os.Open(cacheFile)
	return stream
}
