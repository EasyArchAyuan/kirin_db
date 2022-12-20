package metadata_management

import (
	rm "record_manager"
	"tx"
)

type MetadataInterface interface {
	CreateTable(tblName string, sch *rm.Schema, tx *tx.Transaction)
	GetLayout(tblName string, tx *tx.Transaction) *rm.Layout
}
