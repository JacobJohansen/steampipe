package controlexecute

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/turbot/go-kit/helpers"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe/constants"
	"github.com/turbot/steampipe/dashboard/dashboardtypes"
	"github.com/turbot/steampipe/query/queryresult"
	"github.com/turbot/steampipe/steampipeconfig/modconfig"
	"github.com/turbot/steampipe/utils"
)

type ResultRows []*ResultRow

// ToLeafData converts the result rows to snapshot data format
func (r ResultRows) ToLeafData(dimensionSchema map[string]*dashboardtypes.ColumnSchema) *dashboardtypes.LeafData {
	var res = &dashboardtypes.LeafData{
		Columns: []*dashboardtypes.ColumnSchema{
			{Name: "reason", DataType: "TEXT"},
			{Name: "resource", DataType: "TEXT"},
			{Name: "status", DataType: "TEXT"},
		},
		Rows: make([]map[string]interface{}, len(r)),
	}
	for _, d := range dimensionSchema {
		res.Columns = append(res.Columns, d)
	}
	for i, row := range r {
		res.Rows[i] = map[string]interface{}{
			"reason":   row.Reason,
			"resource": row.Resource,
			"status":   row.Status,
		}
		// flatten dimensions
		for _, d := range row.Dimensions {
			res.Rows[i][d.Key] = d.Value
		}
	}
	return res
}

// ResultRow is the result of a control execution for a single resource
type ResultRow struct {
	// reason for the status
	Reason string `json:"reason" csv:"reason"`
	// resource name
	Resource string `json:"resource" csv:"resource"`
	// status of the row (ok, info, alarm, error, skip)
	Status string `json:"status" csv:"status"`
	// dimensions for this row
	Dimensions []Dimension `json:"dimensions"`
	// parent control run
	Run *ControlRun `json:"-"`
	// source control
	Control *modconfig.Control `json:"-" csv:"control_id:UnqualifiedName,control_title:Title,control_description:Description"`
}

// GetDimensionValue returns the value for a dimension key. Returns an empty string with 'false' if not found
func (r *ResultRow) GetDimensionValue(key string) string {
	for _, dim := range r.Dimensions {
		if dim.Key == key {
			return dim.Value
		}
	}
	return ""
}

// AddDimension checks whether a column value is a scalar type, and if so adds it to the Dimensions map
func (r *ResultRow) AddDimension(c *sql.ColumnType, val interface{}) {
	switch c.ScanType().Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		return
	default:
		r.Dimensions = append(r.Dimensions, Dimension{Key: c.Name(), Value: typehelpers.ToString(val), SqlType: c.DatabaseTypeName()})
	}
}

func NewResultRow(run *ControlRun, row *queryresult.RowResult, colTypes []*sql.ColumnType) (*ResultRow, error) {
	// validate the required columns exist in the result
	if err := validateColumns(colTypes); err != nil {
		return nil, err
	}
	res := &ResultRow{
		Run:     run,
		Control: run.Control,
	}

	// was there a SQL error _executing the control
	// Note: this is different from the control state being 'error'
	if row.Error != nil {
		return nil, row.Error
	}

	for i, c := range colTypes {
		switch c.Name() {
		case "reason":
			res.Reason = typehelpers.ToString(row.Data[i])
		case "resource":
			res.Resource = typehelpers.ToString(row.Data[i])
		case "status":
			status := typehelpers.ToString(row.Data[i])
			if !IsValidControlStatus(status) {
				return nil, fmt.Errorf("invalid control status '%s'", status)
			}
			res.Status = status
		default:
			// if this is a scalar type, add to dimensions
			res.AddDimension(c, row.Data[i])
		}
	}
	return res, nil
}

func IsValidControlStatus(status string) bool {
	return helpers.StringSliceContains([]string{constants.ControlOk, constants.ControlAlarm, constants.ControlInfo, constants.ControlError, constants.ControlSkip}, status)
}

func validateColumns(colTypes []*sql.ColumnType) error {
	requiredColumns := []string{"reason", "resource", "status"}
	var missingColumns []string
	for _, col := range requiredColumns {
		if !columnTypesContainsColumn(col, colTypes) {
			missingColumns = append(missingColumns, col)
		}
	}
	if len(missingColumns) > 0 {
		return fmt.Errorf("control result is missing required %s: %v", utils.Pluralize("column", len(missingColumns)), missingColumns)
	}
	return nil
}

func columnTypesContainsColumn(col string, colTypes []*sql.ColumnType) bool {
	for _, ct := range colTypes {
		if ct.Name() == col {
			return true
		}
	}
	return false
}
