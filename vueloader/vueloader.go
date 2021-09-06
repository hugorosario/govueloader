package vueloader

import (
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
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

func (loader *VueLoader) compileComponents() []VueComponent {
	loader.compiledComponents = ""
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
				loader.compiledComponents += c.Content
				result = append(result, c)
			}
		}
		return nil
	})
	return result
}

//LoadVuePage will write the final compiled page with the Vue instance and all the detected SFC components into w io.Writer.
func (loader *VueLoader) LoadVuePage(w io.Writer, pageTitle string, rootComponent string) {
	loader.compileComponents()
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
		return nil, err
	}
	return &VueLoader{
		layoutTemplate: index,
		componentPath:  componentPath,
	}, nil
}
