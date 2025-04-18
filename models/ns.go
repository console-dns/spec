package models

import (
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordNS struct {
	Ttl  uint32 `json:"ttl"`
	Host string `json:"host"`
}

func (r *RecordNS) Clone() *RecordNS {
	return &RecordNS{
		Ttl:  r.Ttl,
		Host: r.Host,
	}
}

func FromRecordNS(values GetValue) (*RecordNS, error) {
	return NewRecordNS(values("host"), values("ttl"))
}
func NewRecordNS(host, ttl string) (*RecordNS, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexHost.Valid(host)
	if err != nil {
		return nil, err
	}
	return &RecordNS{
		Ttl:  u,
		Host: host,
	}, nil
}

func (r *Record) IndexNS(other *RecordNS) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.NS {
		if a.Host == other.Host && a.Ttl == other.Ttl {
			return i
		}
	}
	return scanner.EOF
}

func (r *Record) ModNS(old, data *RecordNS) error {
	if r.IndexNS(data) != scanner.EOF {
		return errors.New("已存在相同的 NS 记录")
	}
	if old != nil {
		oldIndex := r.IndexNS(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 NS 记录不存在")
		}
		r.NS[oldIndex] = data
	} else {
		r.NS = append(r.NS, data)
	}
	return nil
}

func (r *Record) RemoveNS(data *RecordNS) error {
	index := r.IndexNS(data)
	if index == scanner.EOF {
		return errors.New("当前 NS 记录不存在")
	}
	r.NS, _ = utils.RemoveIndex(r.NS, index)
	return nil
}
