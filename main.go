package main

import (
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
	State      string
}

type Tunnels []Tunnel

type Jumphost struct {
	Id      int
	Name    string
	Tunnels Tunnels
}

type Jumphosts []Jumphost

func main() {
	juphosts := Jumphosts{
		{
			Id:   1,
			Name: "PreProduction",
			Tunnels: Tunnels{
				{
					Id:         1,
					Jumphost:   "PreProduction",
					Name:       "Grafana Cl1",
					Local_port: 123,
					Remote:     "1.2.3.4:443",
					URL:        "https://abcd.jhartman.pl",
					State:      "ok",
				},
				{
					Id:         2,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					State:      "nok",
				},
				{
					Id:         4,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					State:      "ok",
				},
				{
					Id:         5,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					State:      "nok",
				},
				{
					Id:         6,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					State:      "ok",
				},
				{
					Id:         7,
					Jumphost:   "PreProduction",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "2.2.3.4:443",
					URL:        "https://bcde.jhartman.pl",
					State:      "ok",
				},
			},
		},
		{
			Id:   1,
			Name: "Production",
			Tunnels: Tunnels{
				{
					Id:         1,
					Jumphost:   "Production",
					Name:       "Grafana Cl1",
					Local_port: 123,
					Remote:     "11.2.3.4:443",
					URL:        "https://mbcd.jhartman.pl",
					State:      "ok",
				},
				{
					Id:         2,
					Jumphost:   "Production",
					Name:       "Prometheus Cl1",
					Local_port: 123,
					Remote:     "12.2.3.4:443",
					URL:        "https://nbcd.jhartman.pl",
					State:      "nok",
				},
			},
		},
	}

	log.Println("Startup")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /")
		tmpl := template.Must(template.ParseFiles("templates/index.html",
			"templates/tunnels.html",
			"templates/server-group.html"))
		tmpl.Execute(w, juphosts)
	})
	http.ListenAndServe(":3000", r)

	// sweaters := Inventory{"wool", 17}
	// tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")

	// if err != nil {
	// 	panic(err)
	// }
	// err = tmpl.Execute(os.Stdout, sweaters)
	// if err != nil {
	// 	panic(err)
	// }
}
