package main

import (
	"log"
	"os"
	"runtime"

	"github.com/alexflint/go-arg"
	"github.com/brettlangdon/realm"
)

var args struct {
	Zones   []string `arg:"--zone,positional,help:DNS zone files to serve from this server"`
	Bind    string   `arg:"help:[<host>]:<port> to bind too"`
	StatsD  string   `arg:"--statsd,help:<host>:<port> to send StatsD metrics to"`
	Workers int      `arg:"--workers,help:number of workers to start [default: $GOMAXPROCS]`
}

func main() {
	args.Bind = ":53"
	argParser := arg.MustParse(&args)

	// Control number of workers via GOMAXPROCS
	if args.Workers > 0 {
		runtime.GOMAXPROCS(args.Workers)
	}

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
	var server *realm.Server
	var err error
	server, err = realm.NewServer(args.Bind, registry, args.StatsD)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
