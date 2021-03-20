package main

import (
	"strings"
)

// subTargetParams replaces string with the params from the target
func subTargetParams(str string, target Target) string {
	ustr := str
	ustr = strings.ReplaceAll(ustr, "{target}", target.Target)
	ustr = strings.ReplaceAll(ustr, "{host}", target.Host)
	ustr = strings.ReplaceAll(ustr, "{domain}", target.Host)
	ustr = strings.ReplaceAll(ustr, "{hostname}", target.Host)
	ustr = strings.ReplaceAll(ustr, "{proto}", target.Protocol)
	ustr = strings.ReplaceAll(ustr, "{protocol}", target.Protocol)
	ustr = strings.ReplaceAll(ustr, "{query}", target.Querystr)
	ustr = strings.ReplaceAll(ustr, "{querystr}", target.Querystr)
	ustr = strings.ReplaceAll(ustr, "{port}", target.Port)
	ustr = strings.ReplaceAll(ustr, "{basepath}", target.Basepath)
	ustr = strings.ReplaceAll(ustr, "{bpath}", target.Basepath)
	ustr = strings.ReplaceAll(ustr, "{path}", target.Path)
	ustr = strings.ReplaceAll(ustr, "{folder}", target.Folder)
	ustr = strings.ReplaceAll(ustr, "{aws_profile}", target.AWSProfile)
	ustr = strings.ReplaceAll(ustr, "{aws_region}", target.AWSRegion)
	ustr = strings.ReplaceAll(ustr, "{gcp_account}", target.AWSProfile)
	ustr = strings.ReplaceAll(ustr, "{gcp_project}", target.GCPProject)
	ustr = strings.ReplaceAll(ustr, "{gcp_region}", target.GCPRegion)
	ustr = strings.ReplaceAll(ustr, "{gcp_zone}", target.GCPZone)
	return ustr
}