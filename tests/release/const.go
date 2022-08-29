package release

import (
	"time"
)

const (
	environment              string = "environment"
	snapshotName             string = "snapshot"
	releasePlanName          string = "release-plan"
	releasePlanAdmissionName string = "release-plan-admission"
	releaseStrategyName      string = "strategy"
	releaseName              string = "release"
	releasePipelineName      string = "release-pipeline"
	applicationName          string = "application"
	releasePipelineBundle    string = "quay.io/hacbs-release/demo:m5-alpine"
	releaseStrategyPolicy    string = "policy"

	avgPipelineCompletionTime = 10 * time.Minute
	defaultInterval           = 100 * time.Millisecond
)
