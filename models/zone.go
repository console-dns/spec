package models

type Zone struct {
	Records map[string]*Record `json:"records" yaml:"records" toml:"records"`
}

func NewZone() *Zone {
	return &Zone{
		Records: make(map[string]*Record),
	}
}

func (r *Zone) Count() int {
	count := 0
	for _, z := range r.Records {
		count = count + z.Count()
	}
	return count
}

func (r *Zone) ModRecord(name string, f func(r *Record) error) error {
	if r.Records[name] == nil {
		r.Records[name] = NewRecord()
	}
	err := f(r.Records[name])
	if err != nil {
		return err
	}
	return nil
}
