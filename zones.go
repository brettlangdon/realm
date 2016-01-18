package realm

import (
	"github.com/miekg/dns"
)

// Zones is a convenient helper for managing a slice of *Zone.
type Zones []*Zone

// Lookup will find all answer records from across all *Zone
func (z Zones) Lookup(name string, reqType uint16, reqClass uint16) []dns.RR {
	var records []dns.RR
	records = make([]dns.RR, 0)
	for _, zone := range z {
		records = append(records, zone.Lookup(name, reqType, reqClass)...)
	}
	return records
}
