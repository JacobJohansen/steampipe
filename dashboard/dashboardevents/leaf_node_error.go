package dashboardevents

import "github.com/turbot/steampipe/dashboard/dashboardtypes"

type LeafNodeError struct {
	LeafNode    dashboardtypes.DashboardNodeRun
	Session     string
	Error       error
	ExecutionId string
}

// IsDashboardEvent implements DashboardEvent interface
func (*LeafNodeError) IsDashboardEvent() {}
