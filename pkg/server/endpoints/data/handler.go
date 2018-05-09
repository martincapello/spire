package data

import (
	"github.com/sirupsen/logrus"
	"github.com/spiffe/spire/pkg/server/catalog"
	"github.com/spiffe/spire/proto/api/data"
	"github.com/spiffe/spire/proto/server/datastore"
)

type Handler struct {
	Log     logrus.FieldLogger
	Catalog catalog.Catalog
}

func (h *Handler) Dump(req *data.Empty, stream data.Data_DumpServer) error {
	ds := h.Catalog.DataStores()[0]

	datastore.Da
	return ds.Dump(req, stream)
}

func (h *Handler) Restore(stream data.Data_RestoreServer) error {
	ds := h.Catalog.DataStores()[0]
	return ds.Restore(req, stream)
}
