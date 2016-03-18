package agent

import (
	"testing"
)

func TestSetAccessTokensDuringConstructor(t *testing.T) {
	agent := createAgentWithAccessTokensForTest(t)

	if len(agent.AccessTokens) <= 0 {
		t.Errorf("agent.AccessTokens should not be empty")
	}
}

func TestIsAllowed(t *testing.T) {
	agent := createAgentWithAccessTokensForTest(t)

	givenToken := "2fd4e1c67a2d28fced849ee1bb76e7391b93eb12"

	if !agent.IsAllowed(givenToken) {
		t.Errorf("IsAllowed is wrong. GivenToken: %v. AccessTokens: %v", givenToken, agent.AccessTokens)
	}
}
