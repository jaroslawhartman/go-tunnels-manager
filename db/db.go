package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Tunnel struct {
	JumphostId int
	Jumphost   string
	Command    string
	Name       string
	Local_port int
	Remote     string
	URL        string
	Status     int
}

type Tunnels map[int]*Tunnel

type Jumphost struct {
	Name    string
	Command string
}

type Jumphosts map[int]*Jumphost

type JumphostsTunnelsMatrix struct {
	Jumphost Jumphost
	Tunnels  []Tunnel
}

type JumphostsTunnelsMatrixList []JumphostsTunnelsMatrix

func Open(name string) *sql.DB {
	var initialiseDB = false

	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		initialiseDB = true
	}

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatal(err)
	}

	if initialiseDB {
		log.Println("Creating DB sructure")
		sqlStmt := `
		CREATE TABLE "jumphosts" (
			"id"	INTEGER NOT NULL,
			"name"	TEXT,
			"command"	TEXT,
			PRIMARY KEY("id" AUTOINCREMENT)
		);

		CREATE TABLE "tunnels" (
			"id"	INTEGER NOT NULL,
			"jumphost"	INTEGER,
			"name"	INTEGER,
			"local_port"	INTEGER UNIQUE,
			"remote"	INTEGER,
			"url"	TEXT,
			"status"	INTEGER,
			PRIMARY KEY("id" AUTOINCREMENT),
			CONSTRAINT fk_jumphosts
			FOREIGN KEY("jumphost") REFERENCES jumphosts("id") ON DELETE CASCADE
		);

		INSERT INTO jumphosts (name, command) VALUES ("PreProd", "ssh -N -L3000:localhost:3000 192.168.1.200");
		INSERT INTO jumphosts (name, command) VALUES ("Prod", "gcloud compute ssh postman-vm --configuration tef-cloudlab2 -- -NL8080:172.29.17.197:80");
		INSERT INTO jumphosts (name, command) VALUES ("Dev", "gcloud compute ssh postman-vm --configuration tef-cloudlab2 -- -NL8080:172.29.17.197:80");

		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (1, "Grafana", 8000, "192.168.1.100:9900", "http://localhost:8000", 1);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (1, "Prometheus", 8001, "192.168.1.101:9900", "http://localhost:8001", 0);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (1, "Kafka UI", 8002, "192.168.1.102:9900", "http://localhost:8002", 1);

		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (2, "Grafana", 8010, "192.168.1.100:9900", "http://localhost:8000", 1);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (2, "Prometheus", 8011, "192.168.1.101:9900", "http://localhost:8001", 0);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (2, "Kafka UI", 8012, "192.168.1.102:9900", "http://localhost:8002", 0);

		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (3, "Grafana", 8020, "192.168.1.100:9900", "http://localhost:8020", 0);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (3, "Prometheus", 8021, "192.168.1.101:9900", "http://localhost:8021", 1);
		INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) VALUES (3, "Kafka UI", 8022, "192.168.1.102:9900", "http://localhost:8022", 1);
		`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			log.Fatal(err)
		}
	}

	return db
}

func Close(db *sql.DB) {
	db.Close()
}

func GetTunnels(db *sql.DB, filter string) Tunnels {
	tunnels := make(Tunnels)

	rows, err := db.Query(fmt.Sprintf(`SELECT t.id as id, j.id as jid, j.name as jname, j.command as jcommand,
		t.name, t.local_port, t.remote, t.url, t.status
		from tunnels t, jumphosts j
		where t.jumphost=j.id %s`, filter))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var jid int
		var jname string
		var jcommand string
		var name string
		var local_port int
		var remote string
		var url string
		var status int

		err = rows.Scan(&id, &jid, &jname, &jcommand, &name, &local_port, &remote, &url, &status)
		if err != nil {
			log.Fatal(err)
		}

		tunnels[id] = &Tunnel{
			JumphostId: jid,
			Jumphost:   jname,
			Command:    jcommand,
			Name:       name,
			Local_port: local_port,
			Remote:     remote,
			URL:        url,
			Status:     status,
		}
	}

	return tunnels
}

func UpdateTunnel(db *sql.DB, id int, tunnel *Tunnel) string {
	stmt, _ := db.Prepare("UPDATE tunnels set jumphost = ?, name = ?, local_port = ?, remote = ?, url = ? WHERE id = ?")
	defer stmt.Close()

	_, err := stmt.Exec(tunnel.JumphostId,
		tunnel.Name,
		tunnel.Local_port,
		tunnel.Remote,
		tunnel.URL,
		id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return ""
}

func AddTunnel(db *sql.DB, tunnel *Tunnel) string {
	stmt, _ := db.Prepare("INSERT INTO tunnels (jumphost, name, local_port, remote, url, status) values (?, ?, ?, ?, ?, ?)")
	defer stmt.Close()

	_, err := stmt.Exec(tunnel.JumphostId,
		tunnel.Name,
		tunnel.Local_port,
		tunnel.Remote,
		tunnel.URL,
		0)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return ""
}

func DeleteTunnel(db *sql.DB, id int) string {
	stmt, err := db.Prepare("DELETE FROM tunnels WHERE id = ?")

	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return ""
}

func GetJumphosts(db *sql.DB) Jumphosts {
	jumphosts := make(Jumphosts)

	rows, err := db.Query(`SELECT id, name, command	 from jumphosts`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var command string

		err = rows.Scan(&id, &name, &command)
		if err != nil {
			log.Fatal(err)
		}

		jumphosts[id] = &Jumphost{
			Name:    name,
			Command: command,
		}
	}

	return jumphosts
}

func DeleteJumphost(db *sql.DB, id int) string {
	stmt, err := db.Prepare("DELETE FROM jumphosts WHERE id = ?")

	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return ""
}

func GetJumphostsTunnelsMatrix(db *sql.DB) JumphostsTunnelsMatrixList {
	var output JumphostsTunnelsMatrixList

	jumphosts := GetJumphosts(db)

	for ji, j := range jumphosts {
		log.Println("Jumphosts ")
		log.Println(ji, j)

		tunnelsList := make([]Tunnel, 0)
		tunnels := GetTunnels(db, fmt.Sprintf(" AND jid=%d ", ji))

		for ti, t := range tunnels {
			tunnelsList = append(tunnelsList, Tunnel{
				JumphostId: t.JumphostId,
				Jumphost:   t.Jumphost,
				Command:    t.Command,
				Name:       t.Name,
				Local_port: t.Local_port,
				Remote:     t.Remote,
				URL:        t.URL,
				Status:     t.Status,
			})

			log.Println("Tunnel ")
			log.Println(ti, t)
		}

		matrix := JumphostsTunnelsMatrix{
			Jumphost: Jumphost{
				Name:    j.Name,
				Command: j.Command,
			},
			Tunnels: tunnelsList,
		}

		output = append(output, matrix)
	}

	return output
}
