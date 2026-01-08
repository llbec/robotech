package httpapi

import (
	"encoding/json"
	"ether-indexer/internal/model"
	"ether-indexer/internal/storage"
	"ether-indexer/internal/syncer"
	"net/http"
)

type ProjectCreateReq struct {
	ProjectID   string `json:"project_id"`
	RPCEndpoint string `json:"rpc_endpoint"`
	Description string `json:"description"`
	BlockRange  uint64 `json:"block_range"`
}

type ProjectActiveReq struct {
	ProjectID     string `json:"project_id"`
	Active        bool   `json:"active"`
	StartBlock    uint64 `json:"start_block"`
	Step          uint64 `json:"step"`
	IntervalSec   int    `json:"interval_sec"`
	AddressesJSON string `json:"addresses_json"`
	TopicsJSON    string `json:"topics_json"`
}

func createProject(projRepo *storage.ProjectRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p ProjectCreateReq
		_ = json.NewDecoder(r.Body).Decode(&p)

		err := projRepo.Create(
			p.ProjectID,
			p.RPCEndpoint,
			p.Description,
			p.BlockRange,
		)

		if err != nil {
			http.Error(w, err.Error(), 50001)
			return
		}

		if err := syncer.RegisterScheduler(&model.Project{
			ProjectID:   p.ProjectID,
			Active:      false,
			RPCEndpoint: p.RPCEndpoint,
			Description: p.Description,
			BlockRange:  p.BlockRange,
		}); err != nil {
			http.Error(w, err.Error(), 50002)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
