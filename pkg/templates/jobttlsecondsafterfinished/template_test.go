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
	JobKind      = "Job"
	CronJobKind  = "CronJob"
	jobNoTTL     = "job_no_ttl"
	jobTTL       = "job_ttl"
	cronjobNoTTL = "cronjob_no_ttl"
	cronjobTTL   = "cronjob_ttl"
)

type JobTTLSecondsAfterFinishedTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *JobTTLSecondsAfterFinishedTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *JobTTLSecondsAfterFinishedTestSuite) AddJobLike(kind, name string, ttl *int32) {
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
	s.AddJobLike(JobKind, jobTTL, &ttl)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				jobTTL: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestJobNoTTL() {
	s.AddJobLike(JobKind, jobNoTTL, nil)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				jobNoTTL: {{Message: "Standalone Job does not specify ttlSecondsAfterFinished"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestCronJobTTL() {
	ttl := int32(100)
	s.AddJobLike(CronJobKind, cronjobTTL, &ttl)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				cronjobTTL: {{Message: "Managed Job specifies ttlSecondsAfterFinished which might conflict with successfulJobsHistoryLimit and failedJobsHistoryLimit from CronJob that have default values. Final behaviour is determined by the strictest parameter, and therefore, setting ttlSecondsAfterFinished at the job level can result with unexpected behaviour with regard to finished jobs removal"}},
			},
			ExpectInstantiationError: false,
		},
	})
}

func (s *JobTTLSecondsAfterFinishedTestSuite) TestCronJobNoTTL() {
	s.AddJobLike(CronJobKind, cronjobNoTTL, nil)
	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				cronjobNoTTL: nil,
			},
			ExpectInstantiationError: false,
		},
	})
}
