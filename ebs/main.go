package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"strings"
)

func main() {


	mySession := session.Must(session.NewSession())

	// Create a EC2 client from just a session.
	svc := ec2.New(mySession)

	r, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{
		Filters:[]*ec2.Filter{{
			Name:   aws.String("status"),
			Values: []*string{aws.String("available")},
		}},
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}

	fmt.Printf("found %d volumes\n", len(r.Volumes))

	if !(len(os.Args) > 1 && strings.EqualFold(os.Args[1], "delete")) {
		fmt.Printf("dry run only.\n")
		return
	}

	for _, v := range r.Volumes {
		fmt.Printf("%s: %s\n", *v.VolumeId, *v.State)
		fmt.Printf("deleting volume %s\n", *v.VolumeId)
		_, err := svc.DeleteVolume(&ec2.DeleteVolumeInput{
			VolumeId: v.VolumeId,
		})
		if err != nil {
			fmt.Printf("could not delete volume %s, error: %v\n", *v.VolumeId, err)
			continue
		}
	}

	fmt.Printf("done")

}