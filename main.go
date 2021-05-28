package main

// Rough steps
// Read arguments to figure out AWS profile to use
// Read ~/.aws/configuration to tie profile to role ARN
// call aws sts assume-role via SDK
// pluck credentials from response
// Start new process with credentials from response set as environment variables (???)

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"gopkg.in/ini.v1"
)

func main() {
	var profileName string
	var mfaToken string
	flag.StringVar(&profileName, "profile", "", "AWS profile to use for subcommand")
	flag.StringVar(&mfaToken, "token", "", "MFA token to use when assuming role")

	flag.Parse()
	fmt.Printf("Profile is: %v \n", profileName)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homeDir)

	profiles, err := ini.Load(homeDir + "/.aws/config")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		os.Exit(1)
	}

	defaultRegion := profiles.Section("default").Key("region").String()
	targetProfile := profiles.Section("profile " + profileName)
	fmt.Println("region of default profile:", defaultRegion)
	fmt.Println("role arn of selected profile:", targetProfile.Key("role_arn").String())
	if targetProfile.HasKey("mfa_serial") {
		if len(mfaToken) == 0 {
			mfa_serial := targetProfile.Key("mfa_serial")
			fmt.Printf("Please enter MFA token for device: %v \n", mfa_serial)
			fmt.Scanln(&mfaToken)
			if len(mfaToken) != 6 {
				log.Fatalf("Expected 6 characters in MFA token")
			}
		}
	}

	fmt.Println("mfatoken = " + mfaToken)

	// stsOptions := sts.Options{
	// 	Region: defaultRegion,
	// }

	// stsClient := sts.New(stsOptions)
	// stsClient.AssumeRole()
}
