package main

import (
	"log"

	"github.com/brettlangdon/realm"
	"github.com/codegangsta/cli"
)

func main() {
	// Setup our CLI app
	var app *cli.App
	app = cli.NewApp()
	app.Name = "realm"
	app.Usage = "A simple non-recursive DNS server"
	app.Version = "0.1.0"
	app.Author = "Brett Langdon"
	app.Email = "me@brett.is"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "zone, z",
			EnvVar: "REALM_ZONE",
			Value:  &cli.StringSlice{},
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
		var filenames []string
		filenames = c.StringSlice("zone")
		if len(filenames) == 0 {
			log.Fatal("must supply at least 1 zone file via \"--zone\" flag or \"REALM_ZONE\" environment variable")
		}

		// Load and parse the zone file
		var zones realm.Zones
		zones = make(realm.Zones, 0)

		var err error
		for _, filename := range filenames {
			log.Printf("parsing zone file \"%s\"\n", filename)
			var zone *realm.Zone
			zone, err = realm.ParseZone(filename)
			if err != nil {
				log.Fatal(err)
			}
			zones = append(zones, zone)
		}

		// Create and start the server
		var bind string
		bind = c.String("bind")
		log.Printf("starting the server on \"%s\"\n", bind)
		var server *realm.Server
		server = realm.NewServer(bind, zones)
		log.Fatal(server.ListenAndServe())
	}

	// Parse command arguments and run `app.Action`
	app.RunAndExitOnError()
}
