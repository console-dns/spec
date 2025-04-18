package models

import (
	"strconv"
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordCAA struct {
	Flag  uint8  `json:"flag"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (record *RecordCAA) Clone() *RecordCAA {
	return &RecordCAA{
		Flag:  record.Flag,
		Tag:   record.Tag,
		Value: record.Value,
	}
}

func FromRecordCAA(values GetValue) (*RecordCAA, error) {
	return NewRecordCAA("0", values("tag"), values("value"))
}
func NewRecordCAA(flag, tag, value string) (*RecordCAA, error) {
	if flag != "0" && flag != "128" {
		return nil, errors.New("未知 flag 类型")
	}
	if len(value) == 0 {
		return nil, errors.New("value 不能为空")
	}
	switch tag {
	case "issue", "issuewild", "iodef", "contactphone":
	case "contactemail":
		err := utils.RegexMail.Valid(value)
		if err != nil {
			return nil, errors.New("value 不是合法的邮件")
		}
	default:
		return nil, errors.New("未知 tag 类型")
	}
	f, _ := strconv.Atoi(flag)
	return &RecordCAA{
		Flag:  uint8(f),
		Tag:   tag,
		Value: value,
	}, nil
}

func (r *Record) IndexCAA(other *RecordCAA) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.CAA {
		if a.Value == other.Value && a.Tag == other.Tag && other.Flag == a.Flag {
			return i
		}
	}
	return scanner.EOF
}

func (r *Record) ModCAA(old, data *RecordCAA) error {
	if r.IndexCAA(data) != -1 {
		return errors.New("已存在相同的 CAA 记录")
	}
	if old != nil {
		oldIndex := r.IndexCAA(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 CAA 记录不存在")
		}
		r.CAA[oldIndex] = data
	} else {
		r.CAA = append(r.CAA, data)
	}
	return nil
}

func (r *Record) RemoveCAA(data *RecordCAA) error {
	index := r.IndexCAA(data)
	if index == scanner.EOF {
		return errors.New("当前 CAA 记录不存在")
	}
	r.CAA, _ = utils.RemoveIndex(r.CAA, index)
	return nil
}
