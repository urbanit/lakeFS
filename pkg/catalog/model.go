package catalog

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/treeverse/lakefs/pkg/block"
)

const (
	DBEntryFieldChecksum        = "checksum"
	DBEntryFieldPhysicalAddress = "physical_address"

	ContentTypeKey     = "content-type"
	DefaultContentType = "application/octet-stream"
)

type Metadata map[string]string

type Repository struct {
	Name             string    `db:"name"`
	StorageNamespace string    `db:"storage_namespace"`
	DefaultBranch    string    `db:"default_branch"`
	CreationDate     time.Time `db:"creation_date"`
}

type DBEntry struct {
	CommonLevel     bool
	Path            string      `db:"path"`
	PhysicalAddress string      `db:"physical_address"`
	CreationDate    time.Time   `db:"creation_date"`
	Size            int64       `db:"size"`
	Checksum        string      `db:"checksum"`
	Metadata        Metadata    `db:"metadata"`
	Expired         bool        `db:"is_expired"`
	AddressType     AddressType `db:"address_type"`
}

type CommitLog struct {
	Reference    string
	Committer    string    `db:"committer"`
	Message      string    `db:"message"`
	CreationDate time.Time `db:"creation_date"`
	Metadata     Metadata  `db:"metadata"`
	MetaRangeID  string    `db:"meta_range_id"`
	Parents      []string
}

type MergeResult struct {
	Summary   map[DifferenceType]int
	Reference string
}

type Branch struct {
	Name      string `db:"name"`
	Reference string
}

type Tag struct {
	ID       string
	CommitID string
}

// AddressType is the type of entry address
type AddressType int32

const (
	// Deprecated: indicates that the address might be relative or full.
	// Used only for backward compatibility and should not be used for creating entries.
	AddressTypeByPrefixDeprecated AddressType = 0

	// AddressTypeRelative indicates that the address is relative to the storage namespace.
	// For example: "foo/bar"
	AddressTypeRelative AddressType = 1

	// AddressTypeFull indicates that the address is the full address of the object in the object store.
	// For example: "s3://bucket/foo/bar"
	AddressTypeFull AddressType = 2
)

// nolint:staticcheck
func (at AddressType) ToIdentifierType() block.IdentifierType {
	switch at {
	case AddressTypeByPrefixDeprecated:
		return block.IdentifierTypeUnknownDeprecated
	case AddressTypeRelative:
		return block.IdentifierTypeRelative
	case AddressTypeFull:
		return block.IdentifierTypeFull
	default:
		panic(fmt.Sprintf("unknown address type: %d", at))
	}
}

func (m Metadata) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(struct{}{})
	}
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

func (m Metadata) Set(key, value string) {
	k := strings.ToLower(key)
	// make sure content-type holds default value
	if k == ContentTypeKey && value == "" {
		value = DefaultContentType
	}
	m[k] = value
}

func (m Metadata) Get(key string) string {
	k := strings.ToLower(key)
	value := m[key]
	// make sure we always return default for content-type
	if k == ContentTypeKey && value == "" {
		return DefaultContentType
	}
	return value
}

func NewMetadataFromMap(m map[string]string) Metadata {
	meta := make(Metadata)
	for k, v := range m {
		meta[strings.ToLower(k)] = v
	}
	// set default if missing - default type
	if meta[ContentTypeKey] == "" {
		meta[ContentTypeKey] = DefaultContentType
	}
	return meta
}
