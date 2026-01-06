package model

type Project struct {
	ProjectID   string
	Active      bool
	RPCEndpoint string
	Description string
	BlockRange  uint64
	BasePath    string
	CreatedAt   int64
	UpdatedAt   int64
}
