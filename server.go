package realm

import "github.com/miekg/dns"

// A Server listens for DNS requests over UDP and responds with answers from the provided Zone.
type Server struct {
	server   *dns.Server
	registry *Registry
}

// NewServer returns a new initialized *Server that will bind to listen and will look up answers from zone.
func NewServer(listen string, registry *Registry) *Server {
	var s *Server
	s = &Server{registry: registry}
	s.server = &dns.Server{
		Addr:    listen,
		Net:     "udp",
		Handler: s,
	}
	return s
}

// ListenAndServe will start the nameserver on the configured address.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
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
		records = s.registry.Lookup(question.Name, question.Qtype, question.Qclass)
		response.Answer = append(response.Answer, records...)
	}

	// Respond to the request
	w.WriteMsg(response)
}
