package controldisplay

// import (
// 	"bytes"
// 	"context"
// 	"io"
// 	"strings"
// 	"testing"

// 	jsonpatch "github.com/evanphx/json-patch"
// 	"github.com/spf13/viper"
// 	"github.com/turbot/steampipe/constants"
// 	"github.com/turbot/steampipe/control/controlexecute"
// 	"github.com/turbot/steampipe/steampipeconfig/modconfig"
// )

// var rootBenchmark = modconfig.Benchmark{}
// var childBenchmark1 = modconfig.Benchmark{}
// var childBenchmark2 = modconfig.Benchmark{}

// var desc = "Dummy control for unit testing"
// var title = "DummyControl"
// var c11 = modconfig.Control{
// 	Title:       &title,
// 	Description: &desc,
// }
// var c12 = modconfig.Control{
// 	Title:       &title,
// 	Description: &desc,
// }
// var c21 = modconfig.Control{
// 	Title:       &title,
// 	Description: &desc,
// }
// var c22 = modconfig.Control{
// 	Title:       &title,
// 	Description: &desc,
// }

// var tree = &controlexecute.ExecutionTree{
// 	Root: &controlexecute.ResultGroup{
// 		GroupId:     "DummyTest",
// 		Parent:      nil,
// 		Title:       "Test Root Group",
// 		Description: "Description for test root group",
// 		Summary: &controlexecute.GroupSummary{
// 			Status: controlexecute.StatusSummary{
// 				Alarm: 2,
// 			},
// 		},
// 		GroupItem: &rootBenchmark,
// 		Groups: []*controlexecute.ResultGroup{
// 			{
// 				GroupItem: &childBenchmark1,
// 				ControlRuns: []*controlexecute.ControlRun{
// 					{
// 						Control: &c11,
// 						Rows: []*controlexecute.ResultRow{
// 							{
// 								Status:     constants.ControlAlarm,
// 								Reason:     "is pretty insecure",
// 								Resource:   "some other resource",
// 								Dimensions: []controlexecute.Dimension{},
// 								Run:        &controlexecute.ControlRun{Control: &c11},
// 							},
// 						},
// 					},
// 					{
// 						Control: &c12,
// 						Rows: []*controlexecute.ResultRow{
// 							{
// 								Status:     constants.ControlAlarm,
// 								Reason:     "is pretty insecure",
// 								Resource:   "some other resource",
// 								Dimensions: []controlexecute.Dimension{},
// 								Run:        &controlexecute.ControlRun{Control: &c12},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			{
// 				GroupItem: &childBenchmark2,
// 				ControlRuns: []*controlexecute.ControlRun{
// 					{
// 						Control: &c21,
// 						Rows: []*controlexecute.ResultRow{
// 							{
// 								Status:     constants.ControlAlarm,
// 								Reason:     "is pretty insecure",
// 								Resource:   "some other resource",
// 								Dimensions: []controlexecute.Dimension{},
// 								Run:        &controlexecute.ControlRun{Control: &c21},
// 							},
// 						},
// 					},
// 					{
// 						Control: &c22,
// 						Rows: []*controlexecute.ResultRow{
// 							{
// 								Status:     constants.ControlAlarm,
// 								Reason:     "is pretty insecure",
// 								Resource:   "some other resource",
// 								Dimensions: []controlexecute.Dimension{},
// 								Run:        &controlexecute.ControlRun{Control: &c22},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

// const expectedJsonOutput = `{
// 	"group_id": "DummyTest",
// 	"title": "Test Root Group",
// 	"description": "Description for test root group",
// 	"tags": null,
// 	"summary": {
// 		"status": {
// 			"alarm": 2,
// 			"ok": 0,
// 			"info": 0,
// 			"skip": 0,
// 			"error": 0
// 		}
// 	},
// 	"groups": [
// 		{
// 			"group_id": "",
// 			"title": "",
// 			"description": "",
// 			"tags": null,
// 			"summary": {
// 				"status": {
// 					"alarm": 0,
// 					"ok": 0,
// 					"info": 0,
// 					"skip": 0,
// 					"error": 0
// 				}
// 			},
// 			"groups": null,
// 			"controls": [
// 				{
// 					"control_id": "",
// 					"description": "",
// 					"severity": "",
// 					"tags": null,
// 					"title": "",
// 					"results": [
// 						{
// 							"reason": "is pretty insecure",
// 							"resource": "some other resource",
// 							"status": "alarm",
// 							"dimensions": []
// 						}
// 					]
// 				},
// 				{
// 					"control_id": "",
// 					"description": "",
// 					"severity": "",
// 					"tags": null,
// 					"title": "",
// 					"results": [
// 						{
// 							"reason": "is pretty insecure",
// 							"resource": "some other resource",
// 							"status": "alarm",
// 							"dimensions": []
// 						}
// 					]
// 				}
// 			]
// 		},
// 		{
// 			"group_id": "",
// 			"title": "",
// 			"description": "",
// 			"tags": null,
// 			"summary": {
// 				"status": {
// 					"alarm": 0,
// 					"ok": 0,
// 					"info": 0,
// 					"skip": 0,
// 					"error": 0
// 				}
// 			},
// 			"groups": null,
// 			"controls": [
// 				{
// 					"control_id": "",
// 					"description": "",
// 					"severity": "",
// 					"tags": null,
// 					"title": "",
// 					"results": [
// 						{
// 							"reason": "is pretty insecure",
// 							"resource": "some other resource",
// 							"status": "alarm",
// 							"dimensions": []
// 						}
// 					]
// 				},
// 				{
// 					"control_id": "",
// 					"description": "",
// 					"severity": "",
// 					"tags": null,
// 					"title": "",
// 					"results": [
// 						{
// 							"reason": "is pretty insecure",
// 							"resource": "some other resource",
// 							"status": "alarm",
// 							"dimensions": []
// 						}
// 					]
// 				}
// 			]
// 		}
// 	],
// 	"controls": null
// }`

// func TestJsonFormatter(t *testing.T) {
// 	f := new(JSONFormatter)
// 	reader, _ := f.Format(context.Background(), tree)
// 	b := bytes.NewBufferString("")
// 	_, _ = io.Copy(b, reader)
// 	output := b.String()
// 	if !jsonpatch.Equal([]byte(expectedJsonOutput), []byte(output)) {
// 		t.Log(`"expected" is not equal to "output"`)
// 		t.FailNow()
// 	}
// }

// const expectedCsvOutput = `group_id,title,description,control_id,control_title,control_description,reason,resource,status
// ,,,,DummyControl,Dummy control for unit testing,is pretty insecure,some other resource,alarm
// ,,,,DummyControl,Dummy control for unit testing,is pretty insecure,some other resource,alarm
// ,,,,DummyControl,Dummy control for unit testing,is pretty insecure,some other resource,alarm
// ,,,,DummyControl,Dummy control for unit testing,is pretty insecure,some other resource,alarm`

// func TestCsvFormatter(t *testing.T) {
// 	tree.DimensionColorGenerator, _ = controlexecute.NewDimensionColorGenerator(4, 27)
// 	viper.Set(constants.ArgSeparator, ",")
// 	viper.Set(constants.ArgHeader, true)
// 	f := new(CSVFormatter)
// 	reader, _ := f.Format(context.Background(), tree)
// 	b := bytes.NewBufferString("")
// 	_, _ = io.Copy(b, reader)
// 	output := b.String()
// 	spacer := strings.TrimSpace(output)
// 	if spacer != expectedCsvOutput {
// 		t.Log(`"expected" is not equal to "output"`)
// 		t.Logf(spacer)
// 		t.FailNow()
// 	}
// }
