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
		panic(err)
	}

	return vpcs.Vpcs, nil
}

// GetVpc : Get a vpc by vpc-id
func (client *EC2Handler) GetVpc(vpcID *string) (*ec2.Vpc, error) {
	dryRun := false
	vpcInput := &ec2.DescribeVpcsInput{
		DryRun: &dryRun,
		VpcIds: []*string{vpcID},
	}
	vpcs, err := client.conn.DescribeVpcs(vpcInput)
	if err != nil {
		panic(err)
	}

	return vpcs.Vpcs[0], nil
}

// NewEC2Client : Create a new ec2 client
func NewEC2Client(region string) (*EC2Handler, error) {
	sess, err := newSession()
	if err != nil {
		panic(err)
	}

	return &EC2Handler{conn: ec2.New(sess, aws.NewConfig().WithRegion(region))}, nil
}
