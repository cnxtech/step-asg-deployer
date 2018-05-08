package client

import (
	"fmt"

	"github.com/coinbase/step-asg-deployer/aws"
	"github.com/coinbase/step-asg-deployer/deployer/models"
	"github.com/coinbase/step/execution"
	"github.com/coinbase/step/utils/to"
)

// Halt attempts to halt release
func Halt(fileOrJSON *string) error {
	region, accountID := to.RegionAccount()
	release, err := releaseFromFileOrJSON(fileOrJSON, region, accountID)
	if err != nil {
		return err
	}

	deployerARN := to.StepArn(region, accountID, to.Strp("coinbase-step-asg-deployer"))

	return halt(&aws.ClientsStr{}, release, deployerARN)
}

func halt(awsc aws.Clients, release *models.Release, deployerARN *string) error {
	exec, err := execution.FindExecution(awsc.SFNClient(nil, nil, nil), deployerARN, executionPrefix(release))
	if err != nil {
		return err
	}

	if exec == nil {
		return fmt.Errorf("Cannot find current execution of release with prefix %q", executionPrefix(release))
	}

	if err := release.Halt(awsc.S3Client(nil, nil, nil)); err != nil {
		return err
	}

	exec.WaitForExecution(awsc.SFNClient(nil, nil, nil), 1, waiter)
	fmt.Println("")
	return nil
}
