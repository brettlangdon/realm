package realm

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

type Zone struct {
	records []dns.RR
}

func ParseZone(filename string) (*Zone, error) {
	var zone *Zone
	var err error
	zone = &Zone{
		records: make([]dns.RR, 0),
	}

	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not parse zone file \"%s\": \"%s\"", filename, err)
	}
	defer file.Close()

	var tokens chan *dns.Token
	tokens = dns.ParseZone(file, "", "")
	for token := range tokens {
		if token.Error != nil {
			return nil, fmt.Errorf("could not parse zone file \"%s\": \"%s\"", filename, token.Error)
		}

		zone.records = append(zone.records, token.RR)
	}
	return zone, nil
}

func (zone *Zone) Lookup(name string, reqType uint16, reqClass uint16) []dns.RR {
	name = dns.Fqdn(name)
	var records []dns.RR = make([]dns.RR, 0)
	for _, record := range zone.records {
		var header *dns.RR_Header = record.Header()
		if header.Name != name || (header.Class != reqClass && reqClass != dns.ClassANY) {
			continue
		}

		if reqType == dns.TypeANY || reqType == header.Rrtype {
			records = append(records, record)
		} else if header.Rrtype == dns.TypeCNAME {
			records = append(records, record)
			var cname *dns.CNAME = record.(*dns.CNAME)
			var cnameRecords []dns.RR = zone.Lookup(dns.Fqdn(cname.Target), reqType, reqClass)
			records = append(records, cnameRecords...)
		}
	}
	return records
}
