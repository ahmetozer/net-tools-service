package main

import (
	"fmt"
	"regexp"
)

func isMayIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex+`|`+domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayOnlyIPv4(host string) bool {
	match, err := regexp.MatchString(ipv4Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex+`|`+domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}
func isMayOnlyIPv6(host string) bool {
	match, err := regexp.MatchString(ipv6Regex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isHTTPURLScheme(scheme string) bool {
	match, err := regexp.MatchString(`^http$|^https$`, scheme)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isPortValid(port string) bool {
	match, err := regexp.MatchString(portRegex, port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}

func isMayDomain(host string) bool {
	match, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return match
}