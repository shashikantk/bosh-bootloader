package ec2

import (
	"fmt"
	"strings"

	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/genevieve/leftovers/aws/common"
)

type securityGroupsClient interface {
	DescribeSecurityGroups(*awsec2.DescribeSecurityGroupsInput) (*awsec2.DescribeSecurityGroupsOutput, error)
	RevokeSecurityGroupIngress(*awsec2.RevokeSecurityGroupIngressInput) (*awsec2.RevokeSecurityGroupIngressOutput, error)
	RevokeSecurityGroupEgress(*awsec2.RevokeSecurityGroupEgressInput) (*awsec2.RevokeSecurityGroupEgressOutput, error)
	DeleteSecurityGroup(*awsec2.DeleteSecurityGroupInput) (*awsec2.DeleteSecurityGroupOutput, error)
}

type SecurityGroups struct {
	client securityGroupsClient
	logger logger
}

func NewSecurityGroups(client securityGroupsClient, logger logger) SecurityGroups {
	return SecurityGroups{
		client: client,
		logger: logger,
	}
}

func (e SecurityGroups) List(filter string) ([]common.Deletable, error) {
	output, err := e.client.DescribeSecurityGroups(&awsec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return nil, fmt.Errorf("Describing security groups: %s", err)
	}

	var resources []common.Deletable
	for _, sg := range output.SecurityGroups {
		resource := NewSecurityGroup(e.client, sg.GroupId, sg.GroupName, sg.Tags, sg.IpPermissions, sg.IpPermissionsEgress)

		if *sg.GroupName == "default" {
			continue
		}

		if !strings.Contains(resource.identifier, filter) {
			continue
		}

		proceed := e.logger.Prompt(fmt.Sprintf("Are you sure you want to delete security group %s?", resource.identifier))
		if !proceed {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
