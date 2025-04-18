package models

import (
	"math"
	"strconv"
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordSRV struct {
	Ttl      uint32 `json:"ttl"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

func (r *RecordSRV) Clone() *RecordSRV {
	return &RecordSRV{
		Ttl:      r.Ttl,
		Priority: r.Priority,
		Weight:   r.Weight,
		Port:     r.Port,
		Target:   r.Target,
	}
}

func FromRecordSRV(values GetValue) (*RecordSRV, error) {
	return NewRecordSRV(
		values("priority"),
		values("weight"),
		values("port"),
		values("target"),
		values("ttl"))
}
func NewRecordSRV(priority, weight, port, target, ttl string) (*RecordSRV, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	err = utils.RegexHost.Valid(target)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(priority)
	if err != nil {
		return nil, err
	}
	w, err := strconv.Atoi(weight)
	if err != nil {
		return nil, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	if p > math.MaxUint16 || p < 0 {
		return nil, errors.New("priority 区间不合法")
	}
	if w > math.MaxUint16 || w < 0 {
		return nil, errors.New("weight 区间不合法")
	}
	if portInt > math.MaxUint16 || w < 0 {
		return nil, errors.New("port 区间不合法")
	}
	return &RecordSRV{
		Ttl:      u,
		Priority: uint16(p),
		Weight:   uint16(w),
		Port:     uint16(portInt),
		Target:   target,
	}, nil
}

func (r *Record) IndexSRV(other *RecordSRV) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.SRV {
		if a.Priority == other.Priority &&
			a.Ttl == other.Ttl &&
			a.Weight == other.Weight &&
			a.Port == other.Port &&
			other.Target == a.Target {
			return i
		}
	}
	return scanner.EOF
}

func (r *Record) ModSRV(old, data *RecordSRV) error {
	if r.IndexSRV(data) != scanner.EOF {
		return errors.New("已存在相同的 SRV 记录")
	}
	if old != nil {
		oldIndex := r.IndexSRV(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 SRV 记录不存在")
		}
		r.SRV[oldIndex] = data
	} else {
		r.SRV = append(r.SRV, data)
	}
	return nil
}

func (r *Record) RemoveSRV(data *RecordSRV) error {
	index := r.IndexSRV(data)
	if index == scanner.EOF {
		return errors.New("当前 SRV 记录不存在")
	}
	r.SRV, _ = utils.RemoveIndex(r.SRV, index)
	return nil
}
