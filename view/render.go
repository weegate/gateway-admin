//@author wuyong
//@date   2018/1/11
//@desc

package view

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin/render"
)

type TplName string
type TplFiles []string
type RenderTpls map[TplName]TplFiles

type ITplRender interface {
	RegisterRenders(render *Render)
}

type BaseView struct {
	TplDir       string   `eg:"/fe/template/"`
	AppName      string   `eg:"abtest"`
	ModuleName   string   `eg:"policy"`
	OptPageNames []string `eg:"{add,list,update,index}"`
}

type Render map[string]*template.Template

//var _ render.HTMLRender = Render{}

func (v *BaseView) RegisterRenders(render *Render) {
	tplDir := v.TplDir
	layoutFile := strings.Join([]string{v.AppName, "tpl.layout.html"}, "/")

	//register index render
	render.AddFromFiles("index", tplDir+"/tpl.index.html")

	//register app index render
	render.AddFromFiles(v.AppName+"_index", tplDir+"/"+layoutFile, tplDir+"/"+v.AppName+"/tpl."+v.AppName+"_index.html")

	//register app module renders
	for _, optPageName := range v.OptPageNames {
		tplName := strings.Join([]string{v.AppName, v.ModuleName, optPageName}, "_")
		renderFile := strings.Join([]string{v.AppName, v.ModuleName, "tpl." + tplName + ".html"}, "/")
		render.AddFromFiles(tplName, tplDir+"/"+layoutFile, tplDir+"/"+renderFile)
	}
}

func New() Render {
	return make(Render)
}

func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r[name] = tmpl
}

// name must unique
func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.ParseGlob(glob))
	r.Add(name, tmpl)
	return tmpl
}

func (r *Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New("").Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r[name],
		Data:     data,
	}
}
