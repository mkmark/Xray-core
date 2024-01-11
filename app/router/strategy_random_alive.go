package router

import (
	"context"

	"github.com/xtls/xray-core/app/observatory"
	"github.com/xtls/xray-core/common/dice"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/extension"
	"google.golang.org/protobuf/proto"
)

// RandomAliveStrategy represents a random alive balancing strategy
type RandomAliveStrategy struct {
	ctx         context.Context
	observatory extension.Observatory
}

func (s *RandomAliveStrategy) InjectContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *RandomAliveStrategy) PickOutbound(candidates []string) string {
	// candidates are considered alive unless observed otherwise
	if s.observatory == nil {
		core.RequireFeatures(s.ctx, func(observatory extension.Observatory) error {
			s.observatory = observatory
			return nil
		})
	}
	if s.observatory != nil {
		var observeReport proto.Message
		var err error
		observeReport, err = s.observatory.GetObservation(s.ctx)
		if err == nil {
			aliveTags := make([]string, 0)
			if result, ok := observeReport.(*observatory.ObservationResult); ok {
				status := result.Status
				statusMap := make(map[string]*observatory.OutboundStatus)
				for _, outboundStatus := range status {
					statusMap[outboundStatus.OutboundTag] = outboundStatus
				}
				for _, candidate := range candidates {
					if outboundStatus, found := statusMap[candidate]; found {
						if outboundStatus.Alive {
							aliveTags = append(aliveTags, candidate)
						}
					} else {
						// unfound candidate is considered alive
						aliveTags = append(aliveTags, candidate)
					}
				}

				candidates = aliveTags
			}
		}
	}

	count := len(candidates)
	if count == 0 {
		return ""
	}
	return candidates[dice.Roll(count)]
}
