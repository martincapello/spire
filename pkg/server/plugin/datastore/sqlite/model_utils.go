package sqlite

import (
	"bytes"
)

func (b *Bundle) Append(cert CACert) {
	b.CACerts = append(b.CACerts, cert)
}

func (b *Bundle) Contains(cert CACert) bool {
	for _, c := range b.CACerts {
		if bytes.Equal(c.Cert, cert.Cert) {
			return true
		}
	}

	return false
}
