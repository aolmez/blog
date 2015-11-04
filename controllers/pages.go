package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

//PageShow handles /pages/:id route
func PageShow(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := ctx.Value("template").(*template.Template)
	data := helpers.DefaultData(ctx)
	if r.Method == "GET" {

		id := r.URL.Path[len("/pages/"):]
		page, err := models.GetPage(id)
		if err != nil || !page.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Page"] = page
		data["Title"] = page.Name
		data["Active"] = fmt.Sprintf("pages/%s", id)
		tmpl.Lookup("pages/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageIndex handles GET /admin/pages route
func PageIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := ctx.Value("template").(*template.Template)
	data := helpers.DefaultData(ctx)
	if r.Method == "GET" {

		list, err := models.GetPages()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Title"] = "List of pages"
		data["Active"] = "pages"
		data["List"] = list
		tmpl.Lookup("pages/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageCreate handles /admin/new_page route
func PageCreate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := ctx.Value("template").(*template.Template)
	session := ctx.Value("session").(*sessions.Session)
	data := helpers.DefaultData(ctx)
	if r.Method == "GET" {

		data["Title"] = "Create page"
		data["Active"] = "pages"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			Name:        r.PostFormValue("name"),
			Description: r.PostFormValue("description"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := page.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_page", 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageUpdate handles /admin/edit_page/:id route
func PageUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := ctx.Value("template").(*template.Template)
	session := ctx.Value("session").(*sessions.Session)
	data := helpers.DefaultData(ctx)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_page/"):]
		page, err := models.GetPage(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = "Edit page"
		data["Active"] = "pages"
		data["Page"] = page
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			ID:          helpers.Atoi64(r.PostFormValue("id")),
			Name:        r.PostFormValue("name"),
			Description: r.PostFormValue("description"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := page.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageDelete handles /admin/delete_page route
func PageDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl := ctx.Value("template").(*template.Template)

	if r.Method == "POST" {

		page, err := models.GetPage(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := page.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}