package alb

import (
	"sort"
	"testing"

	"github.com/coinbase/step-asg-deployer/aws/mocks"
	"github.com/coinbase/step/utils/to"
	"github.com/stretchr/testify/assert"
)

func Test_FindAll_Empty(t *testing.T) {
	albc := &mocks.ALBClient{}
	am, err := FindAll(albc, []*string{})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(am))
}

func Test_FindAll_NotFound(t *testing.T) {
	albc := &mocks.ALBClient{}
	_, err := FindAll(albc, []*string{to.Strp("tg_name")})
	assert.Error(t, err)
}

func Test_FindAll_Found(t *testing.T) {
	albc := &mocks.ALBClient{}
	albc.AddTargetGroup("tg_name", "project_name", "config_name", "service_name")
	albc.AddTargetGroup("tg_other_name", "project_name", "config_name", "service_name")
	am, err := FindAll(albc, []*string{to.Strp("tg_name"), to.Strp("tg_other_name")})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(am))
	assert.Equal(t, *am[0].TargetGroupArn, "tg_name")
	assert.Equal(t, *am[1].TargetGroupArn, "tg_other_name")
}

func Test_GetInstances(t *testing.T) {
	albc := &mocks.ALBClient{}
	albc.AddTargetGroup("tg_name", "project_name", "config_name", "service_name")

	instances, err := GetInstances(albc, to.Strp("tg_name"), []string{"InstanceId"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(instances))
}

func Test_createDescribeTargetHealthInput(t *testing.T) {
	name := ""

	in := createDescribeTargetHealthInput(&name, []string{})
	assert.Equal(t, len(in.Targets), 0)

	in = createDescribeTargetHealthInput(&name, []string{"a"})
	assert.Equal(t, len(in.Targets), 1)
	assert.Equal(t, *in.Targets[0].Id, "a")

	in = createDescribeTargetHealthInput(&name, []string{"a", "b"})
	tgsIDs := []string{}
	for _, tg := range in.Targets {
		tgsIDs = append(tgsIDs, *tg.Id)
	}

	sort.Strings(tgsIDs) // Sort not Guaranteed by map

	assert.Equal(t, len(tgsIDs), 2)
	assert.Equal(t, tgsIDs[0], "a")
	assert.Equal(t, tgsIDs[1], "b")
}
