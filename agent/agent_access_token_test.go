package agent

import (
	"testing"
)

func TestSetAccessTokensDuringConstructor(t *testing.T) {
	agent := createAgentForTest(t)

	allTokens := len(agent.AccessTokens)

	if allTokens == 0 {
		t.Error("agent.AccessTokens should not be empty.")
	}

	correctTokensFound := 0

	for _, token := range agent.AccessTokens {
		if token == "does-not-exist" {
			t.Errorf("agent.AccessTokens is incorrect: %v", agent.AccessTokens)
		}

		if token == "abc123" {
			correctTokensFound = correctTokensFound + 1
		}
		if token == "password1" {
			correctTokensFound = correctTokensFound + 1
		}
		if token == "these-are-superbad-passwords" {
			correctTokensFound = correctTokensFound + 1
		}
		if token == "pretend-all-of-them-are-hashed" {
			correctTokensFound = correctTokensFound + 1
		}
	}

	if correctTokensFound != allTokens {
		t.Errorf("agent.AccessTokens is incorrect: %v", agent.AccessTokens)
	}
}
