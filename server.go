package realm

import "github.com/miekg/dns"

// A Server listens for DNS requests over UDP and responds with answers from the provided Zone.
type Server struct {
	server *dns.Server
	zone   *Zone
}

// NewServer returns a new initialized *Server that will bind to listen and will look up answers from zone.
func NewServer(listen string, zone *Zone) *Server {
	var server *Server
	s = &Server{zone: zone}
	s.server = &dns.Server{
		Addr:    listen,
		Net:     "udp",
		Handler: server,
	}
	return s
}

// ListenAndServe will start the nameserver on the configured address.
func (s *Server) ListenAndServe() error {
	return server.server.ListenAndServe()
}

// ServeDNS will be called for every DNS request to this server.
// It will attempt to provide answers to all questions from the configured zone.
func (s *Server) ServeDNS(w dns.ResponseWriter, request *dns.Msg) {
	// Setup the default response
	var response *dns.Msg
	response = &dns.Msg{}
	response.SetReply(request)
	response.Compress = true

	// Lookup answers to any of the questions
	for _, question := range request.Question {
		var records []dns.RR
		records = s.zone.Lookup(question.Name, question.Qtype, question.Qclass)
		response.Answer = append(response.Answer, records...)
	}

	// Respond to the request
	w.WriteMsg(response)
}
