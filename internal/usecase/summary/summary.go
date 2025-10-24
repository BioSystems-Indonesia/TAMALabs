package summary_uc

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	summaryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/summary"
)

type SummaryUseCase struct {
	summaryRepo *summaryrepo.SummaryRepository
}

func NewSummaryUsecase(sumaryRepo *summaryrepo.SummaryRepository) *SummaryUseCase {
	return &SummaryUseCase{summaryRepo: sumaryRepo}
}

func (u *SummaryUseCase) SummaryAnalytics(ctx context.Context) entity.SummaryResponse {
	workOrderTrend := u.summaryRepo.GetWorkTrendSummary(ctx)
	abnormalSummary := u.summaryRepo.GetAbnormalSummary(ctx)
	mostOrderdTest := u.summaryRepo.GetMostOrderedTest(ctx)
	testTypeDistribution := u.summaryRepo.GetTestTypeDistribution(ctx)
	ageGroupDistribution := u.summaryRepo.GetAgeGroupDistribution(ctx)
	genderDistribution := u.summaryRepo.GetGenderDistribution(ctx)
	return entity.SummaryResponse{
		WorkOrderTrend:       workOrderTrend,
		AbnormalSummary:      abnormalSummary,
		TopTestOrdered:       mostOrderdTest,
		TestTypeDistribution: testTypeDistribution,
		AgeGroup:             ageGroupDistribution,
		GenderSummary:        genderDistribution,
	}

}

func (u *SummaryUseCase) Summary(ctx context.Context) entity.SummaryCardResponse {
	return u.summaryRepo.GetSummary(ctx)
}
