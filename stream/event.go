package stream

import "fmt"

// Event matches common JSONL agent traces (OpenAI deep-research streams,
// Arbor events.jsonl, LangGraph / OTel-style logs).
type Event struct {
	Type   string         `json:"type"`
	Event  string         `json:"event"`
	Status string         `json:"status"`
	Meta   map[string]any `json:"data"`
}

// CanonicalState returns the Markov state label for this event.
// Prefers type, then event; appends ":status" when a top-level status is set.
// For Arbor-style cycle.phase events, appends data.phase when present.
func (e Event) CanonicalState() string {
	name := e.Type
	if name == "" {
		name = e.Event
	}
	if name == "" {
		return ""
	}
	if e.Status != "" {
		return fmt.Sprintf("%s:%s", name, e.Status)
	}
	if e.Meta != nil {
		if phase, ok := e.Meta["phase"].(string); ok && phase != "" {
			return fmt.Sprintf("%s:%s", name, phase)
		}
	}
	return name
}
