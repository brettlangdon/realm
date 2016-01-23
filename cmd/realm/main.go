package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/brettlangdon/realm"
)

var args struct {
	Zones []string `arg:"--zone,positional,help:DNS zone files to serve from this server"`
	Bind  string   `arg:"help:[<host>]:<port> to bind too"`
}

func main() {
	args.Bind = ":53"
	argParser := arg.MustParse(&args)

	if len(args.Zones) == 0 {
		log.Println("must supply at least 1 zone file to serve")
		argParser.WriteUsage(os.Stderr)
		os.Exit(1)
	}

	var registry *realm.Registry
	registry = realm.NewRegistry()

	for _, filename := range args.Zones {
		// Load and parse the zone file
		var zone *realm.Zone
		var err error
		log.Printf("parsing zone file \"%s\"\n", filename)
		zone, err = realm.ParseZone(filename)
		if err != nil {
			log.Fatal(err)
		}
		registry.AddZone(zone)
	}

	// Create and start the server
	log.Printf("starting the server on \"%s\"\n", args.Bind)
	var server *realm.Server = realm.NewServer(args.Bind, registry)
	log.Fatal(server.ListenAndServe())
}
