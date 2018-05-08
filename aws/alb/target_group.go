package alb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/coinbase/step-asg-deployer/aws"
	"github.com/coinbase/step/utils/to"
)

// TargetGroup struct
type TargetGroup struct {
	ProjectNameTag *string
	ConfigNameTag  *string
	ServiceNameTag *string
	TargetGroupArn *string
}

// ProjectName returns tag
func (s *TargetGroup) ProjectName() *string {
	return s.ProjectNameTag
}

// ConfigName returns tag
func (s *TargetGroup) ConfigName() *string {
	return s.ConfigNameTag
}

// ServiceName returns tag
func (s *TargetGroup) ServiceName() *string {
	return s.ServiceNameTag
}

//////
// Healthy
//////

// GetInstances return instances on the target group
func GetInstances(albc aws.ALBAPI, arn *string, instances []string) (aws.Instances, error) {
	healthOutput, err := albc.DescribeTargetHealth(createDescribeTargetHealthInput(arn, instances))

	if err != nil {
		return nil, err
	}

	tgInstances := aws.Instances{}
	for _, thd := range healthOutput.TargetHealthDescriptions {
		tgInstances.AddTargetGroupInstance(thd)
	}

	return tgInstances, nil
}

func createDescribeTargetHealthInput(arn *string, instances []string) *elbv2.DescribeTargetHealthInput {
	awsInstances := []*elbv2.TargetDescription{}
	for _, id := range instances {
		awsInstances = append(awsInstances, &elbv2.TargetDescription{Id: to.Strp(id)})
	}

	return &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: arn,
		Targets:        awsInstances,
	}
}

//////
// Find
//////

// FindAll returns all target groups in a list
func FindAll(albc aws.ALBAPI, names []*string) ([]*TargetGroup, error) {
	tgs := []*TargetGroup{}
	for _, name := range names {
		tg, err := find(albc, name)
		if err != nil {
			return nil, err
		}
		tgs = append(tgs, tg)
	}

	return tgs, nil
}

func find(alb aws.ALBAPI, targetGroupName *string) (*TargetGroup, error) {
	awsTarget, err := findByName(alb, targetGroupName)
	if err != nil {
		return nil, err
	}

	awsTags, err := findTagsByName(alb, awsTarget.TargetGroupArn)
	if err != nil {
		return nil, err
	}

	return &TargetGroup{
		ProjectNameTag: aws.FetchELBV2Tag(awsTags, to.Strp("ProjectName")),
		ConfigNameTag:  aws.FetchELBV2Tag(awsTags, to.Strp("ConfigName")),
		ServiceNameTag: aws.FetchELBV2Tag(awsTags, to.Strp("ServiceName")),
		TargetGroupArn: awsTarget.TargetGroupArn,
	}, nil
}

func findByName(alb aws.ALBAPI, targetGroupName *string) (*elbv2.TargetGroup, error) {
	elbsOutput, err := alb.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
		Names: []*string{targetGroupName},
	})

	if err != nil {
		return nil, err
	}

	if len(elbsOutput.TargetGroups) != 1 {
		return nil, fmt.Errorf("LoadBalancer Not Found")
	}

	if *elbsOutput.TargetGroups[0].TargetGroupName != *targetGroupName {
		return nil, fmt.Errorf("LoadBalancer Not Found")
	}

	return elbsOutput.TargetGroups[0], nil
}

func findTagsByName(alb aws.ALBAPI, targetGroupARN *string) ([]*elbv2.Tag, error) {
	tagsOutput, err := alb.DescribeTags(&elbv2.DescribeTagsInput{
		ResourceArns: []*string{targetGroupARN},
	})

	if err != nil {
		return nil, err
	}

	if len(tagsOutput.TagDescriptions) != 1 {
		return nil, fmt.Errorf("TargetGroup Not Found")
	}

	if *tagsOutput.TagDescriptions[0].ResourceArn != *targetGroupARN {
		return nil, fmt.Errorf("TargetGroup Not Found")
	}

	return tagsOutput.TagDescriptions[0].Tags, nil
}
