package dodas

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
)

const (
	// GPULabel is the label added to nodes with GPU resource.
	GPULabel = "cloud.google.com/gke-accelerator"
)

var (
	availableGPUTypes = map[string]struct{}{
		"nvidia-tesla-k80":  {},
		"nvidia-tesla-p100": {},
		"nvidia-tesla-v100": {},
	}
)

type dodasCloudProvider struct {
	magnumManager   *dodasManager
	resourceLimiter *cloudprovider.ResourceLimiter
	nodeGroups      []dodasNodeGroup
}

func buildDodasloudProvider(dodasManager dodasManager, resourceLimiter *cloudprovider.ResourceLimiter) (cloudprovider.CloudProvider, error) {
	dcp := &dodasCloudProvider{
		magnumManager:   &dodasManager,
		resourceLimiter: resourceLimiter,
		nodeGroups:      []dodasNodeGroup{},
	}
	return dcp, nil
}

// Name returns the name of the cloud provider.
func (dcp *dodasCloudProvider) Name() string {
	return cloudprovider.DodasProviderName
}

// GPULabel returns the label added to nodes with GPU resource.
func (dcp *dodasCloudProvider) GPULabel() string {
	return GPULabel
}

// GetAvailableGPUTypes return all available GPU types cloud provider supports
func (dcp *dodasCloudProvider) GetAvailableGPUTypes() map[string]struct{} {
	return availableGPUTypes
}

// NodeGroups returns all node groups managed by this cloud provider.
func (dcp *dodasCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	groups := make([]cloudprovider.NodeGroup, len(dcp.nodeGroups))
	for i, group := range dcp.nodeGroups {
		groups[i] = &group
	}
	return groups
}

// AddNodeGroup appends a node group to the list of node groups managed by this cloud provider.
func (dcp *dodasCloudProvider) AddNodeGroup(group dodasNodeGroup) {
	dcp.nodeGroups = append(dcp.nodeGroups, group)
}

// NodeGroupForNode returns the node group that a given node belongs to.
//
// Since only a single node group is currently supported, the first node group is always returned.
func (dcp *dodasCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	// TODO: wait for magnum nodegroup support
	if _, found := node.ObjectMeta.Labels["node-role.kubernetes.io/master"]; found {
		return nil, nil
	}
	return &(dcp.nodeGroups[0]), nil
}

// Pricing is not implemented.
func (dcp *dodasCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetAvailableMachineTypes is not implemented.
func (dcp *dodasCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	return []string{}, nil
}

// NewNodeGroup is not implemented.
func (dcp *dodasCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string,
	taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetResourceLimiter returns resource constraints for the cloud provider
func (dcp *dodasCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	return dcp.resourceLimiter, nil
}

// GetInstanceID gets the instance ID for the specified node.
func (dcp *dodasCloudProvider) GetInstanceID(node *apiv1.Node) string {
	return node.Spec.ProviderID
}

// Refresh is called before every autoscaler main loop.
//
// Currently only prints debug information.
func (dcp *dodasCloudProvider) Refresh() error {
	for _, nodegroup := range dcp.nodeGroups {
		klog.V(3).Info(nodegroup.Debug())
	}
	return nil
}

// Cleanup currently does nothing.
func (dcp *dodasCloudProvider) Cleanup() error {
	return nil
}
