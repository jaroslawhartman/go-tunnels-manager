package main

import (
	"flag"
	"html/template"
	"log"
	"strconv"

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
	Page          string
	JumphostId    int
	TunnelId      int
	CurrentTunnel *Tunnel
	Jumphosts     Jumphosts
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
					Local_port: 124,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "nok",
				},
				{
					Id:         4,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 125,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         5,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 126,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "nok",
				},
				{
					Id:         6,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 127,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         7,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 128,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					Status:     "ok",
				},
			},
		},
		{
			Id:      2,
			Name:    "Production",
			Command: "gcloud beta compute start-iap-tunnel postman-vm 22 --configuration tef-cloudlab2 --listen-on-stdin",
			Tunnels: Tunnels{
				{
					Id:         1,
					Jumphost:   "Production",
					Name:       "Grafana Cl1",
					Local_port: 129,
					Remote:     "11.2.3.4:443",
					URL:        "https://mbcd.jhartman.pl",
					Status:     "ok",
				},
				{
					Id:         2,
					Jumphost:   "Production",
					Name:       "Prometheus Cl1",
					Local_port: 130,
					Remote:     "12.2.3.4:443",
					URL:        "https://nbcd.jhartman.pl",
					Status:     "nok",
				},
			},
		},
	}

	var current *Tunnel = nil

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

		current = &Tunnel{
			Id:         7,
			Jumphost:   "",
			Name:       "",
			Local_port: 0,
			Remote:     "",
			URL:        "",
			Status:     "nok",
		}

		tmpl.Execute(w, TemplateArg{Page: "tunnels", CurrentTunnel: current, Jumphosts: jumphosts})

		current = nil
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

	r.Get("/tunnel/get/{JumphostId}/{TunnelId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/get/")
		jummphostId, _ := strconv.Atoi(chi.URLParam(r, "JumphostId"))
		tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))

		log.Println(jummphostId)
		log.Println(tunnelId)

	OuterLoop:
		for _, j := range jumphosts {
			for _, t := range j.Tunnels {
				if j.Id == jummphostId && t.Id == tunnelId {
					current = &t
					log.Println("Found tunnel")
					log.Println(current)
					break OuterLoop
				}
			}
		}

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/server-group.html",
			"templates/jumphosts.html"))

		w.Header().Set("HX-Trigger-After-Settle", "showModal")

		tmpl.ExecuteTemplate(w, "tunnel-edit-modal", TemplateArg{
			Page:          "jumphosts",
			JumphostId:    jummphostId,
			TunnelId:      tunnelId,
			CurrentTunnel: current,
			Jumphosts:     jumphosts})
	})

	r.Post("/tunnel/update/{TunnelId}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /tunnel/update/")
		tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))
		log.Println(tunnelId)

		inputJumphost := r.FormValue("inputJumphost")
		inputName := r.FormValue("inputName")
		inputPort := r.FormValue("inputPort")
		inputRemote := r.FormValue("inputRemote")
		inputURL := r.FormValue("inputURL")

		log.Println("inputJumphost:" + inputJumphost)
		log.Println("inputName:" + inputName)
		log.Println("inputPort:" + inputPort)
		log.Println("inputRemote:" + inputRemote)
		log.Println("inputURL:" + inputURL)

		jumphosts[0].Tunnels[0].Jumphost = inputJumphost
		jumphosts[0].Tunnels[0].Name = inputName
		jumphosts[0].Tunnels[0].Local_port, _ = strconv.Atoi(inputPort)
		jumphosts[0].Tunnels[0].Remote = inputRemote
		jumphosts[0].Tunnels[0].URL = inputURL

		tmpl := template.Must(template.ParseFiles(
			"templates/index.html",
			"templates/tunnels.html",
			"templates/server-group.html",
			"templates/jumphosts.html"))

		current = &Tunnel{
			Id:         7,
			Jumphost:   "",
			Name:       "",
			Local_port: 0,
			Remote:     "",
			URL:        "",
			Status:     "nok",
		}

		tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{Page: "tunnels", CurrentTunnel: current, Jumphosts: jumphosts})

		current = nil
	})

	http.ListenAndServe(*addr, r)
}
