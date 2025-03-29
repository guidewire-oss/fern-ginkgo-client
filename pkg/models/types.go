package models

import (
	"time"
)

type TimeLog struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type TestRun struct {
	ID                uint64     `json:"id" gorm:"primaryKey"`
	TestProjectName   string     `json:"test_project_name"`
	TestSeed          uint64     `json:"test_seed"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           time.Time  `json:"end_time"`
	GitBranch         string     `json:"git_branch"`
	GitSha            string     `json:"git_sha"`
	BuildTriggerActor string     `json:"build_trigger_actor"`
	BuildUrl          string     `json:"build_url"`
	SuiteRuns         []SuiteRun `json:"suite_runs" gorm:"foreignKey:TestRunID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
  EnableGeminiInsights bool    `json:"enable_gemini_insights"`
}

type SuiteRun struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	TestRunID uint64    `json:"test_run_id"`
	SuiteName string    `json:"suite_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SpecRuns  []SpecRun `json:"spec_runs" gorm:"foreignKey:SuiteID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SpecRun struct {
	ID              uint64    `json:"id" gorm:"primaryKey"`
	SuiteID         uint64    `json:"suite_id"`
	SpecDescription string    `json:"spec_description"`
	Status          string    `json:"status"`
	Message         string    `json:"message"`
	Tags            []Tag     `json:"tags" gorm:"many2many:spec_run_tags;"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
}

type Tag struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
