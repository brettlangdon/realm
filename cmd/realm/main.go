package main

import (
	"log"

	"github.com/brettlangdon/realm"
	"github.com/codegangsta/cli"
)

func main() {
	// Setup our CLI app
	var app *cli.App = cli.NewApp()
	app.Name = "realm"
	app.Version = "0.1.0"
	app.Author = "Brett Langdon"
	app.Email = "me@brett.is"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "zone, z",
			EnvVar: "REALM_ZONE",
			Usage:  "location to DNS zone file [required]",
		},
		cli.StringFlag{
			Name:   "bind, b",
			EnvVar: "REALM_BIND",
			Value:  ":53",
			Usage:  "'[<host>]:<port>' to bind too",
		},
	}

	// This action is called for all commands
	app.Action = func(c *cli.Context) {
		// Ensure that a zone filename was provided
		var filename string = c.String("zone")
		if filename == "" {
			log.Fatal("must supply zone file via \"--zone\" flag or \"REALM_ZONE\" environment variable")
		}

		// Load and parse the zone file
		var zone *realm.Zone
		var err error
		log.Printf("parsing zone file \"%s\"\n", filename)
		zone, err = realm.ParseZone(filename)
		if err != nil {
			log.Fatal(err)
		}

		// Create and start the server
		var bind string = c.String("bind")
		log.Printf("starting the server on \"%s\"\n", bind)
		var server *realm.Server = realm.NewServer(bind, zone)
		log.Fatal(server.ListenAndServe())
	}

	// Parse command arguments and run `app.Action`
	app.RunAndExitOnError()
}
