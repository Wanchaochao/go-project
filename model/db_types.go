package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

type JsonStringSlice []string

func (v *JsonStringSlice) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	return json.Unmarshal(b, v) // receiver必须为指针
}

func (v JsonStringSlice) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'[', ']'}, nil
	}
	return json.Marshal(v) // receiver不能为指针
}

type JsonMapStringAny map[string]any

func (v *JsonMapStringAny) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v) // receiver必须为指针
}

func (v JsonMapStringAny) Value() (driver.Value, error) {
	if v == nil {
		return []byte{'{', '}'}, nil
	}
	return json.Marshal(v) // receiver不能为指针
}

type BoolString bool

func (v *BoolString) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	*v = string(b) == "Y"
	return nil
}

func (v BoolString) Value() (driver.Value, error) {
	s := "N"
	if v {
		s = "Y"
	}
	return s, nil // receiver不能为指针
}
