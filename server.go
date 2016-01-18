package realm

import "github.com/miekg/dns"

type Server struct {
	server *dns.Server
	zone   *Zone
}

func NewServer(listen string, zone *Zone) *Server {
	var server *Server = &Server{zone: zone}
	server.server = &dns.Server{
		Addr:    listen,
		Net:     "udp",
		Handler: server,
	}
	return server
}

func (server *Server) ListenAndServe() error {
	return server.server.ListenAndServe()
}

func (server *Server) ServeDNS(w dns.ResponseWriter, request *dns.Msg) {
	var response *dns.Msg = &dns.Msg{}
	response.SetReply(request)
	response.Compress = true

	for _, question := range request.Question {
		var records []dns.RR = server.zone.Lookup(question.Name, question.Qtype, question.Qclass)
		response.Answer = append(response.Answer, records...)
	}

	w.WriteMsg(response)
}
