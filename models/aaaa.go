package models

import (
	"net"
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordAAAA struct {
	Ttl uint32 `json:"ttl" yaml:"ttl" toml:"ttl"`
	Ip  net.IP `json:"ip" yaml:"ip" toml:"ip"`
}

func (record *RecordAAAA) Clone() *RecordAAAA {
	return &RecordAAAA{
		Ttl: record.Ttl,
		Ip:  record.Ip,
	}
}

func FromRecordAAAA(values GetValue) (*RecordAAAA, error) {
	return NewRecordAAAA(values("ip"), values("ttl"))
}
func NewRecordAAAA(ip, ttl string) (*RecordAAAA, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexIPv6.Valid(ip)
	if err != nil {
		return nil, err
	}
	return &RecordAAAA{
		Ttl: u,
		Ip:  net.ParseIP(ip),
	}, nil
}

func (r *Record) IndexAAAA(other *RecordAAAA) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.AAAA {
		if a.Ip.Equal(other.Ip) && a.Ttl == other.Ttl {
			return i
		}
	}
	return scanner.EOF
}

func (r *Record) ModAAAA(old, data *RecordAAAA) error {
	if r.IndexAAAA(data) != scanner.EOF {
		return errors.New("已存在相同的 AAAA 记录")
	}
	if old != nil {
		oldIndex := r.IndexAAAA(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 AAAA 记录不存在")
		}
		r.AAAA[oldIndex] = data
	} else {
		r.AAAA = append(r.AAAA, data)
	}
	return nil
}

func (r *Record) RemoveAAAA(data *RecordAAAA) error {
	index := r.IndexAAAA(data)
	if index == scanner.EOF {
		return errors.New("当前 AAAA 记录不存在")
	}
	r.AAAA, _ = utils.RemoveIndex(r.AAAA, index)
	return nil
}
