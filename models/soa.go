package models

import (
	"fmt"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordSOA struct {
	Ttl     uint32 `json:"ttl"`
	MName   string `json:"mname"`
	RName   string `json:"rname"`
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minimum uint32 `json:"minimum"`
}

func (r *RecordSOA) Clone() *RecordSOA {
	return &RecordSOA{
		Ttl:     r.Ttl,
		MName:   r.MName,
		RName:   r.RName,
		Serial:  r.Serial,
		Refresh: r.Refresh,
		Retry:   r.Retry,
		Expire:  r.Expire,
		Minimum: r.Minimum,
	}
}

func FromRecordSOA(values GetValue) (*RecordSOA, error) {
	return NewRecordSOA(
		values("mname"),
		values("rname"),
		values("serial"),
		values("refresh"),
		values("retry"),
		values("expire"),
		values("minimum"),
		values("ttl"))
}
func NewRecordSOA(mName, rName, serial, refresh, retry, expire, minimum, ttl string) (*RecordSOA, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	if err = utils.RegexHost.Valid(mName); err != nil {
		return nil, err
	}
	if err = utils.RegexHost.Valid(rName); err != nil {
		return nil, err
	}
	serialInt, err := utils.AtoUint32(serial)
	if err != nil || serialInt > 65535 || serialInt < 0 {
		return nil, err
	}
	refreshInt, err := utils.AtoUint32(refresh)
	if err != nil {
		return nil, err
	}
	// 小于 refresh
	retryInt, err := utils.AtoUint32(retry)
	if err != nil {
		return nil, err
	}
	//大于 Refresh 和 Retry 的总和
	expireInt, err := utils.AtoUint32(expire)
	if err != nil {
		return nil, err
	}
	minimumInt, err := utils.AtoUint32(minimum)
	if err != nil {
		return nil, err
	}
	if retryInt >= refreshInt {
		return nil, fmt.Errorf("retry > refresh")
	}

	if expireInt <= refreshInt+retryInt {
		return nil, fmt.Errorf("expire < retry + refresh")
	}
	return &RecordSOA{
		Ttl:     u,
		MName:   mName,
		RName:   rName,
		Serial:  serialInt,
		Refresh: refreshInt,
		Retry:   retryInt,
		Expire:  expireInt,
		Minimum: minimumInt,
	}, nil
}

func (r *Record) ModSOA(old, data *RecordSOA) error {
	if old == nil {
		if r.SOA != nil {
			return errors.New("已存在相同的 SOA 记录")
		} else {
			r.SOA = data
		}
	} else {
		if r.SOA == nil {
			return errors.New("SOA 记录不存在")
		} else {
			r.SOA = data
		}
	}
	return nil
}

func (r *Record) RemoveSOA(data *RecordSOA) error {
	if r.SOA == nil {
		return errors.New("当前 SOA 记录不存在")
	}
	r.SOA = nil
	return nil
}
