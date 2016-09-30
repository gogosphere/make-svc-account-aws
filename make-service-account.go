package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

var (
	fail        = syntax()
	accountrole = os.Args[1]
	svcaccount  = os.Args[2]
	credpath    = os.Getenv("HOME") + "/.aws/credentials"
)

// AWSreq carrier pigeon
type AWSreq struct {
	service *iam.IAM
	params  *string
}

func main() {
	var setup = &AWSreq{}
	setup.authAWS()
	setup.createAccount(setup.service)
}

func assignkey(svc *iam.IAM, svcaccount *string) {
	params := &iam.CreateAccessKeyInput{UserName: svcaccount}
	resp, err := svc.CreateAccessKey(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func (a AWSreq) createAccount(svc *iam.IAM) {
	params := &iam.CreateUserInput{UserName: a.params}
	_, err := svc.CreateUser(params)
	if err != nil {
		fmt.Println("Err in creating the user ", err)
	}
	assignkey(svc, a.params)
}

func syntax() int {
	var k int
	for k, v := range os.Args {
		if v == "" {
			log.Println("Accountrole or ServiceAccount Name not specified.")
			return k
		}
	}
	return k
}

func (a *AWSreq) authAWS() {
	credentialObject := credentials.NewSharedCredentials(credpath, accountrole)
	a.service = iam.New(session.New(aws.NewConfig().WithRegion("us-east-1").
		WithMaxRetries(2).WithCredentials(credentialObject)))
	a.params = aws.String(svcaccount)

}
