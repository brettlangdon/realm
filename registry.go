package realm

import "github.com/miekg/dns"

type RecordsEntry map[uint16][]dns.RR

func (entry RecordsEntry) GetRecords(rrType uint16) []dns.RR {
	var records []dns.RR
	records = make([]dns.RR, 0)

	if rrType == dns.TypeANY {
		for _, rrs := range entry {
			records = append(records, rrs...)
		}
	} else if rrs, ok := entry[rrType]; ok {
		records = append(records, rrs...)
	}

	return records
}

type DomainEntry map[uint16]RecordsEntry

func (entry DomainEntry) AddEntry(record dns.RR) {
	var header *dns.RR_Header
	header = record.Header()

	if _, ok := entry[header.Class]; !ok {
		entry[header.Class] = make(RecordsEntry)
	}
	if _, ok := entry[header.Class][header.Rrtype]; !ok {
		entry[header.Class][header.Rrtype] = make([]dns.RR, 0)
	}

	entry[header.Class][header.Rrtype] = append(entry[header.Class][header.Rrtype], record)
}

func (entry DomainEntry) GetEntries(rrClass uint16) []RecordsEntry {
	var entries []RecordsEntry
	entries = make([]RecordsEntry, 0)

	if rrClass == dns.ClassANY {
		for _, entry := range entry {
			entries = append(entries, entry)
		}
	} else if entry, ok := entry[rrClass]; ok {
		entries = append(entries, entry)
	}

	return entries
}

type Registry struct {
	records map[string]DomainEntry
}

func NewRegistry() *Registry {
	return &Registry{
		records: make(map[string]DomainEntry),
	}
}

func (r *Registry) addRecord(record dns.RR) {
	var header *dns.RR_Header
	header = record.Header()

	var name string
	name = dns.Fqdn(header.Name)

	if _, ok := r.records[name]; !ok {
		r.records[name] = make(DomainEntry)
	}
	r.records[name].AddEntry(record)

	// If this record is an SOA record then also store under the Mbox name
	if header.Rrtype == dns.TypeSOA {
		var soa *dns.SOA
		soa = record.(*dns.SOA)

		if _, ok := r.records[soa.Mbox]; !ok {
			r.records[soa.Mbox] = make(DomainEntry)
		}
		r.records[soa.Mbox].AddEntry(record)
	}
}

func (r *Registry) AddZone(z *Zone) {
	for _, record := range z.Records() {
		r.addRecord(record)
	}
}

// Lookup will find all records which we should respond with for the given name, request type, and request class.
func (r *Registry) Lookup(name string, reqType uint16, reqClass uint16) []dns.RR {
	name = dns.Fqdn(name)
	var records []dns.RR
	records = make([]dns.RR, 0)

	var domainEntry DomainEntry
	var ok bool
	domainEntry, ok = r.records[name]
	if !ok {
		return records
	}

	var recordEntries []RecordsEntry
	recordEntries = domainEntry.GetEntries(reqClass)

	for _, recordEntry := range recordEntries {
		var rrs []dns.RR
		rrs = recordEntry.GetRecords(reqType)
		records = append(records, rrs...)

		if len(rrs) == 0 && reqType == dns.TypeA {
			rrs = recordEntry.GetRecords(dns.TypeCNAME)
			for _, rr := range rrs {
				records = append(records, rr)
				var header *dns.RR_Header
				header = rr.Header()
				if header.Rrtype == dns.TypeCNAME && reqType != dns.TypeCNAME {
					// Attempt to resolve this CNAME record
					var cname *dns.CNAME
					cname = rr.(*dns.CNAME)
					var cnameRecords []dns.RR
					cnameRecords = r.Lookup(dns.Fqdn(cname.Target), reqType, reqClass)
					records = append(records, cnameRecords...)
				}

			}
		}

	}
	return records
}
