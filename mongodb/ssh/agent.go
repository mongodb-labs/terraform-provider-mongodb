package ssh

// Agent holds SSH agent configuration parameters
type Agent struct {
	Agent         bool   `json:"agent,string,omitempty"`
	AgentIdentity string `json:"agent_identity,omitempty"`
}
