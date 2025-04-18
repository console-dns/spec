package models

import (
	"net"
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordA struct {
	Ttl uint32 `json:"ttl" yaml:"ttl" toml:"ttl"`
	Ip  net.IP `json:"ip" yaml:"ip" toml:"ip"`
}

func (record *RecordA) Clone() *RecordA {
	return &RecordA{
		Ttl: record.Ttl,
		Ip:  record.Ip,
	}
}

func FromRecordA(values GetValue) (*RecordA, error) {
	return NewRecordA(values("ip"), values("ttl"))
}
func NewRecordA(ip, ttl string) (*RecordA, error) {
	t, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexIPv4.Valid(ip)
	if err != nil {
		return nil, err
	}
	return &RecordA{
		Ttl: t,
		Ip:  net.ParseIP(ip),
	}, nil
}

func (r *Record) IndexA(other *RecordA) int {
	if other == nil {
		return -1
	}
	for i, a := range r.A {
		if a.Ip.Equal(other.Ip) && a.Ttl == other.Ttl {
			return i
		}
	}
	return -1
}

func (r *Record) ModA(old, data *RecordA) error {
	if r.IndexA(data) != scanner.EOF {
		return errors.New("已存在相同的 A 记录")
	}
	if old != nil {
		oldIndex := r.IndexA(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 A 记录不存在")
		}
		r.A[oldIndex] = data
	} else {
		r.A = append(r.A, data)
	}
	return nil
}

func (r *Record) RemoveA(data *RecordA) error {
	index := r.IndexA(data)
	if index == scanner.EOF {
		return errors.New("当前 A 记录不存在")
	}
	r.A, _ = utils.RemoveIndex(r.A, index)
	return nil
}
