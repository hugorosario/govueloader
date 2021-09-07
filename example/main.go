package main

import (
	"log"
	"net/http"

	"github.com/hugorosario/govueloader/vueloader"
)

func main() {
	/*
		Create the VueLoader instance with the default config.
		It will scan the "views" folder for components and use an internal default base layout with Vue loaded from a CDN.
	*/
	defaultloader, err := vueloader.New()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//call LoadVuePage to write the compiled template into the ResponseWriter.
		defaultloader.LoadVuePage(w, "GoVueLoader Demo", "<hello-world></hello-world>")
	})

	/*
		Create a VueLoaderConfig instance with:
		- A custom layout from file "layout.html".
		- A different component path from folder "vuetify".
		- A flag to force recompilation on every request to make development easier.
	*/
	config := vueloader.NewConfig()
	//set your layout file or content
	config.Layout = "./layout.html"
	//set your component path
	config.ComponentPath = "./vuetify"
	//compile with every request, good for development
	config.CompileEveryRequest = true
	//create a VueLoader instance with our custom config
	loader, err := vueloader.NewWithConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/vuetify", func(w http.ResponseWriter, r *http.Request) {
		//call LoadVuePage to write the compiled template into the ResponseWriter.
		loader.LoadVuePage(w, "GoVueLoader Vuetify Demo", "<app></app>")
	})

	http.ListenAndServe(":8082", nil)
}
