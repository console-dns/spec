package models

type Record struct {
	A     []*RecordA     `json:"a,omitempty"`
	AAAA  []*RecordAAAA  `json:"aaaa,omitempty"`
	TXT   []*RecordTXT   `json:"txt,omitempty"`
	CNAME []*RecordCNAME `json:"cname,omitempty"`
	NS    []*RecordNS    `json:"ns,omitempty"`
	MX    []*RecordMX    `json:"mx,omitempty"`
	SRV   []*RecordSRV   `json:"srv,omitempty"`
	CAA   []*RecordCAA   `json:"caa,omitempty"`
	SOA   *RecordSOA     `json:"soa,omitempty"`
}

func NewRecord() *Record {
	return &Record{
		A:     make([]*RecordA, 0),
		AAAA:  make([]*RecordAAAA, 0),
		TXT:   make([]*RecordTXT, 0),
		CNAME: make([]*RecordCNAME, 0),
		NS:    make([]*RecordNS, 0),
		MX:    make([]*RecordMX, 0),
		SRV:   make([]*RecordSRV, 0),
		CAA:   make([]*RecordCAA, 0),
		SOA:   nil,
	}
}

func (r *Record) Count() int {
	count := len(r.A) +
		len(r.AAAA) +
		len(r.TXT) +
		len(r.CNAME) +
		len(r.NS) +
		len(r.MX) +
		len(r.SRV) +
		len(r.CAA)
	if r.SOA != nil {
		count = count + 1
	}
	return count
}
