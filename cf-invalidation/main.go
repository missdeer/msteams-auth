package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := cloudfront.New(sess)
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String("E31ILSCIMHAQ6K"),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			Paths: &cloudfront.Paths{
				Items:    []*string{aws.String("/.well-known/microsoft-identity-association.json")},
				Quantity: aws.Int64(1),
			},
			CallerReference: aws.String(fmt.Sprint(time.Now().Unix())),
		},
	}

	result, err := svc.CreateInvalidation(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeAccessDenied:
				fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
			case cloudfront.ErrCodeMissingBody:
				fmt.Println(cloudfront.ErrCodeMissingBody, aerr.Error())
			case cloudfront.ErrCodeInvalidArgument:
				fmt.Println(cloudfront.ErrCodeInvalidArgument, aerr.Error())
			case cloudfront.ErrCodeNoSuchDistribution:
				fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
			case cloudfront.ErrCodeBatchTooLarge:
				fmt.Println(cloudfront.ErrCodeBatchTooLarge, aerr.Error())
			case cloudfront.ErrCodeTooManyInvalidationsInProgress:
				fmt.Println(cloudfront.ErrCodeTooManyInvalidationsInProgress, aerr.Error())
			case cloudfront.ErrCodeInconsistentQuantities:
				fmt.Println(cloudfront.ErrCodeInconsistentQuantities, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
