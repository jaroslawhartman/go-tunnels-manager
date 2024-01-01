package tunnelsmgr

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"strconv"
	"sync"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Tunnelmgr struct {
	sync.Mutex
	db  *sql.DB
	cfg TunnelmgrCfg
}

type TunnelmgrCfg struct {
	DBFilename string
	Listener   string
}

type TemplateArg struct {
	Page          string
	Error         string
	JumphostId    int
	TunnelId      int
	CurrentTunnel int
	Jumphosts     Jumphosts
	Tunnels       Tunnels
}

///////////////////////////
//
// Chi Handlers
//
///////////////////////////

func (t *Tunnelmgr) HandlerGetRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /")
	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/tunnels.html",
		"templates/jumphosts.html"))

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	tmpl.Execute(w, TemplateArg{Page: "tunnels", CurrentTunnel: 0, Jumphosts: jumphosts, Tunnels: tunnels})
}

func (t *Tunnelmgr) HandlerGetJumphosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /jumphosts")
	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/tunnels.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.Execute(w, TemplateArg{Page: "jumphosts", Jumphosts: jumphosts, Tunnels: nil})
}

func (t *Tunnelmgr) HandlerDeleteJumphosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /jumphost/delete/")
	id, _ := strconv.Atoi(chi.URLParam(r, "Id"))
	log.Println(id)

	err := t.DeleteJumphost(id)

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/tunnels.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "jumphosts", TemplateArg{
		Page:      "jumphosts",
		Error:     err,
		Jumphosts: jumphosts,
		Tunnels:   nil})
}

func (t *Tunnelmgr) HandlerGetTunnel(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /tunnel/get/")
	tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))

	log.Println(tunnelId)

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/tunnels.html",
		"templates/jumphosts.html"))

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	w.Header().Set("HX-Trigger-After-Settle", "showModal")

	tmpl.ExecuteTemplate(w, "tunnel-edit-modal", TemplateArg{
		CurrentTunnel: tunnelId,
		Tunnels:       tunnels,
		Jumphosts:     jumphosts})
}

func (t *Tunnelmgr) HandlerPostUpdateTunnel(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /tunnel/update/")
	tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))
	log.Println(tunnelId)

	inputJumphost, _ := strconv.Atoi(r.FormValue("inputJumphost"))
	inputName := r.FormValue("inputName")
	inputPort, _ := strconv.Atoi(r.FormValue("inputPort"))
	inputRemote := r.FormValue("inputRemote")
	inputURL := r.FormValue("inputURL")

	err := t.UpdateTunnel(tunnelId, &Tunnel{
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

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
		Page:          "tunnels",
		Error:         err,
		CurrentTunnel: 0,
		Tunnels:       tunnels,
		Jumphosts:     jumphosts})
}

func (t *Tunnelmgr) HandlerPostAddTunnel(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /tunnel/add/")

	inputJumphost, _ := strconv.Atoi(r.FormValue("inputJumphost"))
	inputName := r.FormValue("inputName")
	inputPort, _ := strconv.Atoi(r.FormValue("inputPort"))
	inputRemote := r.FormValue("inputRemote")
	inputURL := r.FormValue("inputURL")

	err := t.AddTunnel(&Tunnel{
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

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
		Page:          "tunnels",
		Error:         err,
		CurrentTunnel: 0,
		Tunnels:       tunnels,
		Jumphosts:     jumphosts})
}

func (t *Tunnelmgr) HandlerDeleteTunnel(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /tunnel/delete/")
	tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))
	log.Println(tunnelId)

	err := t.DeleteTunnel(tunnelId)

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/tunnels.html",
		"templates/jumphosts.html"))

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "tunnels", TemplateArg{
		Page:          "tunnels",
		Error:         err,
		CurrentTunnel: 0,
		Tunnels:       tunnels,
		Jumphosts:     jumphosts})
}

///////////////////////////
//
// Init helpers
//
///////////////////////////

func (t *Tunnelmgr) Config() {
	var addr = flag.String("l", ":3000", "Listening [<host>]:<port>")
	flag.Parse()

	t.cfg = TunnelmgrCfg{
		DBFilename: "tunnels.sqlite",
		Listener:   *addr,
	}
}

///////////////////////////
//
// Main Run
//
///////////////////////////

func (t *Tunnelmgr) Run() {
	t.Config()

	t.Open(t.cfg.DBFilename)
	defer t.Close()

	log.Printf("Starting server at %s", t.cfg.Listener)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", t.HandlerGetRoot)
	r.Get("/jumphosts", t.HandlerGetJumphosts)
	r.Delete("/jumphost/delete/{Id}", t.HandlerDeleteJumphosts)
	r.Get("/tunnel/get/{TunnelId}", t.HandlerGetTunnel)
	r.Post("/tunnel/update/{TunnelId}", t.HandlerPostUpdateTunnel)
	r.Post("/tunnel/add", t.HandlerPostAddTunnel)
	r.Delete("/tunnel/delete/{TunnelId}", t.HandlerDeleteTunnel)

	if err := http.ListenAndServe(t.cfg.Listener, r); err != nil {
		log.Fatal(err)
	}
}
