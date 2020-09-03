package ginapp

import (
	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
	"net/http"
	"pokemon/pkg/config"
)

var TemplatePath = "./templates"
var Debug = true

var View *jet.Set
var htmlContentType = []string{"text/html; charset=utf-8"}

type HTML struct {
	//Engine *jet.VarMap
	TemplateName string
	Data         interface{}
}

func init() {
	Debug := !config.Config().EnvProd
	View = jet.NewHTMLSet(TemplatePath).SetDevelopmentMode(Debug)
}

func CleanViewCache(c *gin.Context) {
	View = jet.NewHTMLSet(TemplatePath).SetDevelopmentMode(Debug)
	c.String(http.StatusOK, "clean template cache ok!")
}

func NewHTML(TemplateName string, data interface{}) *HTML {
	return &HTML{
		TemplateName: TemplateName,
		Data:         data,
	}
}

func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	//return r.Template.ExecuteTemplate(w, r.Name, r.Data)
	t, err := View.GetTemplate(r.TemplateName)
	if t == nil || err != nil {
		// template could not be loaded
		panic("template file doesn't exist!:err:" + err.Error())
	}

	vars := make(jet.VarMap)
	for i, v := range r.Data.(gin.H) {
		vars.Set(i, v)
	}
	return t.Execute(w, vars, nil)
}

func (r HTML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}
}
