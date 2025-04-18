package models

import (
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordCNAME struct {
	Ttl  uint32 `json:"ttl"`
	Host string `json:"host"`
}

func (r *RecordCNAME) Clone() *RecordCNAME {
	return &RecordCNAME{
		Ttl:  r.Ttl,
		Host: r.Host,
	}
}

func FromRecordCNAME(values GetValue) (*RecordCNAME, error) {
	return NewRecordCNAME(values("host"), values("ttl"))
}
func NewRecordCNAME(host, ttl string) (*RecordCNAME, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexHost.Valid(host)
	if err != nil {
		return nil, err
	}
	return &RecordCNAME{
		Ttl:  u,
		Host: host,
	}, nil
}

func (r *Record) IndexCNAME(other *RecordCNAME) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.CNAME {
		if a.Host == other.Host && a.Ttl == other.Ttl {
			return i
		}
	}
	return scanner.EOF
}
func (r *Record) ModCNAME(old, data *RecordCNAME) error {
	if r.IndexCNAME(data) != scanner.EOF {
		return errors.New("已存在相同的 CNAME 记录")
	}
	if old != nil {
		oldIndex := r.IndexCNAME(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 CNAME 记录不存在")
		}
		r.CNAME[oldIndex] = data
	} else {
		r.CNAME = append(r.CNAME, data)
	}
	return nil
}

func (r *Record) RemoveCNAME(data *RecordCNAME) error {
	index := r.IndexCNAME(data)
	if index == scanner.EOF {
		return errors.New("当前 CNAME 记录不存在")
	}
	r.CNAME, _ = utils.RemoveIndex(r.CNAME, index)
	return nil
}
