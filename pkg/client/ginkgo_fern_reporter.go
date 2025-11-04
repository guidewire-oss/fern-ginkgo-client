package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/guidewire-oss/fern-ginkgo-client/pkg/utils"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"

	gt "github.com/onsi/ginkgo/v2/types"
)

func getRunLevelTags() []models.Tag {
	envTags := os.Getenv("TEST_RUN_TAGS")
	if envTags == "" {
		return nil
	}

	parts := strings.Split(envTags, ",")
	tags := make([]models.Tag, 0, len(parts))
	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name != "" {
			tags = append(tags, models.Tag{Name: name})
		}
	}
	return tags
}

func (f *FernApiClient) Report(report gt.Report) error {

	var suiteRuns []models.SuiteRun
	suiteRun := models.SuiteRun{
		SuiteName: report.SuiteDescription,
		StartTime: report.StartTime,
		EndTime:   report.EndTime,
		Tags:      convertTags(report.SuiteLabels),
	}

	var specRuns []models.SpecRun
	for _, spec := range report.SpecReports {
		specRun := models.SpecRun{
			SpecDescription: spec.FullText(),
			Status:          spec.State.String(),
			Message:         spec.Failure.Message,
			StartTime:       spec.StartTime,
			EndTime:         spec.EndTime,
			Tags:            convertTags(spec.Labels()),
		}
		if spec.LeafNodeType != gt.NodeTypeIt {
			// It's a setup/teardown node
			fmt.Printf("Current node %s is not an It node. Skipping \n", spec.LeafNodeType)
		} else {
			specRuns = append(specRuns, specRun)
		}
	}

	suiteRun.SpecRuns = specRuns
	suiteRuns = append(suiteRuns, suiteRun)

	testRun := models.TestRun{
		TestProjectID: f.id, // Set this to your project id
		TestSeed:      uint64(report.SuiteConfig.RandomSeed),
		StartTime:     report.StartTime,
		EndTime:       time.Now(), // or report.EndTime if available
		Tags:          getRunLevelTags(),
		Environment:   os.Getenv("TEST_ENVIRONMENT"),
		SuiteRuns:     suiteRuns,
	}

	addMetadataInfo(&testRun)

	testJson, err := json.Marshal(testRun)
	if err != nil {
		panic(err)
	}

	bodyReader := bytes.NewReader(testJson)

	var reportUrl string
	if f.token != "" {
		reportUrl, err = url.JoinPath(f.baseURL, "api/v1/test-runs")
	} else {
		reportUrl, err = url.JoinPath(f.baseURL, "api/testrun")
	}
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, reportUrl, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// Add authorization token if available
	if f.token != "" {
		req.Header.Set("Authorization", "Bearer "+f.token)
	}
	resp, err := f.httpClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read the response body
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("client: error reading response body: %w", readErr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("client: Response status code: %d, Response: %s", resp.StatusCode, string(body))
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
	// Open the repository from the provided path - if not set to the root, go up the tree.
	repoPath, err = utils.FindGitRoot(repoPath)
	if err != nil {
		gitRootError := fmt.Errorf("⚠️ Warning: No Git repository found (%v)", err)
		fmt.Printf("%v\n", gitRootError)
		return "NA", "NA", gitRootError
	}
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
