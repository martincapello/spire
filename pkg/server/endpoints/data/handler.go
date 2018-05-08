package data

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spiffe/spire/pkg/server/catalog"
	"github.com/spiffe/spire/proto/api/data"
)

type Handler struct {
	Log     logrus.FieldLogger
	Catalog catalog.Catalog
}

func (h *Handler) Dump(req *data.Empty, stream data.Data_DumpServer) error {
	//ds := h.Catalog.DataStores()[0]

	return nil
}

func (h *Handler) Replay(stream data.Data_ReplayServer) error {
	return fmt.Errorf("not implemented")
}
