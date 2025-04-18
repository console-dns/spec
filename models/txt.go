package models

import (
	"text/scanner"

	"github.com/console-dns/spec/utils"
	"github.com/pkg/errors"
)

type RecordTXT struct {
	Ttl  uint32 `json:"ttl"`
	Text string `json:"text"`
}

func (r *RecordTXT) Clone() *RecordTXT {
	return &RecordTXT{
		Ttl:  r.Ttl,
		Text: r.Text,
	}
}

func FromRecordTXT(values GetValue) (*RecordTXT, error) {
	return NewRecordTXT(values("text"), values("ttl"))
}
func NewRecordTXT(txt, ttl string) (*RecordTXT, error) {
	u, err := utils.ParseTtl(ttl)
	if err != nil {
		return nil, err
	}
	return &RecordTXT{
		Ttl:  u,
		Text: txt,
	}, nil
}
func (r *Record) IndexTXT(other *RecordTXT) int {
	if other == nil {
		return scanner.EOF
	}
	for i, a := range r.TXT {
		if a.Text == other.Text && a.Ttl == other.Ttl {
			return i
		}
	}
	return scanner.EOF
}
func (r *Record) ModTXT(old, data *RecordTXT) error {
	if r.IndexTXT(data) != scanner.EOF {
		return errors.New("已存在相同的 TXT 记录")
	}
	if old != nil {
		oldIndex := r.IndexTXT(old)
		if oldIndex == scanner.EOF {
			return errors.New("当前 TXT 记录不存在")
		}
		r.TXT[oldIndex] = data
	} else {
		r.TXT = append(r.TXT, data)
	}
	return nil
}

func (r *Record) RemoveTXT(data *RecordTXT) error {
	index := r.IndexTXT(data)
	if index == scanner.EOF {
		return errors.New("当前 TXT 记录不存在")
	}
	r.TXT, _ = utils.RemoveIndex(r.TXT, index)
	return nil
}
