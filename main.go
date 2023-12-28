package main

import (
	"flag"
	"html/template"
	"log"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	d "jhartman.pl/go-tunnels-ui/db"
)

type TemplateArg struct {
	Page          string
	Error         string
	JumphostId    int
	TunnelId      int
	CurrentTunnel int
	Jumphosts     d.Jumphosts
	Tunnels       d.Tunnels
}

func main() {
	db := d.Open("tunnels.sqlite")
	defer d.Close(db)

	var addr = flag.String("l", ":3000", "Listening [<host>]:<port>")
	flag.Parse()

	log.Printf("Starting server at %s", *addr)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /")
		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		tunnels := d.GetTunnels(db, "")
		jumphosts := d.GetJumphosts(db)

		tmpl.Execute(w, TemplateArg{Page: "tunnels", CurrentTunnel: 0, Jumphosts: jumphosts, Tunnels: tunnels})
	})

	r.Get("/jumphosts", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /jumphosts")
		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		jumphosts := d.GetJumphosts(db)

		tmpl.Execute(w, TemplateArg{Page: "jumphosts", Jumphosts: jumphosts, Tunnels: nil})
	})

	r.Delete("/jumphost/delete/{Id}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /jumphost/delete/")
		id, _ := strconv.Atoi(chi.URLParam(r, "Id"))
		log.Println(id)

		err := d.DeleteJumphost(db, id)

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		jumphosts := d.GetJumphosts(db)

		tmpl.ExecuteTemplate(w, "jumphosts", TemplateArg{
			Page:      "jumphosts",
			Error:     err,
			Jumphosts: jumphosts,
			Tunnels:   nil})
	})

	r.Get("/tunnel/get/{TunnelId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/get/")
		tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))

		log.Println(tunnelId)

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		tunnels := d.GetTunnels(db, "")
		jumphosts := d.GetJumphosts(db)

		w.Header().Set("HX-Trigger-After-Settle", "showModal")

		tmpl.ExecuteTemplate(w, "tunnel-edit-modal", TemplateArg{
			CurrentTunnel: tunnelId,
			Tunnels:       tunnels,
			Jumphosts:     jumphosts})
	})

	r.Post("/tunnel/update/{TunnelId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/update/")
		tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))
		log.Println(tunnelId)

		inputJumphost, _ := strconv.Atoi(r.FormValue("inputJumphost"))
		inputName := r.FormValue("inputName")
		inputPort, _ := strconv.Atoi(r.FormValue("inputPort"))
		inputRemote := r.FormValue("inputRemote")
		inputURL := r.FormValue("inputURL")

		err := d.UpdateTunnel(db, tunnelId, &d.Tunnel{
			JumphostId: inputJumphost,
			Name:       inputName,
			Local_port: inputPort,
			Remote:     inputRemote,
			URL:        inputURL,
		})

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		tunnels := d.GetTunnels(db, "")
		jumphosts := d.GetJumphosts(db)

		tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
			Page:          "tunnels",
			Error:         err,
			CurrentTunnel: 0,
			Tunnels:       tunnels,
			Jumphosts:     jumphosts})
	})

	r.Post("/tunnel/add", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/add/")

		inputJumphost, _ := strconv.Atoi(r.FormValue("inputJumphost"))
		inputName := r.FormValue("inputName")
		inputPort, _ := strconv.Atoi(r.FormValue("inputPort"))
		inputRemote := r.FormValue("inputRemote")
		inputURL := r.FormValue("inputURL")

		err := d.AddTunnel(db, &d.Tunnel{
			JumphostId: inputJumphost,
			Name:       inputName,
			Local_port: inputPort,
			Remote:     inputRemote,
			URL:        inputURL,
		})

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		tunnels := d.GetTunnels(db, "")
		jumphosts := d.GetJumphosts(db)

		tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
			Page:          "tunnels",
			Error:         err,
			CurrentTunnel: 0,
			Tunnels:       tunnels,
			Jumphosts:     jumphosts})
	})

	r.Delete("/tunnel/delete/{TunnelId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/delete/")
		tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))
		log.Println(tunnelId)

		err := d.DeleteTunnel(db, tunnelId)

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/jumphosts.html"))

		tunnels := d.GetTunnels(db, "")
		jumphosts := d.GetJumphosts(db)

		tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
			Page:          "tunnels",
			Error:         err,
			CurrentTunnel: 0,
			Tunnels:       tunnels,
			Jumphosts:     jumphosts})
	})

	http.ListenAndServe(*addr, r)
	// go func() {
	// 	http.ListenAndServe(*addr, r)
	// }()

	// main loop

	// for {
	// 	log.Printf("Help")

	// 	time.Sleep(1 * time.Second)
	// }
}
