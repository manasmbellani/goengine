package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// normalizeTarget is used to convert raw target string to target structs
func normalizeTarget(rawTarget string, target *Target) {
	// Check the format of the rawTarget
	rawTargetParts := strings.Split(rawTarget, "://")
	var rawTargetWithoutProtocol string
	if len(rawTargetParts) < 2 {
		log.Fatalf("rawTarget must follow format protocol://....")
	} else {
		rawTargetWithoutProtocol = strings.Split(rawTarget, "://")[1]
	}

	// Get the protocol
	target.Protocol = rawTargetParts[0]
	
	// Parse the folder input if we are dealing with a folder, and add trailing
	// '/' if not present
	if target.Protocol == "folder" {
		target.Folder = rawTargetParts[1]
		if !strings.HasSuffix(target.Folder, "/") {
			target.Folder = target.Folder + "/"
		}
	} else if target.Protocol == "aws" {
		// Parse AWS account information - Profile/region
		targetParts := strings.Split(rawTargetWithoutProtocol, ":")
		if len(targetParts) < 2 {
			log.Fatalf("AWS target spec: aws://<aws_profile>:<aws_region>")
		} else {
			target.AWSProfile = targetParts[0]
			target.AWSRegion = targetParts[1]
		}
	} else if target.Protocol == "gcp" {
		// Parse as GCP account information - GCP account, project, region and 
		// zone
		targetParts := strings.Split(rawTargetWithoutProtocol, ":")
		if len(targetParts) < 4 {
			log.Fatalf("GCP target spec: gcp://<gcp_account>:<gcp_project>:<gcp_region>:<gcp_zone>")
		} else {
			target.GCPAccount = targetParts[0]
			target.GCPProject = targetParts[1]
			target.GCPRegion = targetParts[2]
			target.GCPZone = targetParts[3]
		}
	} else {
		// Parse as URL
		u, err := url.Parse(rawTarget)
		if err != nil {
			log.Printf("Error parsing target: %s as URL. Error: %s\n", rawTarget, err)
		}
		target.Target = rawTarget
		target.Host = u.Host
		target.Port = u.Port()
		if target.Protocol == "" {
			target.Protocol = DefProtocol
		}
		if target.Port == "" {
			if target.Protocol == "https" {
				target.Port = "443"
			} else if target.Protocol == "http" {
				target.Port = "80"
			} else {
				target.Port = DefPort
			}
		}
		target.Path = target.Protocol + "://" + target.Host + ":" + target.Port + u.Path
		queryMap := u.Query()
		var qm []string
		s := ""
		for k, v := range queryMap {
			qm = append(qm, fmt.Sprintf(s, "%s=%s", k, v))
		}
		target.Querystr = strings.Join(qm, "&")
		// Check if path contains "/" and it is a filename containing a dot
		// If it is calculate the basepath, if not, then basepath is
		basepathOnly := ""
		if strings.Contains(u.Path, "/") {
			pathSplits := strings.Split(u.Path, "/")
			pathSplitsLen := len(pathSplits)
			filename := pathSplits[pathSplitsLen-1]

			if strings.Contains(filename, ".") {
				basepathOnly = strings.Join(pathSplits[0:pathSplitsLen-1], "/")
			} else {
				basepathOnly = u.Path
			}
		} else {
			basepathOnly = u.Path
		}
		target.Basepath = target.Protocol + "://" + target.Host + ":" + target.Port + basepathOnly

		// Remove the last char '/' in basepath and path
		if strings.HasSuffix(target.Basepath, "/") {
			target.Basepath = target.Basepath[0 : len(target.Basepath)-1]
		}
		if strings.HasSuffix(target.Path, "/") {
			target.Path = target.Path[0 : len(target.Path)-1]
		}
	}
}
