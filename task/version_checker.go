package task

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	SemVer "github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/turbot/steampipe/constants"
	"github.com/turbot/steampipe/utils"
	"github.com/turbot/steampipe/version"
)

// the current version of the Steampipe CLI application
var currentVersion = version.SteampipeVersion.String()

type versionCheckResponse struct {
	NewVersion   string    `json:"latest_version,omitempty"` // `json:"current_version"`
	DownloadURL  string    `json:"download_url,omitempty"`   // `json:"download_url"`
	ChangelogURL string    `json:"html,omitempty"`           // `json:"changelog_url"`
	Alerts       []*string `json:"alerts,omitempty"`
}

// VersionChecker :: the version checker struct composition container.
// This MUST not be instantiated manually. Use `CreateVersionChecker` instead
type versionChecker struct {
	checkResult *versionCheckResponse // a channel to store the HTTP response
	signature   string                // flags whether update check should be done
}

// check if there is a new version
func checkSteampipeVersion(id string) []string {
	var notificationLines []string

	if !viper.GetBool(constants.ArgUpdateCheck) {
		return notificationLines
	}

	v := new(versionChecker)
	v.signature = id
	v.doCheckRequest()
	notificationLines, _ = v.notificationMessage()
	return notificationLines
}

// notificationMessage returns any required update notification as an array of strings
func (c *versionChecker) notificationMessage() ([]string, error) {
	info := c.checkResult
	if info == nil {
		return nil, nil
	}

	if info.NewVersion == "" {
		return nil, nil
	}

	newVersion, err := SemVer.NewVersion(info.NewVersion)
	if err != nil {
		return nil, err
	}
	currentVersion, err := SemVer.NewVersion(currentVersion)

	if err != nil {
		fmt.Println(fmt.Errorf("there's something wrong with the Current Version"))
		fmt.Println(err)
	}

	if newVersion.GreaterThan(currentVersion) {
		var downloadURLColor = color.New(color.FgYellow)
		var notificationLines = []string{
			"",
			fmt.Sprintf("A new version of Steampipe is available! %s → %s", constants.Bold(currentVersion), constants.Bold(newVersion)),
			fmt.Sprintf("You can update by downloading from %s", downloadURLColor.Sprint("https://steampipe.io/downloads")),
			"",
		}
		return notificationLines, nil
	}
	return nil, nil
}

func displayUpdateNotification(notificationLines []string) {
	// convert notificationLines into an array of arrays
	var notificationTable = make([][]string, len(notificationLines))
	for i, line := range notificationLines {
		notificationTable[i] = []string{line}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{})                // no headers please
	table.SetAlignment(tablewriter.ALIGN_LEFT) // we align to the left
	table.SetAutoWrapText(false)               // let's not wrap the text
	table.SetBorder(true)                      // there needs to be a border to give the dialog feel
	table.AppendBulk(notificationTable)        // Add Bulk Data

	fmt.Println()
	table.Render()
	fmt.Println()
}

// contact the Turbot Artifacts Server and retrieve the latest released version
func (c *versionChecker) doCheckRequest() {
	payload := utils.BuildRequestPayload(c.signature, map[string]interface{}{})
	sendRequestTo := c.versionCheckURL()
	timeout := 5 * time.Second
	resp, err := utils.SendRequest(c.signature, "POST", sendRequestTo, payload, timeout)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("[TRACE] Unknown response during version check: %d\n", resp.StatusCode)
		return
	}

	c.checkResult = c.decodeResult(bodyString)
}

func (c *versionChecker) decodeResult(body string) *versionCheckResponse {
	var result versionCheckResponse

	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil
	}
	return &result
}

func (c *versionChecker) versionCheckURL() url.URL {
	var u url.URL
	//https://hub.steampipe.io/api/cli/version/latest
	u.Scheme = "https"
	u.Host = "hub.steampipe.io"
	u.Path = "api/cli/version/latest"
	return u
}
