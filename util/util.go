// Package contains miscelaneous utility functions
package util

import (
	"net"
	"strings"
)

// Parse a string with comma separated CIDR's; if given string is empty, return
// a slice with a single 'default' 0.0.0.0 CIDR. Return slice of net.IPNet objs
func ParseCIDRs(cidrs string) ([]*net.IPNet, error) {
	if cidrs == "" {
		_, defaultCIDR, _ := net.ParseCIDR("0.0.0.0/0")
		return []*net.IPNet{defaultCIDR}, nil
	}

	// Get rid of spaces
	cidrs = strings.Replace(cidrs, " ", "", -1)

	// Convert cidr strings to net.IPNet objects
	converted := []*net.IPNet{}

	for _, value := range strings.Split(cidrs, ",") {
		_, newCIDR, err := net.ParseCIDR(value)
		if err != nil {
			return converted, err
		}

		converted = append(converted, newCIDR)
	}

	return converted, nil
}
