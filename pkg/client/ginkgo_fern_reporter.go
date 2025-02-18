package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"

	gt "github.com/onsi/ginkgo/v2/types"
)

func (f *FernApiClient) Report(testName string, report gt.Report) error {

	var suiteRuns []models.SuiteRun
	suiteRun := models.SuiteRun{
		SuiteName: report.SuiteDescription,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
	}

	var specRuns []models.SpecRun
	for _, spec := range report.SpecReports {
		specRun := models.SpecRun{
			SpecDescription: spec.FullText(),
			Status:          spec.State.String(),
			Message:         spec.Failure.Message,
			StartTime:       spec.StartTime,
			EndTime:         spec.EndTime,
		}

		// Accessing the suite labels
		labels := report.SuiteLabels
		// Add gitSHA to tags
		if gitSha := getGitSha(); gitSha != "" {
			labels = append(labels, getGitSha())
		}
		// logic to convert suite labels string to []Tag
		specRun.Tags = convertTags(labels)
		specRuns = append(specRuns, specRun)
	}

	suiteRun.SpecRuns = specRuns
	suiteRuns = append(suiteRuns, suiteRun)

	testRun := models.TestRun{
		TestProjectName: f.name, // Set this to your project name
		TestSeed:        uint64(report.SuiteConfig.RandomSeed),
		StartTime:       report.StartTime,
		EndTime:         time.Now(), // or report.EndTime if available
		SuiteRuns:       suiteRuns,
	}

	testJson, err := json.Marshal(testRun)
	if err != nil {
		panic(err)
	}

	bodyReader := bytes.NewReader(testJson)

	url, err := url.JoinPath(f.baseURL, "api/testrun")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	_, err = f.httpClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}

	return nil
}

func convertTags(specLabels []string) []models.Tag {
	// Convert Ginkgo tags to Tag struct
	var tags []models.Tag
	for _, label := range specLabels {
		tags = append(tags, models.Tag{
			Name: label, // Or however you want to define the tag
		})
	}
	return tags
}

func getGitSha() string {
	// Check for uncommitted changes
	cmd := exec.Command("git", "diff", "--quiet")
	err := cmd.Run()
	if err != nil {
		return "gitSHA:NA"
	}

	// Check for unpushed commits
	cmd = exec.Command("git", "rev-list", "--count", "HEAD..origin/$(git rev-parse --abbrev-ref HEAD)")
	unpushed, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(unpushed)) != "0" {
		return "gitSHA:NA"
	}

	// get GIT SHA
	cmd = exec.Command("git", "rev-parse", "HEAD")
	sha, err := cmd.Output()
	if err != nil {
		fmt.Printf("client error: could not find the Git SHA of last commit")
		return ""
	}
	return "gitSHA:" + strings.TrimSpace(string(sha))
}
