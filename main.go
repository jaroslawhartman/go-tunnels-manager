package main

import (
	"flag"
	"html/template"
	"log"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Inventory struct {
	Material string
	Count    uint
}

type Tunnel struct {
	Id         int
	Jumphost   string
	Name       string
	Local_port int
	Remote     string
	URL        string
	Status     string
}

type Tunnels []Tunnel

type Jumphost struct {
	Id      int
	Name    string
	Command string
	Tunnels Tunnels
}

type Jumphosts []Jumphost

type TemplateArg struct {
	Page      string
	Jumphosts Jumphosts
}

func main() {
	var addr = flag.String("l", ":3000", "Listening [<host>]:<port>")
	flag.Parse()

	jumphosts := Jumphosts{
		{
			Id:      1,
			Name:    "PreProduction",
			Command: "ssh -NL {Local_port}:{Remote} docker-mtx",
			Tunnels: Tunnels{
				{
					Id:         1,
					Jumphost:   "PreProduction",
					Name:       "Grafana Cl1",
					Local_port: 123,
					Remote:     "1.2.3.4:443",
					URL:        "https://abcd.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         2,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "nok",
				},
				{
					Id:         4,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         5,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "nok",
				},
				{
					Id:         6,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         7,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
			},
		},
		{
			Id:      1,
			Name:    "Production",
			Command: "gcloud beta compute start-iap-tunnel postman-vm 22 --configuration tef-cloudlab2 --listen-on-stdin",
			Tunnels: Tunnels{
				{
					Id:         1,
					Jumphost:   "Production",
					Name:       "Grafana Cl1",
					Local_port: 123,
					Remote:     "11.2.3.4:443",
					URL:        "https://mbcd.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         2,
					Jumphost:   "Production",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "12.2.3.4:443",
					URL:        "https://nbcd.jhartman.pl",
					Status:     "nok",
				},
			},
		},
	}

	log.Printf("Starting server at %s", *addr)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /")
		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/server-group.html",
			"templates/jumphosts.html"))

		tmpl.Execute(w, TemplateArg{Page: "tunnels", Jumphosts: jumphosts})
	})

	r.Get("/jumphosts", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /jumphosts")
		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/server-group.html",
			"templates/jumphosts.html"))

		tmpl.Execute(w, TemplateArg{Page: "jumphosts", Jumphosts: jumphosts})
	})

	http.ListenAndServe(*addr, r)
}
