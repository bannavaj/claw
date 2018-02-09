package claw

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/opsworks"
)

// OpsWorksHandler : OpsWorks struct
type OpsWorksHandler struct {
	conn *opsworks.OpsWorks
}

// InstanceIDMap : Map of instance ids with statuses
type InstanceIDMap struct {
	Hostname        *string
	InstanceID      *string
	Status          *string
	AutoScalingType *string
}

// GetAllStacks : Get the number of instances on a cluster
func (client *OpsWorksHandler) GetAllStacks() ([]*opsworks.Stack, error) {
	stackInput := &opsworks.DescribeStacksInput{}
	stacks, err := client.conn.DescribeStacks(stackInput)
	if err != nil {
		return nil, err
	}

	return stacks.Stacks, nil
}

// GetStackSummary : get stack summary
func (client *OpsWorksHandler) GetStackSummary(stackID *string) (*opsworks.DescribeStackSummaryOutput, error) {
	stackSummaryInput := &opsworks.DescribeStackSummaryInput{
		StackId: stackID,
	}
	stackSummary, err := client.conn.DescribeStackSummary(stackSummaryInput)
	if err != nil {
		return nil, err
	}

	return stackSummary, nil
}

// GetStack : Get a stack by id
func (client *OpsWorksHandler) GetStack(stackID *string) (*opsworks.Stack, error) {
	stackInput := &opsworks.DescribeStacksInput{
		StackIds: []*string{stackID},
	}
	stacks, err := client.conn.DescribeStacks(stackInput)
	if err != nil {
		return nil, err
	}

	return stacks.Stacks[0], nil
}

// GetStackID : Get the stack id
func (client *OpsWorksHandler) GetStackID(clusterName *string) (string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return "", errors.New("Error getting all stacks")
	}

	for _, stack := range stacks {
		if *stack.Name == *clusterName {
			return *stack.StackId, nil
		}
	}

	return "", errors.New("Invalid cluster name")
}

// GetStackNames : Get all stack names
func (client *OpsWorksHandler) GetStackNames() ([]string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return []string{""}, errors.New("Error getting all stacks")
	}

	stackNames := make([]string, 0)
	for _, stack := range stacks {
		stackNames = append(stackNames, *stack.Name)
	}

	return stackNames, nil
}

// GetRegionStackNames : Get all stack names in a region
func (client *OpsWorksHandler) GetRegionStackNames(region *string) ([]string, error) {
	stacks, err := client.GetAllStacks()
	if err != nil {
		return []string{""}, errors.New("Error getting all stacks")
	}

	stackNames := make([]string, 0)
	for _, stack := range stacks {
		if *stack.Region == *region {
			stackNames = append(stackNames, *stack.Name)
		}
	}

	return stackNames, nil
}

// GetAllApps : Get all apps on a stack
func (client *OpsWorksHandler) GetAllApps(stackID *string) ([]*opsworks.App, error) {
	stackIDInput := &opsworks.DescribeAppsInput{
		StackId: stackID,
	}
	apps, err := client.conn.DescribeApps(stackIDInput)
	if err != nil {
		return []*opsworks.App{}, errors.New("Error getting all apps")
	}

	return apps.Apps, nil
}

// GetAppID : Get app id in a stack
func (client *OpsWorksHandler) GetAppID(stackID *string, appName *string) (string, error) {
	apps, err := client.GetAllApps(stackID)
	if err != nil {
		return "", errors.New("Error getting all apps")
	}

	for _, app := range apps {
		if *app.Name == *appName {
			return *app.AppId, nil
		}
	}

	return "", errors.New("Invalid stack id or app name")
}

// GetAllLayers : Get all layers in a stack
func (client *OpsWorksHandler) GetAllLayers(stackID *string) ([]*opsworks.Layer, error) {
	stackIDInput := &opsworks.DescribeLayersInput{
		StackId: stackID,
	}
	layers, err := client.conn.DescribeLayers(stackIDInput)
	if err != nil {
		return []*opsworks.Layer{}, errors.New("Error getting all layers")
	}

	return layers.Layers, nil
}

// GetLayerID : Get layer id in a stack
func (client *OpsWorksHandler) GetLayerID(stackID *string, appName *string) (string, error) {
	layers, err := client.GetAllLayers(stackID)
	if err != nil {
		return "", errors.New("Error getting all layers")
	}

	for _, layer := range layers {
		if *layer.Name == *appName {
			return *layer.LayerId, nil
		}
	}

	return "", errors.New("Invalid stack id or app name")
}

// GetAllInstances : Get all instances in a layer
func (client *OpsWorksHandler) GetAllInstances(layerID *string) ([]*opsworks.Instance, error) {
	layerIDInput := &opsworks.DescribeInstancesInput{
		LayerId: layerID,
	}
	instances, err := client.conn.DescribeInstances(layerIDInput)
	if err != nil {
		return []*opsworks.Instance{}, errors.New("Error getting all instances")
	}

	return instances.Instances, nil
}

// GetInstanceIDs : Get all instance ids in a layer
func (client *OpsWorksHandler) GetInstanceIDs(layerID *string) ([]*InstanceIDMap, error) {
	instances, err := client.GetAllInstances(layerID)
	if err != nil {
		return []*InstanceIDMap{}, errors.New("Error getting all instances")
	}

	instanceList := make([]*InstanceIDMap, 0)
	for _, instance := range instances {
		instanceMap := &InstanceIDMap{
			instance.Hostname,
			instance.InstanceId,
			instance.Status,
			instance.AutoScalingType,
		}

		instanceList = append(instanceList, instanceMap)
	}

	return instanceList, nil
}

// NewOpsworksClient : Create a new opsworks client
func NewOpsworksClient() (*OpsWorksHandler, error) {
	sess, err := newSession()
	if err != nil {
		return &OpsWorksHandler{}, err
	}

	return &OpsWorksHandler{conn: opsworks.New(sess)}, nil
}
