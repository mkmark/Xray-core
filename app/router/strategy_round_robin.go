package router

import (
	reflect "reflect"
	sync "sync"
)

type RoundRobinStrategy struct {
	mu         sync.Mutex
	tags       []string
	index      int
	roundRobin *RoundRobinStrategy
}

func NewRoundRobin(tags []string) *RoundRobinStrategy {
	return &RoundRobinStrategy{
		tags: tags,
	}
}
func (r *RoundRobinStrategy) NextTag() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	tags := r.tags[r.index]
	r.index = (r.index + 1) % len(r.tags)
	return tags
}

func (s *RoundRobinStrategy) PickOutbound(tags []string) string {
	if len(tags) == 0 {
		panic("0 tags")
	}
	if s.roundRobin == nil || !reflect.DeepEqual(s.roundRobin.tags, tags) {
		s.roundRobin = NewRoundRobin(tags)
	}
	tag := s.roundRobin.NextTag()

	return tag
}
