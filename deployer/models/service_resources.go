package models

import (
	"fmt"

	"github.com/coinbase/step-asg-deployer/aws"
	"github.com/coinbase/step-asg-deployer/aws/alb"
	"github.com/coinbase/step-asg-deployer/aws/ami"
	"github.com/coinbase/step-asg-deployer/aws/asg"
	"github.com/coinbase/step-asg-deployer/aws/elb"
	"github.com/coinbase/step-asg-deployer/aws/iam"
	"github.com/coinbase/step-asg-deployer/aws/sg"
	"github.com/coinbase/step-asg-deployer/aws/subnet"
	"github.com/coinbase/step/utils/is"
	"github.com/coinbase/step/utils/to"
)

// This models all resources referenced for a release
type serviceIface interface {
	ProjectName() *string
	ConfigName() *string
	Name() *string
	ReleaseID() *string
}

// ServiceResources struct
type ServiceResources struct {
	Image          *ami.Image
	Profile        *iam.Profile
	PrevASG        *asg.ASG
	SecurityGroups []*sg.SecurityGroup
	ELBs           []*elb.LoadBalancer
	TargetGroups   []*alb.TargetGroup
	Subnets        []*subnet.Subnet
}

// ServiceResourceNames struct
type ServiceResourceNames struct {
	Image          *string   `json:"image,omitempty"`
	Profile        *string   `json:"profile_arn,omitempty"`
	PrevASG        *string   `json:"prev_asg_arn,omitempty"`
	SecurityGroups []*string `json:"security_groups,omitempty"`
	ELBs           []*string `json:"elbs,omitempty"`
	TargetGroups   []*string `json:"target_group_arns,omitempty"`
	Subnets        []*string `json:"subnets,omitempty"`
}

// ToServiceResourceNames returns
func (sr *ServiceResources) ToServiceResourceNames() *ServiceResourceNames {
	var im *string
	if sr.Image != nil {
		im = sr.Image.ImageID
	}

	var profile *string
	if sr.Profile != nil {
		profile = sr.Profile.Arn
	}

	var prevASG *string
	if sr.PrevASG != nil {
		prevASG = sr.PrevASG.AutoScalingGroupName
	}

	sgs := []*string{}
	for _, sg := range sr.SecurityGroups {
		if sg == nil || is.EmptyStr(sg.GroupID) {
			continue
		}

		sgs = append(sgs, sg.GroupID)
	}

	elbs := []*string{}
	for _, elb := range sr.ELBs {
		if elb == nil || is.EmptyStr(elb.LoadBalancerName) {
			continue
		}

		elbs = append(elbs, elb.LoadBalancerName)
	}

	tgs := []*string{}
	for _, tg := range sr.TargetGroups {
		if tg == nil || is.EmptyStr(tg.TargetGroupArn) {
			continue
		}

		tgs = append(tgs, tg.TargetGroupArn)
	}

	subnets := []*string{}
	for _, subnet := range sr.Subnets {
		if subnet == nil || is.EmptyStr(subnet.SubnetID) {
			continue
		}

		subnets = append(subnets, subnet.SubnetID)
	}

	return &ServiceResourceNames{
		Image:          im,
		Profile:        profile,
		PrevASG:        prevASG,
		SecurityGroups: sgs,
		ELBs:           elbs,
		TargetGroups:   tgs,
		Subnets:        subnets,
	}
}

// Validate returns
func (sr *ServiceResources) Validate(service *Service) error {

	if err := sr.validateAttributes(service); err != nil {
		return err
	}

	if err := ValidateImage(service, sr.Image); err != nil {
		return err
	}

	// Now the Easy Validations are over time to validate Tags and Paths
	if err := ValidateIAMProfile(service, sr.Profile); err != nil {
		return err
	}

	if err := ValidatePrevASG(service, sr.PrevASG); err != nil {
		return err
	}

	for _, r := range sr.SecurityGroups {
		if err := ValidateSecurityGroup(service, r); err != nil {
			return err
		}
	}

	for _, r := range sr.ELBs {
		if err := ValidateELB(service, r); err != nil {
			return err
		}
	}

	for _, r := range sr.TargetGroups {
		if err := ValidateTargetGroup(service, r); err != nil {
			return err
		}
	}

	for _, r := range sr.Subnets {
		if err := ValidateSubnet(service, r); err != nil {
			return err
		}
	}

	return nil
}

func (sr *ServiceResources) validateAttributes(service *Service) error {
	names := sr.ToServiceResourceNames()

	// Must have Image
	if sr.Image == nil {
		return fmt.Errorf("Image is nil")
	}

	// Must have the correct amount of security groups, ELBS and Target Groups
	if len(service.SecurityGroups) != len(sr.SecurityGroups) {
		return fmt.Errorf("Security Group Not Found actual %v expected %v", to.StrSlice(names.SecurityGroups), to.StrSlice(service.SecurityGroups))
	}

	if len(service.ELBs) != len(sr.ELBs) {
		return fmt.Errorf("ELB Not Found actual %v expected %v", to.StrSlice(names.ELBs), to.StrSlice(service.ELBs))
	}

	if len(service.TargetGroups) != len(sr.TargetGroups) {
		return fmt.Errorf("TargetGroup Not Found actual %v expected %v", to.StrSlice(names.TargetGroups), to.StrSlice(service.TargetGroups))
	}

	if len(service.Subnets()) != len(sr.Subnets) {
		return fmt.Errorf("Subnets Not Found actual %v expected %v", to.StrSlice(names.Subnets), to.StrSlice(service.Subnets()))
	}

	return nil
}

