package vueloader

import (
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type VueLoader struct {
	layoutTemplate     *template.Template
	componentPath      string
	compiledComponents template.HTML
}

type VueComponent struct {
	FileName string
	Content  template.HTML
}

type VuePage struct {
	Title       string
	RootElement template.HTML
	Components  template.HTML
}

const defaultLayout = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>{{.Title}}</title>
        <script src="https://cdn.jsdelivr.net/npm/vue@2.6.14/dist/vue.min.js"></script>
        {{range .Components}}            
            {{.Content}}            
        {{end}}
    </head>
    <body>
        <main id="main-vue" v-cloak>
            {{.Root}}
        </main>
        <script>
            new Vue({el: "#main-vue"});
        </script>
    </body>
</html>
`

func (loader *VueLoader) loadComponents() []VueComponent {
	var result []VueComponent
	filepath.WalkDir(loader.componentPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(info.Name())) == ".vue" {
			c := VueComponent{
				FileName: info.Name(),
			}
			b, err := ioutil.ReadFile(path)
			if err == nil {
				c.Content = template.HTML("<!--####### BEGIN " + c.FileName + " #######-->\n" + string(b) + "<!--####### END " + c.FileName + " #######-->\n")
				result = append(result, c)
			}
		}
		return nil
	})
	return result
}

//LoadVuePage will write the html template with the Vue instance and all the detected SFC components into w io.Writer.
func (loader *VueLoader) LoadVuePage(w io.Writer, pageTitle string, rootComponent string) {
	loader.loadComponents()
	loader.layoutTemplate.Execute(w, VuePage{
		Title:       pageTitle,
		RootElement: template.HTML(rootComponent),
		Components:  loader.compiledComponents,
	})
}

/*
New creates a VueLoader instance which will parse your layout file.
If the layout does not exists, a default one will be used which will fetch Vue.js from jsdelivr CDN.
*/
func New(layoutFilename string, componentPath string) (*VueLoader, error) {
	index, err := template.ParseFiles(layoutFilename)
	if err != nil {
		log.Printf("Error loading layout '%s', using built-in template.", layoutFilename)
		index, err = template.New("index.html").Parse(defaultLayout)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &VueLoader{
		layoutTemplate: index,
		componentPath:  componentPath,
	}, nil
}
