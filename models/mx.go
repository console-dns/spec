package models

import (
	"math"
	"strconv"
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordMX struct {
	Ttl        uint32 `json:"ttl"`
	Host       string `json:"host"`
	Preference uint16 `json:"preference"`
}

func (r *RecordMX) Clone() *RecordMX {
	return &RecordMX{
		Ttl:        r.Ttl,
		Host:       r.Host,
		Preference: r.Preference,
	}
}

func FromRecordMX(values GetValue) (*RecordMX, error) {
	return NewRecordMX(values("host"), values("preference"), values("ttl"))
}
func NewRecordMX(host, preference, ttl string) (*RecordMX, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexHost.Valid(host)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(preference)
	if err != nil {
		return nil, err
	}
	if p > math.MaxUint16 || p < 0 {
		return nil, errors.New("preference 区间不合法")
	}
	return &RecordMX{
		Ttl:        u,
		Host:       host,
		Preference: uint16(p),
	}, nil
}

func (r *Record) IndexMX(other *RecordMX) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.MX {
		if a.Host == other.Host && a.Ttl == other.Ttl && other.Preference == a.Preference {
			return i
		}
	}
	return scanner.EOF
}

func (r *Record) ModMX(old, data *RecordMX) error {
	if r.IndexMX(data) != scanner.EOF {
		return errors.New("已存在相同的 MX 记录")
	}
	if old != nil {
		oldIndex := r.IndexMX(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 MX 记录不存在")
		}
		r.MX[oldIndex] = data
	} else {
		r.MX = append(r.MX, data)
	}
	return nil
}

func (r *Record) RemoveMX(data *RecordMX) error {
	index := r.IndexMX(data)
	if index == scanner.EOF {
		return errors.New("当前 MX 记录不存在")
	}
	r.MX, _ = utils.RemoveIndex(r.MX, index)
	return nil
}
