package files

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const baseDir = "/srv/reports"

func ServeReport(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	f, err := os.Open(filepath.Join(baseDir, name))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func FetchAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user")
	resp, err := http.Get("https://avatars.example.com/u/" + url.PathEscape(id) + ".png")
	if err != nil {
		http.Error(w, "upstream", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
