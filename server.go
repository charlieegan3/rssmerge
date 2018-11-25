package main

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

func main() {
	box := packr.NewBox("./web")

	http.Handle("/build", http.HandlerFunc(RSSMergeHandler))
	http.Handle("/", http.FileServer(box))
	http.ListenAndServe(":3000", nil)
}
