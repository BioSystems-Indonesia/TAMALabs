package runner

import (
	"context"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	a15 "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/smb/A15"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/tcp/ba400"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order/runner/postrun"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order/runner/prerun"
)

type Strategy struct {
	runActionPreRunner    *prerun.RunAction
	cancelActionPreRunner *prerun.CancelAction

	runActionPostRunner            *postrun.RunAction
	cancelActionPostRunner         *postrun.CancelAction
	incompleteSendActionPostRunner *postrun.IncompleteSendAction

	ba400 *ba400.Ba400
	a15   *a15.A15
}

func NewStrategy(
	runActionPreRunner *prerun.RunAction,
	cancelActionPreRunner *prerun.CancelAction,

	runActionPostRunner *postrun.RunAction,
	cancelActionPostRunner *postrun.CancelAction,
	incompleteSendActionPostRunner *postrun.IncompleteSendAction,

	ba400 *ba400.Ba400,
	a15 *a15.A15,
) *Strategy {
	return &Strategy{
		runActionPreRunner:             runActionPreRunner,
		cancelActionPreRunner:          cancelActionPreRunner,
		runActionPostRunner:            runActionPostRunner,
		cancelActionPostRunner:         cancelActionPostRunner,
		incompleteSendActionPostRunner: incompleteSendActionPostRunner,
		ba400:                          ba400,
		a15:                            a15,
	}
}

func (s *Strategy) ChoosePreRunner(ctx context.Context, action constant.WorkOrderRunAction) (usecase.WorkOrderPreRunner, error) {
	switch action {
	case constant.WorkOrderRunActionRun:
		return s.runActionPreRunner, nil
	case constant.WorkOrderRunActionCancel:
		return s.cancelActionPreRunner, nil
	default:
		return nil, fmt.Errorf("unknown action %s", action)
	}
}

func (s *Strategy) ChoosePostRunner(
	ctx context.Context,
	action constant.WorkOrderRunAction,
) (usecase.WorkOrderPostRunner, error) {
	switch action {
	case constant.WorkOrderRunActionRun:
		return s.runActionPostRunner, nil
	case constant.WorkOrderRunActionCancel:
		return s.cancelActionPostRunner, nil
	case constant.WorkOrderRunActionIncompleteSend:
		return s.incompleteSendActionPostRunner, nil
	default:
		return nil, fmt.Errorf("unknown action %s", action)
	}
}
