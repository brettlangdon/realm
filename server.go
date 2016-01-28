package realm

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/miekg/dns"
)

// A Server listens for DNS requests over UDP and responds with answers from the provided Zone.
type Server struct {
	server   *dns.Server
	registry *Registry
	statsd   *statsd.Client
}

// NewServer returns a new initialized *Server that will bind to listen and will look up answers from zone.
func NewServer(listen string, registry *Registry, statsdHost string) (*Server, error) {
	var err error
	var s *Server
	s = &Server{registry: registry}
	s.server = &dns.Server{
		Addr:    listen,
		Net:     "udp",
		Handler: s,
	}
	if statsdHost != "" {
		s.statsd, err = statsd.New(statsdHost)
		if err != nil {
			return nil, err
		}
		s.statsd.Namespace = "realm."
	}
	return s, err
}

// ListenAndServe will start the nameserver on the configured address.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// ServeDNS will be called for every DNS request to this server.
// It will attempt to provide answers to all questions from the configured zone.
func (s *Server) ServeDNS(w dns.ResponseWriter, request *dns.Msg) {
	// Call `Hijack` since we will handle closing `dns.ResponseWriter` ourselves
	w.Hijack()
	// Handle the request
	go s.handle(w, request)
}

func (s *Server) handle(w dns.ResponseWriter, request *dns.Msg) {
	// Always close the writer
	defer w.Close()

	// Capture starting time for measuring message response time
	var start time.Time
	start = time.Now()

	// Setup the default response
	var response *dns.Msg
	response = &dns.Msg{}
	response.SetReply(request)
	response.Compress = true

	// Lookup answers to any of the questions
	for _, question := range request.Question {
		// Capture starting time for measuring lookup
		var lookupStart time.Time
		lookupStart = time.Now()

		// Perform lookup for this question
		var records []dns.RR
		records = s.registry.Lookup(question.Name, question.Qtype, question.Qclass)

		// Capture ending and elapsed time
		var lookupElapsed time.Duration
		lookupElapsed = time.Since(lookupStart)

		// Append results to the response
		response.Answer = append(response.Answer, records...)

		// If StatsD is enabled, record some metrics
		if s.statsd != nil {
			var tags []string
			tags = []string{
				fmt.Sprintf("name:%s", question.Name),
				fmt.Sprintf("qtype:%s", dns.TypeToString[question.Qtype]),
				fmt.Sprintf("qclass:%s", dns.ClassToString[question.Qclass]),
			}

			s.statsd.TimeInMilliseconds("lookup.time", lookupElapsed.Seconds()*1000.0, tags, 1)
			s.statsd.Histogram("lookup.answer", float64(len(records)), tags, 1)
			s.statsd.Count("request.question", 1, tags, 1)
		}
	}

	// Respond to the request
	w.WriteMsg(response)

	// Record any ending metrics
	if s.statsd != nil {
		var elapsed time.Duration
		elapsed = time.Since(start)
		s.statsd.TimeInMilliseconds("request.time", elapsed.Seconds()*1000.0, nil, 1)
	}
}
