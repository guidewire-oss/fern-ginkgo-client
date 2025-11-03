package models

import (
	"time"
)

type TimeLog struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type TestRun struct {
	ID                uint64     `json:"id"`
	TestProjectName   string     `json:"test_project_name"`
	TestProjectID     string     `json:"test_project_id"`
	TestSeed          uint64     `json:"test_seed"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           time.Time  `json:"end_time"`
	GitBranch         string     `json:"git_branch"`
	GitSha            string     `json:"git_sha"`
	BuildTriggerActor string     `json:"build_trigger_actor"`
	BuildUrl          string     `json:"build_url"`
	Tags              []Tag      `json:"tags"`
	Environment       string     `json:"environment"`
	SuiteRuns         []SuiteRun `json:"suite_runs"`
}

type SuiteRun struct {
	ID        uint64    `json:"id"`
	TestRunID uint64    `json:"test_run_id"`
	SuiteName string    `json:"suite_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Tags      []Tag     `json:"tags"`
	SpecRuns  []SpecRun `json:"spec_runs"`
}

type SpecRun struct {
	ID              uint64    `json:"id"`
	SuiteID         uint64    `json:"suite_id"`
	SpecDescription string    `json:"spec_description"`
	Status          string    `json:"status"`
	Message         string    `json:"message"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Tags            []Tag     `json:"tags"`
}

type Tag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
