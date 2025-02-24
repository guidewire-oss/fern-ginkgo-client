package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"net/http"
	"net/url"
	"os"
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

	addMetadataInfo(&testRun)

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

func addMetadataInfo(testRun *models.TestRun) {
	if os.Getenv("GITHUB_ACTION") != "" {
		testRun.GitBranch = os.Getenv("GITHUB_REF_NAME")
		testRun.GitSha = os.Getenv("GITHUB_SHA")
		testRun.BuildTriggerActor = os.Getenv("GITHUB_TRIGGERING_ACTOR")
		testRun.BuildUrl = fmt.Sprintf("%s/%s/actions/runs/%s", os.Getenv("GITHUB_SERVER_URL"), os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_RUN_ID"))
	} else {
		repoPath := os.Getenv("GIT_REPO_PATH")
		if repoPath == "" {
			repoPath = "."
		}
		branch, commitSHA, _ := GetBranchAndCommit(repoPath)
		testRun.GitBranch = branch
		testRun.GitSha = commitSHA
	}
}

func GetBranchAndCommit(repoPath string) (branch string, commitSHA string, err error) {
	// Open the repository from the provided path.
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "NA", "NA", fmt.Errorf("failed to open repository: %w", err)
	}

	// Retrieve the HEAD reference.
	ref, err := repo.Head()
	if err != nil {
		return "NA", "NA", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commit object corresponding to the HEAD reference.
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "NA", "NA", fmt.Errorf("failed to get commit object: %w", err)
	}

	// Extract the short name of the branch from the reference.
	branch = ref.Name().Short()

	// Convert the commit hash to a string.
	commitSHA = commit.Hash.String()

	return branch, commitSHA, nil
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
