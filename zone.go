package realm

import (
	"fmt"
	"os"

	"github.com/miekg/dns"
)

// A Zone a container for records parsed from a zone file.
type Zone struct {
	records []dns.RR
}

// ParseZone will attempt to parse a zone file from the provided filename and return a Zone.
// ParseZone will return an error if the file provided does not exist or could not be properly parsed.
func ParseZone(filename string) (*Zone, error) {
	var zone *Zone
	var err error
	zone = &Zone{
		records: make([]dns.RR, 0),
	}

	// Open the file
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not parse zone file \"%s\": \"%s\"", filename, err)
	}
	defer file.Close()

	// Parse the file into records
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

// Lookup will find all records which we should respond with for the given name, request type, and request class.
func (zone *Zone) Lookup(name string, reqType uint16, reqClass uint16) []dns.RR {
	name = dns.Fqdn(name)
	var records []dns.RR
	records = make([]dns.RR, 0)
	for _, record := range zone.records {
		var header *dns.RR_Header
		header = record.Header()

		// Skip this record if the name or class do not match
		if header.Name != name || (header.Class != reqClass && reqClass != dns.ClassANY) {
			continue
		}

		// Collect this record if the types match or this record is a CNAME
		if reqType == dns.TypeANY || reqType == header.Rrtype {
			records = append(records, record)
		} else if header.Rrtype == dns.TypeCNAME {
			// Append this CNAME record as a response
			records = append(records, record)

			// Attempt to resolve this CNAME record
			var cname *dns.CNAME
			cname = record.(*dns.CNAME)
			var cnameRecords []dns.RR
			cnameRecords = zone.Lookup(dns.Fqdn(cname.Target), reqType, reqClass)
			records = append(records, cnameRecords...)
		}
	}
	return records
}
