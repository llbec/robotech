package shard

type RouteResult struct {
	ProjectID string
	Month     string
	Path      string
}

type Router struct {
	basePath string
}

func NewRouter(basePath string) *Router {
	return &Router{basePath: basePath}
}

func (r *Router) Route(projectID string, blockTime int64) RouteResult {
	month := MonthKey(blockTime)
	return RouteResult{
		ProjectID: projectID,
		Month:     month,
		Path:      r.basePath + "/" + projectID + "/" + month,
	}
}
