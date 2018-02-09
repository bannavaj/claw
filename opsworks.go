package claw

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

// OpsWorksHandler : OpsWorks struct
type OpsWorksHandler struct {
	conn *opsworks.OpsWorks
}

// GetAllStacks : Get the number of instances on a cluster
func (client *OpsWorksHandler) GetAllStacks() ([]*opsworks.Stack, error) {
	params := &opsworks.DescribeStacksInput{}
	stacks, err := client.conn.DescribeStacks(params)
	if err != nil {
		return nil, err
	}

	return stacks.Stacks, nil
}

// GetStackSummary : get stack summary
func (client *OpsWorksHandler) GetStackSummary(stackID string) (*opsworks.DescribeStackSummaryOutput, error) {
	params := &opsworks.DescribeStackSummaryInput{
		StackId: aws.String(stackID),
	}
	stackSummary, err := client.conn.DescribeStackSummary(params)
	if err != nil {
		return nil, err
	}

	return stackSummary, nil
}

// GetStack : Get a stack by id
func (client *OpsWorksHandler) GetStack(stackID string) (*opsworks.Stack, error) {
	params := &opsworks.DescribeStacksInput{
		StackIds: []*string{aws.String(stackID)},
	}
	stacks, err := client.conn.DescribeStacks(params)
	if err != nil {
		return nil, err
	}

	return stacks.Stacks[0], nil
}

// GetStackID : Get the stack id
func (client *OpsWorksHandler) GetStackID(clusterName string) (string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return "", err
	}

	for _, stack := range stacks {
		if *stack.Name == clusterName {
			return *stack.StackId, nil
		}
	}

	return "", nil
}

// GetStackNames : Get all stack names
func (client *OpsWorksHandler) GetStackNames() ([]string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return []string{""}, err
	}

	stackNames := make([]string, 0)
	for _, stack := range stacks {
		stackNames = append(stackNames, *stack.Name)
	}

	return stackNames, nil
}

// GetRegionStackNames : Get all stack names in a region
func (client *OpsWorksHandler) GetRegionStackNames(region string) ([]string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return []string{""}, err
	}

	stackNames := make([]string, 0)
	for _, stack := range stacks {
		if *stack.Region == region {
			stackNames = append(stackNames, *stack.Name)
		}
	}

	return stackNames, nil
}

// GetAllApps : Get all apps on a stack
func (client *OpsWorksHandler) GetAllApps(stackID string) ([]*opsworks.App, error) {
	params := &opsworks.DescribeAppsInput{
		StackId: aws.String(stackID),
	}
	apps, err := client.conn.DescribeApps(params)
	if err != nil {
		return []*opsworks.App{}, err
	}

	return apps.Apps, nil
}

// GetAppID : Get app id in a stack
func (client *OpsWorksHandler) GetAppID(stackID string, appName string) (string, error) {
	apps, err := client.GetAllApps(stackID)
	if err != nil {
		return "", err
	}

	for _, app := range apps {
		if *app.Name == appName {
			return *app.AppId, nil
		}
	}

	return "", nil
}

// GetAllLayers : Get all layers in a stack
func (client *OpsWorksHandler) GetAllLayers(stackID string) ([]*opsworks.Layer, error) {
	params := &opsworks.DescribeLayersInput{
		StackId: aws.String(stackID),
	}
	layers, err := client.conn.DescribeLayers(params)
	if err != nil {
		return []*opsworks.Layer{}, err
	}

	return layers.Layers, nil
}

// GetLayerID : Get layer id in a stack
func (client *OpsWorksHandler) GetLayerID(stackID string, appName string) (string, error) {
	layers, err := client.GetAllLayers(stackID)
	if err != nil {
		return "", err
	}

	for _, layer := range layers {
		if *layer.Name == appName {
			return *layer.LayerId, nil
		}
	}

	return "", nil
}

// NewOpsworksClient : Create a new opsworks client
func NewOpsworksClient() (*OpsWorksHandler, error) {
	sess, err := newSession()
	if err != nil {
		return &OpsWorksHandler{}, err
	}

	return &OpsWorksHandler{conn: opsworks.New(sess)}, nil
}
