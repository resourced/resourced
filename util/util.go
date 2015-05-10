// Package contains miscelaneous utility functions
package util

import (
	"net"
	"strings"
)

// Parse a string with comma separated CIDR's; return slice of net.IPNet objs;
// if given string is empty, return an empty slice of net.IPNet objs instead.
func ParseCIDRs(cidrs string) ([]*net.IPNet, error) {
	if cidrs == "" {
		return []*net.IPNet{}, nil
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
