package geo

import (
	_ "embed"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

var (
	//go:embed "countries.mmdb"
	countries []byte
)

type DB struct {
	reader *maxminddb.Reader
}

func New() (*DB, error) {
	reader, err := maxminddb.FromBytes(countries)
	if err != nil {
		return nil, err
	}

	return &DB{reader: reader}, nil
}

func (db *DB) Close() error {
	return db.reader.Close()
}

func (db *DB) LookupCountryCode(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ""
	}

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err := db.reader.Lookup(parsed, &record)
	if err != nil {
		return ""
	}

	return record.Country.ISOCode
}
