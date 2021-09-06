package main

import (
	"log"
	"net/http"

	"github.com/hugorosario/govueloader/vueloader"
)

func main() {
	//create the VueLoader instance
	loader, err := vueloader.New("layout.html", "./components")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//call LoadVuePage to write the compiled template into the ResponseWriter.
		loader.LoadVuePage(w, "GoVueLoader Demo", "<hello-world></hello-world>")
	})
	http.ListenAndServe(":8082", nil)
}
