package gormkit

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

type BitBool struct {
	Bool  bool
	Valid bool
}

func NewBitBool(bool bool) BitBool {
	return BitBool{Bool: bool, Valid: true}
}

var (
	bitTrue  = []byte{0x1}
	bitFalse = []byte{0x0}
)

func (b *BitBool) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		if bytes.Compare(src.([]byte), bitTrue) == 0 {
			b.Bool, b.Valid = true, true
		}
		if bytes.Compare(src.([]byte), bitFalse) == 0 {
			b.Bool, b.Valid = false, true
		}
	}
	return nil
}
func (b *BitBool) Value() (driver.Value, error) {
	if b == nil {
		return nil, nil
	}
	if b.Valid {
		if b.Bool {
			return bitTrue, nil
		} else {
			return bitFalse, nil
		}
	}

	return bitFalse, nil
}
func (b BitBool) MarshalJSON() ([]byte, error) {
	if b.Valid {
		return json.Marshal(b.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (b *BitBool) UnmarshalJSON(data []byte) error {
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		b.Valid = true
		b.Bool = *x
	} else {
		b.Valid = false
	}
	return nil
}