// ValidateImage returns
func ValidateImage(service serviceIface, im *ami.Image) error {
	if im == nil {
		return fmt.Errorf("Image is nil")
	}

	if im.DeployWithTag == nil {
		return fmt.Errorf("Image %v DeployWith Tag nil", *im.ImageID)
	}

	if *im.DeployWithTag != "step-asg-deployer" {
		return fmt.Errorf("Image %v DeployWith Tag expected: %v actual: %v", *im.ImageID, "step-asg-deployer", *im.DeployWithTag)
	}

	return nil
}

// ValidateSubnet returns
func ValidateSubnet(service serviceIface, subnet *subnet.Subnet) error {
	if subnet == nil {
		return fmt.Errorf("Subnet is nil")
	}

	if subnet.DeployWithTag == nil {
		return fmt.Errorf("Subnet %v DeployWith Tag nil", *subnet.SubnetID)
	}

	if *subnet.DeployWithTag != "step-asg-deployer" {
		return fmt.Errorf("Subnet %v DeployWith Tag expected: %v actual: %v", *subnet.SubnetID, "step-asg-deployer", *subnet.DeployWithTag)
	}

	return nil
}

// ValidatePrevASG returns
func ValidatePrevASG(service serviceIface, as *asg.ASG) error {
	if as == nil {
		return nil // Allowed to not have previous ASG
	}

	// None of these should happen but it is just being extra safe
	if !aws.HasProjectName(as, service.ProjectName()) {
		return fmt.Errorf("Previous ASG incorrect ProjectName requires %q has %q", to.Strs(service.ProjectName()), to.Strs(as.ProjectName()))
	}

	if !aws.HasConfigName(as, service.ConfigName()) {
		return fmt.Errorf("Previous ASG incorrect ConfigName requires %q has %q", to.Strs(service.ConfigName()), to.Strs(as.ConfigName()))
	}

	if !aws.HasServiceName(as, service.Name()) {
		return fmt.Errorf("Previous ASG incorrect ServiceName requires %q has %q", to.Strs(service.Name()), to.Strs(as.ServiceName()))
	}

	if as.ReleaseID() == nil {
		return fmt.Errorf("Previous ASG ReleaseID nil")
	}

	// Existing ASG must not have the same Release ID
	if *as.ReleaseID() == *service.ReleaseID() {
		return fmt.Errorf("Previous ASG incorrect ReleaseID requires %q has %q", to.Strs(service.ReleaseID()), to.Strs(as.ReleaseID()))
	}

	return nil
}

// ValidateIAMProfile returns
func ValidateIAMProfile(service serviceIface, profile *iam.Profile) error {
	if profile == nil {
		return nil // Profile is allowed to be nil
	}

	if profile.Path == nil {
		// Again should never happen
		return fmt.Errorf("Iam Profile Path not found")
	}

	validPath := fmt.Sprintf("/%v/%v/%v/", *service.ProjectName(), *service.ConfigName(), *service.Name())
	if *profile.Path != validPath {
		// Again should never happen
		return fmt.Errorf("Iam Profile Path incorrect, it is %q and requires %q", *profile.Path, validPath)
	}

	return nil
}

// ValidateSecurityGroup returns
func ValidateSecurityGroup(service serviceIface, sc *sg.SecurityGroup) error {
	if sc == nil {
		return fmt.Errorf("SecurityGroup is nil")
	}

	if !aws.HasProjectName(sc, service.ProjectName()) {
		return fmt.Errorf("Security Group incorrect ProjectName requires %q has %q", to.Strs(service.ProjectName()), to.Strs(sc.ProjectName()))
	}

	if !aws.HasConfigName(sc, service.ConfigName()) {
		return fmt.Errorf("Security Group incorrect ConfigName requires %q has %q", to.Strs(service.ConfigName()), to.Strs(sc.ConfigName()))
	}

	if !aws.HasServiceName(sc, service.Name()) {
		return fmt.Errorf("Security Group incorrect ServiceName requires %q has %q", to.Strs(service.Name()), to.Strs(sc.ServiceName()))
	}

	return nil
}

// ValidateELB returns
func ValidateELB(service serviceIface, lb *elb.LoadBalancer) error {
	if lb == nil {
		return fmt.Errorf("LoadBalancer is nil")
	}

	if !aws.HasProjectName(lb, service.ProjectName()) {
		return fmt.Errorf("ELB incorrect ProjectName requires %q has %q", to.Strs(service.ProjectName()), to.Strs(lb.ProjectName()))
	}

	if !aws.HasConfigName(lb, service.ConfigName()) {
		return fmt.Errorf("ELB incorrect ConfigName requires %q has %q", to.Strs(service.ConfigName()), to.Strs(lb.ConfigName()))
	}

	if !aws.HasServiceName(lb, service.Name()) {
		return fmt.Errorf("ELB incorrect ServiceName requires %q has %q", to.Strs(service.Name()), to.Strs(lb.ServiceName()))
	}

	return nil
}

// ValidateTargetGroup returns
func ValidateTargetGroup(service serviceIface, tg *alb.TargetGroup) error {
	if tg == nil {
		return fmt.Errorf("TargetGroup is nil")
	}

	if !aws.HasProjectName(tg, service.ProjectName()) {
		return fmt.Errorf("Target Group incorrect ProjectName requires %q has %q", to.Strs(service.ProjectName()), to.Strs(tg.ProjectName()))
	}

	if !aws.HasConfigName(tg, service.ConfigName()) {
		return fmt.Errorf("Target Group incorrect ConfigName requires %q has %q", to.Strs(service.ConfigName()), to.Strs(tg.ConfigName()))
	}

	if !aws.HasServiceName(tg, service.Name()) {
		return fmt.Errorf("Target Group incorrect ServiceName requires %q has %q", to.Strs(service.Name()), to.Strs(tg.ServiceName()))
	}

	return nil
}
