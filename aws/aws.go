package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sfn/sfniface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	ar "github.com/coinbase/step/aws"
)

// FetchEc2Tag extracts tags
func FetchEc2Tag(tags []*ec2.Tag, tagKey *string) *string {
	if tagKey == nil {
		return nil
	}

	for _, tag := range tags {
		if tag.Key == nil {
			continue
		}
		if *tag.Key == *tagKey {
			return tag.Value
		}
	}

	return nil
}

// FetchELBTag extracts tags
func FetchELBTag(tags []*elb.Tag, tagKey *string) *string {
	if tagKey == nil {
		return nil
	}

	for _, tag := range tags {
		if tag.Key == nil {
			continue
		}
		if *tag.Key == *tagKey {
			return tag.Value
		}
	}

	return nil
}

// FetchELBV2Tag extracts tags
func FetchELBV2Tag(tags []*elbv2.Tag, tagKey *string) *string {
	if tagKey == nil {
		return nil
	}

	for _, tag := range tags {
		if tag.Key == nil {
			continue
		}
		if *tag.Key == *tagKey {
			return tag.Value
		}
	}

	return nil
}

// FetchASGTag extracts tags
func FetchASGTag(tags []*autoscaling.TagDescription, tagKey *string) *string {
	if tagKey == nil {
		return nil
	}

	for _, tag := range tags {
		if tag.Key == nil {
			continue
		}
		if *tag.Key == *tagKey {
			return tag.Value
		}
	}

	return nil
}

// HasProjectName checks value
func HasProjectName(r interface {
	ProjectName() *string
}, projectName *string) bool {
	if r.ProjectName() == nil || projectName == nil {
		return false
	}
	return *r.ProjectName() == *projectName
}

// HasConfigName checks value
func HasConfigName(r interface {
	ConfigName() *string
}, configName *string) bool {
	if r.ConfigName() == nil || configName == nil {
		return false
	}
	return *r.ConfigName() == *configName
}

// HasServiceName checks value
func HasServiceName(r interface {
	ServiceName() *string
}, serviceName *string) bool {
	if r.ServiceName() == nil || serviceName == nil {
		return false
	}
	return *r.ServiceName() == *serviceName
}

// HasReleaseID checks value
func HasReleaseID(r interface {
	ReleaseID() *string
}, releaseID *string) bool {
	if r.ReleaseID() == nil || releaseID == nil {
		return false
	}
	return *r.ReleaseID() == *releaseID
}

// S3API aws API
type S3API s3iface.S3API

// ASGAPI aws API
type ASGAPI autoscalingiface.AutoScalingAPI

// ELBAPI aws API
type ELBAPI elbiface.ELBAPI

// EC2API aws API
type EC2API ec2iface.EC2API

// ALBAPI aws API
type ALBAPI elbv2iface.ELBV2API

// CWAPI aws API
type CWAPI cloudwatchiface.CloudWatchAPI

// IAMAPI aws API
type IAMAPI iamiface.IAMAPI

// SNSAPI aws API
type SNSAPI snsiface.SNSAPI

// SFNAPI aws API
type SFNAPI sfniface.SFNAPI

// Clients for AWS
type Clients interface {
	S3Client(region *string, accountID *string, role *string) S3API
	ASGClient(region *string, accountID *string, role *string) ASGAPI
	ELBClient(region *string, accountID *string, role *string) ELBAPI
	EC2Client(region *string, accountID *string, role *string) EC2API
	ALBClient(region *string, accountID *string, role *string) ALBAPI
	CWClient(region *string, accountID *string, role *string) CWAPI
	IAMClient(region *string, accountID *string, role *string) IAMAPI
	SNSClient(region *string, accountID *string, role *string) SNSAPI
	SFNClient(region *string, accountID *string, role *string) SFNAPI
}

// ClientsStr implementation
type ClientsStr struct {
	session *session.Session
	configs map[string]*aws.Config

	S3  S3API
	ASG ASGAPI
	ELB ELBAPI
	EC2 EC2API
	ALB ALBAPI
	CW  CWAPI
	IAM IAMAPI
	SNS SNSAPI
}

// GetSession get session
func (awsc *ClientsStr) GetSession() *session.Session {
	return awsc.session
}

// SetSession assings session
func (awsc *ClientsStr) SetSession(sess *session.Session) {
	awsc.session = sess
}

// GetConfig retrieves config for key
func (awsc *ClientsStr) GetConfig(key string) *aws.Config {
	if awsc.configs == nil {
		return nil
	}

	config, ok := awsc.configs[key]
	if ok && config != nil {
		return config
	}

	return nil
}

// SetConfig assigns config for key
func (awsc *ClientsStr) SetConfig(key string, config *aws.Config) {
	if awsc.configs == nil {
		awsc.configs = map[string]*aws.Config{}
	}
	awsc.configs[key] = config
}

// S3Client returns client for region account and role
func (awsc *ClientsStr) S3Client(region *string, accountID *string, role *string) S3API {
	return s3.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// ASGClient returns client for region account and role
func (awsc *ClientsStr) ASGClient(region *string, accountID *string, role *string) ASGAPI {
	return autoscaling.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// ELBClient returns client for region account and role
func (awsc *ClientsStr) ELBClient(region *string, accountID *string, role *string) ELBAPI {
	return elb.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// EC2Client returns client for region account and role
func (awsc *ClientsStr) EC2Client(region *string, accountID *string, role *string) EC2API {
	return ec2.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// ALBClient returns client for region account and role
func (awsc *ClientsStr) ALBClient(region *string, accountID *string, role *string) ALBAPI {
	return elbv2.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// CWClient returns client for region account and role
func (awsc *ClientsStr) CWClient(region *string, accountID *string, role *string) CWAPI {
	return cloudwatch.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// IAMClient returns client for region account and role
func (awsc *ClientsStr) IAMClient(region *string, accountID *string, role *string) IAMAPI {
	return iam.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// SNSClient returns client for region account and role
func (awsc *ClientsStr) SNSClient(region *string, accountID *string, role *string) SNSAPI {
	return sns.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}

// SFNClient returns client for region account and role
func (awsc *ClientsStr) SFNClient(region *string, accountID *string, role *string) SFNAPI {
	return sfn.New(ar.Session(awsc), ar.Config(awsc, region, accountID, role))
}
