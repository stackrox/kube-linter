package jobttlsecondsafterfinished

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/jobttlsecondsafterfinished/internal/params"
	batchV1 "k8s.io/api/batch/v1"
)

const (
	JobKind        = "Job"
	CronJobKind    = "CronJob"
	job_no_ttl     = "job_no_ttl"
	job_ttl        = "job_ttl"
	cronjob_no_ttl = "cronjob_no_ttl"
	cronjob_ttl    = "cronjob_ttl"
)


type JobTTLSecondsAfterFinishedTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *JobTTLSecondsAfterFinishedTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *JobTTLSecondsAfterFinishedTestSuite) AddJobLike(kind string, name string, ttl *int32) {
	switch kind {
	case JobKind:
		s.ctx.AddMockJob(s.T(), name)
		if ttl != nil {
			s.ctx.ModifyJob(s.T(), name, func(job *batchV1.Job) {
				job.Spec.TTLSecondsAfterFinished = ttl
			})
		}
	case CronJobKind:
		s.ctx.AddMockCronJob(s.T(), name)
		if ttl != nil {
			s.ctx.ModifyCronJob(s.T(), name, func(cronjob *batchV1.CronJob) {
				cronjob.Spec.JobTemplate.Spec.TTLSecondsAfterFinished = ttl
			})
		}
	}
}

func TestJobTTLSecondsAfterFinished(t *testing.T) {
	suite.Run(t, new(JobTTLSecondsAfterFinishedTestSuite))
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestJobTTL() {
	ttl := int32(100)
	s.AddJobLike(JobKind, job_ttl, &ttl)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				job_ttl: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestJobNoTTL() {
	s.AddJobLike(JobKind, job_no_ttl, nil)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				job_no_ttl: {{Message: "Standalone Job does not specify ttlSecondsAfterFinished"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestCronJobTTL() {
	ttl := int32(100)
	s.AddJobLike(CronJobKind, cronjob_ttl, &ttl)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				cronjob_ttl: {{Message: "Managed Job specifies ttlSecondsAfterFinished which might conflict with successfulJobsHistoryLimit and failedJobsHistoryLimit from CronJob. Final behaviour is determined by the stricktier"}}, // must be nil
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestCronJobNoTTL() {
	s.AddJobLike(CronJobKind, cronjob_no_ttl, nil)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				cronjob_no_ttl: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
