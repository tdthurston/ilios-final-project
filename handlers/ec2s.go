package handlers

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Ec2Info() string {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String("us-east-1")
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return "Unable to retrieve EC2 data"
	}

	// Extract EC2 Data
	ec2s := []map[string]string{}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			ec2Name := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					ec2Name = *tag.Value
					break
				}
			}

			ec2Data := map[string]string{

				"name":  ec2Name,
				"id":    *instance.InstanceId,
				"type":  *instance.InstanceType,
				"state": *instance.State.Name,
				
			}
			ec2s = append(ec2s, ec2Data)
		}
	}

	// Convert the EC2 data into JSON format
	ec2sJSON, err := json.MarshalIndent(map[string]interface{}{"EC2 Instances": ec2s}, "", "    ")
	if err != nil {
		log.Println("Error marshaling EC2s to JSON:", err)
		return `{"error": "Error processing EC2 data"}`
	}

	return string(ec2sJSON)
}
