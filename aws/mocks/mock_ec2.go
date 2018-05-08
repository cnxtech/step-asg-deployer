package mocks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/coinbase/step-asg-deployer/aws"
	"github.com/coinbase/step/utils/to"
)

// DescribeSubnetsResponse returns
type DescribeSubnetsResponse struct {
	Resp  *ec2.DescribeSubnetsOutput
	Error error
}

// DescribeImagesResponse returns
type DescribeImagesResponse struct {
	Resp  *ec2.DescribeImagesOutput
	Error error
}

// DescribeSecurityGroupsResponse returns
type DescribeSecurityGroupsResponse struct {
	Resp  *ec2.DescribeSecurityGroupsOutput
	Error error
}

// EC2Client returns
type EC2Client struct {
	aws.EC2API
	DescribeSecurityGroupsResp map[string]*DescribeSecurityGroupsResponse
	DescribeSubnetsResp        *DescribeSubnetsResponse
	DescribeImagesResp         *DescribeImagesResponse
}

func (m *EC2Client) init() {
	if m.DescribeSecurityGroupsResp == nil {
		m.DescribeSecurityGroupsResp = map[string]*DescribeSecurityGroupsResponse{}
	}
}

// AddSecurityGroup returns
func (m *EC2Client) AddSecurityGroup(name string, projectName string, configName string, serviceName string, err error) {
	m.init()
	m.DescribeSecurityGroupsResp[name] = &DescribeSecurityGroupsResponse{
		Resp: &ec2.DescribeSecurityGroupsOutput{
			SecurityGroups: []*ec2.SecurityGroup{
				MakeMockSecurityGroup(projectName, configName, serviceName),
			},
		},
		Error: err,
	}
}

// AddImage returns
func (m *EC2Client) AddImage(nameTag string, id string) {
	m.DescribeImagesResp = &DescribeImagesResponse{
		Resp: &ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				&ec2.Image{
					ImageId: to.Strp(id),
					Tags: []*ec2.Tag{
						&ec2.Tag{Key: to.Strp("Name"), Value: to.Strp(nameTag)},
						&ec2.Tag{Key: to.Strp("DeployWith"), Value: to.Strp("step-asg-deployer")},
					},
				},
			},
		},
	}
}

// AddSubnet returns
func (m *EC2Client) AddSubnet(nameTag string, id string) {
	m.DescribeSubnetsResp = &DescribeSubnetsResponse{
		Resp: &ec2.DescribeSubnetsOutput{
			Subnets: []*ec2.Subnet{
				&ec2.Subnet{
					SubnetId: to.Strp(id),
					Tags: []*ec2.Tag{
						&ec2.Tag{Key: to.Strp("Name"), Value: to.Strp(nameTag)},
						&ec2.Tag{Key: to.Strp("DeployWith"), Value: to.Strp("step-asg-deployer")},
					},
				},
			},
		},
	}
}

// DescribeSecurityGroups returns
func (m *EC2Client) DescribeSecurityGroups(in *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	m.init()
	sgName := in.Filters[1].Values[0]
	resp := m.DescribeSecurityGroupsResp[*sgName]
	if resp == nil {
		return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: []*ec2.SecurityGroup{}}, nil
	}
	return resp.Resp, resp.Error
}

// MakeMockSecurityGroup returns
func MakeMockSecurityGroup(projectName string, configName string, serviceName string) *ec2.SecurityGroup {
	return &ec2.SecurityGroup{
		GroupId: to.Strp("group-id"),
		Tags: []*ec2.Tag{
			&ec2.Tag{Key: to.Strp("ProjectName"), Value: to.Strp(projectName)},
			&ec2.Tag{Key: to.Strp("ConfigName"), Value: to.Strp(configName)},
			&ec2.Tag{Key: to.Strp("ServiceName"), Value: to.Strp(serviceName)},
		},
	}
}

// DescribeSubnets returns
func (m *EC2Client) DescribeSubnets(in *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	if m.DescribeSubnetsResp == nil {
		return nil, fmt.Errorf("Add Subnets")
	}

	return m.DescribeSubnetsResp.Resp, m.DescribeSubnetsResp.Error
}

// DescribeImages returns
func (m *EC2Client) DescribeImages(in *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	if m.DescribeImagesResp == nil {
		return nil, fmt.Errorf("Add Image")
	}

	return m.DescribeImagesResp.Resp, m.DescribeImagesResp.Error
}
