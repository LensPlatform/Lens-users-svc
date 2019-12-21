package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonEmbeddable struct{}

// ser struct implements the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (u JsonEmbeddable) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// User struct implements the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (u JsonEmbeddable) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}
