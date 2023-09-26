package toggl

import (
	"math"
	"time"
)

type ReportService struct {
	timeService    TimeService
	projectService ProjectService
}

func NewReportService(timeService TimeService, projectService ProjectService) ReportService {
	return ReportService{
		timeService:    timeService,
		projectService: projectService,
	}
}

func (r *ReportService) BuildReport(sd, ed time.Time) (*Report, error) {
	builder := NewReportBuilder(r.timeService, r.projectService)

	return builder.BuildReport(sd, ed)
}

func (r *ReportService) BuildStringReport(model *Report) (*string, error) {
	rpt := NewStringReport()
	return rpt.Generate(model), nil
}

func (r *ReportService) BuildJsonReport(model *Report) (*[]byte, error) {
	rpt := NewJsonReport()
	return rpt.Generate(model)

}

func roundToQuarterHour(duration float64) float64 {
	rounded := math.Ceil(duration*4) / 4.0
	return rounded
}
