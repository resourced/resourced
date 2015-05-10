package util

import (
	"testing"
)

func TestParseCIDRs(t *testing.T) {
	goodCIDRs := []string{"127.0.0.1/8", "127.0.0.1/8, 0.0.0.0/0", "  0.0.0.0/0,  127.0.0.1/24 "}
	badCIDRs := []string{"127.0.0.1", "127.0.0.1/99", "127.0.0.2/8 127.0.0.1/24"}

	for _, goodCIDR := range goodCIDRs {
		_, err := ParseCIDRs(goodCIDR)
		if err != nil {
			t.Errorf("'%v' should pass as a good CIDR value. Err: %s", goodCIDR, err)
		}
	}

	for _, badCIDR := range badCIDRs {
		_, err := ParseCIDRs(badCIDR)
		if err == nil {
			t.Errorf("'%v' should NOT pass as proper CIDR value. Err: %s", badCIDR, err)
		}
	}

	cidrs, err := ParseCIDRs(goodCIDRs[0])
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(cidrs) != 1 {
		t.Errorf("'cidrs' should have 1 element, but has %v", len(cidrs))
	}
}
