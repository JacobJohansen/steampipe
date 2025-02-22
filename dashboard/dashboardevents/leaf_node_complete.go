package dashboardevents

import "github.com/turbot/steampipe/dashboard/dashboardtypes"

type LeafNodeComplete struct {
	LeafNode    dashboardtypes.DashboardNodeRun
	Session     string
	ExecutionId string
}

// IsDashboardEvent implements DashboardEvent interface
func (*LeafNodeComplete) IsDashboardEvent() {}
