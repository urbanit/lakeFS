package multiparts

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/treeverse/lakefs/pkg/db"
)

type Metadata map[string]string

type MultipartUpload struct {
	UploadID        string            `db:"upload_id"`
	Path            string            `db:"path"`
	CreationDate    time.Time         `db:"creation_date"`
	PhysicalAddress string            `db:"physical_address"`
	Metadata        map[string]string `db:"metadata"`
}

type Tracker interface {
	Create(ctx context.Context, uploadID, path, physicalAddress string, creationTime time.Time, metadata Metadata) error
	Get(ctx context.Context, uploadID string) (*MultipartUpload, error)
	Delete(ctx context.Context, uploadID string) error
}

type tracker struct {
	db db.Database
}

var (
	ErrMultipartUploadNotFound  = fmt.Errorf("multipart upload not found")
	ErrInvalidUploadID          = errors.New("invalid upload id")
	ErrInvalidMetadataSrcFormat = errors.New("invalid metadata source format")
)

func (m Metadata) Set(k, v string) {
	m[strings.ToLower(k)] = v
}

func (m Metadata) Get(k string) string {
	return m[strings.ToLower(k)]
}
func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	data, ok := src.([]byte)
	if !ok {
		return ErrInvalidMetadataSrcFormat
	}
	return json.Unmarshal(data, m)
}

func NewTracker(adb db.Database) Tracker {
	return &tracker{
		db: adb,
	}
}

func (m *tracker) Create(ctx context.Context, uploadID, path, physicalAddress string, creationTime time.Time, metadata Metadata) error {
	if uploadID == "" {
		return ErrInvalidUploadID
	}
	_, err := m.db.Transact(ctx, func(tx db.Tx) (interface{}, error) {
		_, err := tx.Exec(`INSERT INTO gateway_multiparts (upload_id,path,creation_date,physical_address,metadata)
			VALUES ($1, $2, $3, $4, $5)`,
			uploadID, path, creationTime, physicalAddress, metadata)
		return nil, err
	})
	return err
}

func (m *tracker) Get(ctx context.Context, uploadID string) (*MultipartUpload, error) {
	if uploadID == "" {
		return nil, ErrInvalidUploadID
	}
	res, err := m.db.Transact(ctx, func(tx db.Tx) (interface{}, error) {
		var m MultipartUpload
		if err := tx.Get(&m, `
			SELECT upload_id, path, creation_date, physical_address, metadata 
			FROM gateway_multiparts
			WHERE upload_id = $1`,
			uploadID); err != nil {
			return nil, err
		}
		return &m, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*MultipartUpload), nil
}

func (m *tracker) Delete(ctx context.Context, uploadID string) error {
	if uploadID == "" {
		return ErrInvalidUploadID
	}
	_, err := m.db.Transact(ctx, func(tx db.Tx) (interface{}, error) {
		res, err := tx.Exec(`DELETE FROM gateway_multiparts WHERE upload_id = $1`, uploadID)
		if err != nil {
			return nil, err
		}
		affected := res.RowsAffected()
		if affected != 1 {
			return nil, ErrMultipartUploadNotFound
		}
		return nil, nil
	})
	return err
}
