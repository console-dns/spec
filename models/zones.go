package models

import (
	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type Zones struct {
	Zones map[string]*Zone `json:"zones" yaml:"zones" toml:"zones"`
}

func (z *Zones) ListZones() []string {
	var zones = make([]string, 0, len(z.Zones))
	for zone := range z.Zones {
		zones = append(zones, zone)
	}
	return zones
}

func (z *Zones) ListRecords() map[string]map[string]Record {
	var result = make(map[string]map[string]Record)
	for zone, records := range z.Zones {
		result[zone] = make(map[string]Record)
		for r, rd := range records.Records {
			result[zone][r] = *rd
		}
	}
	return result
}

func (z *Zones) GetZone(name string) *Zone {
	return z.Zones[name]
}

func (z *Zones) GetRecords(zone string) map[string]Record {
	result := make(map[string]Record)
	if z.GetZone(zone) == nil {
		return result
	}
	for s, record := range z.Zones[zone].Records {
		result[s] = *record
	}
	return result
}

// Clean 清理
func (z *Zones) Clean(zone bool) {
	rzs := make([]string, 0)
	for rzn, response := range z.Zones {
		rds := make([]string, 0)
		for rn, record := range response.Records {
			if record.Count() == 0 {
				rds = append(rds, rn)
			}
		}
		for _, rd := range rds {
			delete(response.Records, rd)
		}
		if len(response.Records) == 0 {
			rzs = append(rzs, rzn)
		}
	}
	if zone {
		for _, rz := range rzs {
			delete(z.Zones, rz)
		}
	}
}

func NewZones() *Zones {
	return &Zones{
		Zones: make(map[string]*Zone),
	}
}

func (z *Zones) AddZone(zone string) error {
	err := utils.RegexHost.Valid(zone)
	if err != nil {
		return errors.New("区域名称不合法")
	}
	if _, ok := z.Zones[zone]; ok {
		return errors.New("区域已存在")
	}
	z.Zones[zone] = NewZone()
	return nil
}

func (z *Zones) RemoveZone(zone string) error {
	err := utils.RegexHost.Valid(zone)
	if err != nil {
		return errors.New("区域名称不合法")
	}
	if _, ok := z.Zones[zone]; !ok {
		return errors.New("区域不存在")
	}
	delete(z.Zones, zone)
	return nil
}

func (z *Zones) CopyFrom(src *Record, zone string, record string, rType string) {
	if z.Zones[zone] == nil {
		z.Zones[zone] = NewZone()
	}
	if z.Zones[zone].Records[record] == nil {
		z.Zones[zone].Records[record] = NewRecord()
	}
	switch rType {
	case "A":
		z.Zones[zone].Records[record].A = fastCopy(src.A)
	case "AAAA":
		z.Zones[zone].Records[record].AAAA = fastCopy(src.AAAA)
	case "TXT":
		z.Zones[zone].Records[record].TXT = fastCopy(src.TXT)
	case "CNAME":
		z.Zones[zone].Records[record].CNAME = fastCopy(src.CNAME)
	case "NS":
		z.Zones[zone].Records[record].NS = fastCopy(src.NS)
	case "MX":
		z.Zones[zone].Records[record].MX = fastCopy(src.MX)
	case "SRV":
		z.Zones[zone].Records[record].SRV = fastCopy(src.SRV)
	case "CAA":
		z.Zones[zone].Records[record].CAA = fastCopy(src.CAA)
	case "SOA":
		z.Zones[zone].Records[record].SOA = *&src.SOA
	}
}

func fastCopy[V Clone[V], T []V](src T) T {
	dst := make(T, len(src))
	for i, a := range src {
		dst[i] = a.Clone()
	}
	return dst
}
