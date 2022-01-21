package gmutex

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
)

// LockJSON calls LockData with the JSON encoding of v.
func (m *Mutex) LockJSON(ctx context.Context, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.LockData(ctx, bytes.NewReader(b))
}

// TryLockJSON calls TryLockData with the JSON encoding of v.
// Parses JSON-encoded data into the value pointed to by v,
// if the lock is already in use and v is a pointer.
func (m *Mutex) TryLockJSON(ctx context.Context, v interface{}) (bool, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return false, err
	}

	if rv := reflect.ValueOf(v); rv.Kind() != reflect.Ptr || rv.IsNil() {
		return m.TryLockData(ctx, bytes.NewReader(b))
	}

	buf := bytes.NewBuffer(b)
	locked, err := m.TryLockData(ctx, buf)
	if locked || err != nil {
		return locked, err
	}
	return false, json.Unmarshal(buf.Bytes(), v)
}

// UpdateJSON calls Update with the JSON encoding of v.
func (m *Mutex) UpdateJSON(ctx context.Context, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Update(ctx, bytes.NewReader(b))
}

// AdoptJSON calls Adopt with the JSON encoding of v.
func (m *Mutex) AdoptJSON(ctx context.Context, id string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return m.Adopt(ctx, id, bytes.NewReader(b))
}

// InspectJSON calls Inspect parsing JSON-encoded data
// into the value pointed to by v.
func (m *Mutex) InspectJSON(ctx context.Context, v interface{}) (bool, error) {
	var buf bytes.Buffer
	locked, err := m.Inspect(ctx, &buf)
	if err == nil {
		err = json.Unmarshal(buf.Bytes(), v)
	}
	return locked, err
}
