package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Item struct {
	Id      int
	Text    string
	Checked bool
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	tmpl, err := template.ParseGlob("./views/*.html")

	if err != nil {
		log.Fatalf("unable to parse templates %e\n", err)
	}

	items := make([]Item, 0)

	e := echo.New()

	e.Renderer = &TemplateRenderer{
		templates: tmpl,
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", items)
	})

	e.POST("/todos", func(c echo.Context) error {
		text := c.FormValue("todo")
		items = append(items, Item{Id: len(items), Text: text, Checked: false})
		return c.Render(http.StatusOK, "items", items)
	})

	e.POST("/todos/:id/toggle", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "failed to parse id as integer")
		}

		if id > len(items) || id < 0 {
			return c.String(http.StatusNotFound, "item not found")
		}

		item := &items[id]
		item.Checked = !item.Checked

		return c.Render(http.StatusOK, "item", item)
	})

	e.Logger.Fatal(e.Start("localhost:3000"))
}
