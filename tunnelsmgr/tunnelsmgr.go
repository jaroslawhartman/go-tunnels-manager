package tunnelsmgr

import (
	"database/sql"
	"embed"
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
	Page            string
	Error           string
	JumphostId      int
	TunnelId        int
	CurrentTunnel   int
	CurrentJumphost int
	Jumphosts       Jumphosts
	Tunnels         Tunnels
}

var (
	//go:embed templates/*
	files embed.FS
)

///////////////////////////
//
// Chi Handlers
//
///////////////////////////

func (t *Tunnelmgr) HandlerGetRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /")
	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	tmpl.Execute(w, TemplateArg{Page: "tunnels", CurrentTunnel: 0, Jumphosts: jumphosts, Tunnels: tunnels})
}

///////////////////////
//
// Jumphosts Handlers
//
///////////////////////

func (t *Tunnelmgr) HandlerGetJumphosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /jumphosts")
	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.Execute(w, TemplateArg{Page: "jumphosts", Jumphosts: jumphosts, Tunnels: nil})
}

func (t *Tunnelmgr) HandlerDeleteJumphosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /jumphost/delete/")
	id, _ := strconv.Atoi(chi.URLParam(r, "Id"))
	log.Println(id)

	err := t.DeleteJumphost(id)

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "jumphosts", TemplateArg{
		Page:      "jumphosts",
		Error:     err,
		Jumphosts: jumphosts,
		Tunnels:   nil})
}

func (t *Tunnelmgr) HandlerGetJumphost(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /jumphost/get/")
	jumphostId, _ := strconv.Atoi(chi.URLParam(r, "JumphostId"))

	log.Println(jumphostId)

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	tunnels := t.GetTunnels("")
	jumphosts := t.GetJumphosts()

	w.Header().Set("HX-Trigger-After-Settle", "showModal")

	tmpl.ExecuteTemplate(w, "jumphost-edit-modal", TemplateArg{
		CurrentJumphost: jumphostId,
		Tunnels:         tunnels,
		Jumphosts:       jumphosts})
}

func (t *Tunnelmgr) HandlerPostUpdateJumphost(w http.ResponseWriter, r *http.Request) {
	log.Println("Post /jumphost/update/")
	jumphostId, _ := strconv.Atoi(chi.URLParam(r, "JumphostId"))
	log.Println(jumphostId)

	inputName := r.FormValue("inputName")
	inputCommand := r.FormValue("inputCommand")

	err := t.UpdateJumphost(jumphostId, &Jumphost{
		Name:    inputName,
		Command: inputCommand,
	})

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "jumphosts", TemplateArg{
		Page:      "jumphosts",
		Error:     err,
		Jumphosts: jumphosts,
		Tunnels:   nil})
}

func (t *Tunnelmgr) HandlerPostAddJumphost(w http.ResponseWriter, r *http.Request) {
	log.Println("Post /jumphost/add/")

	inputName := r.FormValue("inputName")
	inputCommand := r.FormValue("inputCommand")

	err := t.AddJumphost(&Jumphost{
		Name:    inputName,
		Command: inputCommand,
	})

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	jumphosts := t.GetJumphosts()

	tmpl.ExecuteTemplate(w, "jumphosts", TemplateArg{
		Page:      "jumphosts",
		Error:     err,
		Jumphosts: jumphosts,
		Tunnels:   nil})
}

///////////////////////
//
// Tunnels Handlers
//
///////////////////////

func (t *Tunnelmgr) HandlerGetTunnel(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /tunnel/get/")
	tunnelId, _ := strconv.Atoi(chi.URLParam(r, "TunnelId"))

	log.Println(tunnelId)

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
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

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
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

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
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

	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
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

///////////////////////
//
// Settings Handlers
//
///////////////////////

func (t *Tunnelmgr) HandlerGetSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("Get /settings")
	tmpl := template.Must(template.ParseFS(files,
		"templates/index.html",
		"templates/tunnels.html",
		"templates/settings.html",
		"templates/jumphosts.html"))

	// jumphosts := t.GetJumphosts()

	tmpl.Execute(w, TemplateArg{Page: "settings", Jumphosts: nil, Tunnels: nil})
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

	t.OpenDB()
	defer t.CloseDB()

	log.Printf("Starting server at %s", t.cfg.Listener)

	// go t.Watchdog()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", t.HandlerGetRoot)
	r.Get("/jumphosts", t.HandlerGetJumphosts)
	r.Get("/jumphost/get/{JumphostId}", t.HandlerGetJumphost)
	r.Post("/jumphost/update/{JumphostId}", t.HandlerPostUpdateJumphost)
	r.Post("/jumphost/add", t.HandlerPostAddJumphost)
	r.Delete("/jumphost/delete/{Id}", t.HandlerDeleteJumphosts)
	r.Get("/tunnel/get/{TunnelId}", t.HandlerGetTunnel)
	r.Post("/tunnel/update/{TunnelId}", t.HandlerPostUpdateTunnel)
	r.Post("/tunnel/add", t.HandlerPostAddTunnel)
	r.Delete("/tunnel/delete/{TunnelId}", t.HandlerDeleteTunnel)
	r.Get("/settings", t.HandlerGetSettings)

	if err := http.ListenAndServe(t.cfg.Listener, r); err != nil {
		log.Fatal(err)
	}
}

// func (t *Tunnelmgr) Watchdog() {

// 	for {
// 		// Loop every 10 seconds
// 		time.Sleep(10 * time.Second)

// 		// tunnels := t.GetTunnels("")

// 		log.Println("...watchdog")
// 	}
// }
