package claw

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2Handler : EC2 struct
type EC2Handler struct {
	conn *ec2.EC2
}

// GetAllVpcs : Get all vps
func (client *EC2Handler) GetAllVpcs() ([]*ec2.Vpc, error) {
	vpcInput := &ec2.DescribeVpcsInput{}
	vpcs, err := client.conn.DescribeVpcs(vpcInput)
	if err != nil {
		return nil, err
	}

	return vpcs.Vpcs, nil
}

// GetVpc : Get a vpc by vpc-id
func (client *EC2Handler) GetVpc(vpcID string) (*ec2.Vpc, error) {
	params := &ec2.DescribeVpcsInput{
		DryRun: aws.Bool(false),
		VpcIds: []*string{
			aws.String(vpcID),
		},
	}
	vpcs, err := client.conn.DescribeVpcs(params)
	if err != nil {
		return nil, err
	}

	return vpcs.Vpcs[0], nil
}

// GetStackInstances : Get all instances on a stack
func (client *EC2Handler) GetStackInstances(clusterName string) ([]*ec2.Instance, error) {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:opsworks:stack"),
				Values: []*string{
					aws.String(clusterName),
				},
			},
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	resp, err := client.conn.DescribeInstances(params)
	if err != nil {
		return nil, err
	}

	var instances []*ec2.Instance
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// NewEC2Client : Create a new ec2 client
func NewEC2Client(region string) (*EC2Handler, error) {
	sess, err := newSession()
	if err != nil {
		panic(err)
	}

	return &EC2Handler{conn: ec2.New(sess, aws.NewConfig().WithRegion(region))}, nil
}
