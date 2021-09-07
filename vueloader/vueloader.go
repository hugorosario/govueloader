package vueloader

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
		{{.Styles}}         
    </head>
    <body>
		{{.Templates}}         
		
		<main id="main" v-cloak>
            {{.RootElement}}
        </main>

		{{.Scripts}}		
        <script>
            new Vue({el: "#main"});
        </script>   
    </body>
</html>
`

//VueLoaderConfig contains all the configuration properties for a VueLoader instance.
type VueLoaderConfig struct {
	/*
		Layout is the base template for a Vue page.
		Can be provided as HTML or as a file path whose content will be loaded automatically.
	*/
	Layout string
	/*
		ComponentPath is the path where the compiler will search for Vue Single File Components.
	*/
	ComponentPath string
	/*
		CompileEveryRequest can be set to compile every time the LoadVuePage method is executed.
		Makes development easier without the need to restart the server every time you change something.
	*/
	CompileEveryRequest bool
}

type VueLoader struct {
	Config          *VueLoaderConfig
	layoutTemplate  *template.Template
	compiledHTML    string
	compiledScripts string
	compiledStyles  string
}

type vueComponent struct {
	FileName      string
	ID            string
	HtmlContent   string
	ScriptContent string
	StyleContent  string
}

type vuePage struct {
	Title       string
	RootElement template.HTML
	Scripts     template.HTML
	Styles      template.HTML
	Templates   template.HTML
}

//load the single file components from the filesystem
func (loader *VueLoader) load() []vueComponent {
	var components []vueComponent
	filepath.WalkDir(loader.Config.ComponentPath, func(path string, info fs.DirEntry, err error) error {
		defer recover()
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(info.Name())) == ".vue" {
			c := vueComponent{
				FileName:      info.Name(),
				ID:            "",
				HtmlContent:   "",
				ScriptContent: "",
				StyleContent:  "",
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
			if err != nil {
				return err
			}
			tmpl := content.Find("template").First()
			c.ID = tmpl.AttrOr("id", info.Name())
			c.HtmlContent, _ = goquery.OuterHtml(tmpl)
			c.ScriptContent = content.Find("script").First().Text()
			//there might be multiple style tags
			c.StyleContent = ""
			content.Find("style").Each(func(i int, s *goquery.Selection) {
				style := s.Text()
				//TODO process scoped styles
				if style != "" {
					c.StyleContent += style + "\n"
				}
			})
			components = append(components, c)
		}
		return nil
	})
	return components
}

func (loader *VueLoader) compile() error {
	//check if config.Layout is a file path and replace it with the file contents
	stat, err := os.Stat(loader.Config.Layout)
	if stat != nil && err == nil {
		content, err := os.ReadFile(loader.Config.Layout)
		if err != nil {
			return err
		}
		loader.Config.Layout = string(content)
	}
	//check if the layout is valid HTML
	if !strings.Contains(loader.Config.Layout, "<html>") {
		return errors.New("template is not valid HTML")
	}
	//parse the template at config.Layout
	loader.layoutTemplate, err = template.New("layout").Parse(loader.Config.Layout)
	if err != nil {
		return err
	}
	//load components and compile all the separate tags
	components := loader.load()
	loader.compiledHTML = "\n<!--####### Templates #######-->\n"
	loader.compiledScripts = "\n<!--####### Scripts #######-->\n"
	loader.compiledStyles = "\n<!--####### Styles #######-->\n"
	loader.compiledScripts += "<script>\n"
	loader.compiledStyles += "<style>\n"
	for _, v := range components {
		loader.compiledHTML += v.HtmlContent + "\n"
		loader.compiledScripts += v.ScriptContent + "\n"
		loader.compiledStyles += v.StyleContent + "\n"
	}
	loader.compiledScripts += "</script>\n"
	loader.compiledStyles += "</style>\n"
	return nil
}

//LoadVuePage will write the final compiled page with the Vue instance and all the detected components into w io.Writer.
func (loader *VueLoader) LoadVuePage(w io.Writer, pageTitle string, rootComponent string) {
	if loader.Config.CompileEveryRequest {
		if err := loader.compile(); err != nil {
			log.Println(err)
			return
		}
	}
	err := loader.layoutTemplate.Execute(w, vuePage{
		Title:       pageTitle,
		RootElement: template.HTML(rootComponent),
		Scripts:     template.HTML(loader.compiledScripts),
		Styles:      template.HTML(loader.compiledStyles),
		Templates:   template.HTML(loader.compiledHTML),
	})
	if err != nil {
		log.Println(err)
	}
}

/*
NewWithConfig creates a VueLoader instance with the provided VueLoaderConfig.
If the layout cannot be loaded, a default one will be used which will fetch Vue.js from jsdelivr CDN.
*/
func NewWithConfig(config *VueLoaderConfig) (*VueLoader, error) {
	loader := &VueLoader{
		Config: config,
	}
	if err := loader.compile(); err != nil {
		return nil, err
	}
	return loader, nil
}

/*
New creates a VueLoader instance with the default VueLoaderConfig.
*/
func New() (*VueLoader, error) {
	return NewWithConfig(NewConfig())
}

//NewConfig returns a new VueLoaderConfig instance with all the default properties assigned
func NewConfig() *VueLoaderConfig {
	return &VueLoaderConfig{
		Layout:              defaultLayout,
		ComponentPath:       "./views",
		CompileEveryRequest: false,
	}
}
