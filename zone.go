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

func (zone *Zone) Records() []dns.RR {
	return zone.records
}
