package dodas

import (
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config"
)

const (
	clusterStatusUpdateInProgress = "UPDATE_IN_PROGRESS"
	clusterStatusUpdateComplete   = "UPDATE_COMPLETE"
	clusterStatusUpdateFailed     = "UPDATE_FAILED"

	waitForStatusTimeStep       = 30 * time.Second
	waitForUpdateStatusTimeout  = 2 * time.Minute
	waitForCompleteStatusTimout = 10 * time.Minute

	// Could move to property of magnumManager implementations if needed
	scaleToZeroSupported = false

	// Time that the goroutine that first acquires clusterUpdateMutex
	// in deleteNodes should wait for other synchronous calls to deleteNodes.
	deleteNodesBatchingDelay = 2 * time.Second
)

const (
	stackStatusUpdateInProgress = "UPDATE_IN_PROGRESS"
	stackStatusUpdateComplete   = "UPDATE_COMPLETE"
)

// statusesPreventingUpdate is a set of statuses that would prevent
// the cluster from successfully scaling.
//
// TODO: If it becomes possible to update even in UPDATE_FAILED state then it can be removed here
// https://storyboard.openstack.org/#!/story/2005056
var statusesPreventingUpdate = sets.NewString(
	clusterStatusUpdateInProgress,
	clusterStatusUpdateFailed,
)

// magnumManagerHeat implements the magnumManager interface.
//
// Most interactions with the cluster are done directly with magnum,
// but scaling down requires an intermediate step using heat to
// delete the specific nodes that the autoscaler has picked for removal.
type dodasManager struct {
	clusterClient *dodasclient.ServiceClient
	clusterName   string

	kubeMinionsStackName string

	waitTimeStep time.Duration
}

// createMagnumManagerHeat sets up cluster and stack clients and returns
// an magnumManagerHeat.
func createDodasManager(configReader io.Reader, discoverOpts cloudprovider.NodeGroupDiscoveryOptions, opts config.AutoscalingOptions) (*dodasManager, error) {

	clusterClient = dodasclient{}

	manager := dodasManager{
		clusterClient: clusterClient,
		clusterName:   opts.ClusterName,
		waitTimeStep:  waitForStatusTimeStep,
	}

	return &manager, nil
}

func (mgr *dodasManager) updateNodeCount(nodegroup string, nodes int) error {
	// TODO: scale IM

	// updateOpts := []clusters.UpdateOptsBuilder{
	// 	UpdateOptsInt{Op: clusters.ReplaceOp, Path: "/node_count", Value: nodes},
	// }
	// _, err := clusters.Update(mgr.clusterClient, mgr.clusterName, updateOpts).Extract()
	// if err != nil {
	// 	return fmt.Errorf("could not update cluster: %v", err)
	// }
	// return nil
	return nil
}

func (mgr *dodasManager) getNodes(nodegroup string) ([]string, error) {
	// TODO: get node ProviderIDs by getting nova instance IDs from heat
	// Waiting for https://review.openstack.org/#/c/639053/ to be able to get
	// nova instance IDs from the kube_minions stack resource.
	// This works fine being empty for now anyway.
	return []string{}, nil
}

// NodeRef ...
type NodeRef struct {
	Name       string
	MachineID  string
	ProviderID string
	IPs        []string
}

// deleteNodes deletes nodes by passing a comma separated list of names or IPs
// of minions to remove to heat, and simultaneously sets the new number of minions on the stack.
// The magnum node_count is then set to the new value (does not cause any more nodes to be removed).
//
// TODO: The two step process is required until https://storyboard.openstack.org/#!/story/2005052
// is complete, which will allow resizing with specific nodes to be deleted as a single Magnum operation.
func (mgr *dodasManager) deleteNodes(nodegroup string, nodes []NodeRef, updatedNodeCount int) error {
	// delete node

	return nil
}

// getClusterStatus returns the current status of the magnum cluster.
func (mgr *dodasManager) getClusterStatus() (string, error) {

	//return cluster.Status, nil
	return "", nil
}

// canUpdate checks if the cluster status is present in a set of statuses that
// prevent the cluster from being updated.
// Returns if updating is possible and the status for convenience.
func (mgr *dodasManager) canUpdate() (bool, string, error) {
	// clusterStatus, err := mgr.getClusterStatus()
	// if err != nil {
	// 	return false, "", fmt.Errorf("could not get cluster status: %v", err)
	// }
	// return !statusesPreventingUpdate.Has(clusterStatus), clusterStatus, nil
	return true, "", nil
}

// templateNodeInfo returns a NodeInfo with a node template based on the VM flavor
// that is used to created minions in a given node group.
func (mgr *dodasManager) templateNodeInfo(nodegroup string) (*schedulernodeinfo.NodeInfo, error) {
	// TODO: create a node template by getting the minion flavor from the heat stack.
	return nil, cloudprovider.ErrNotImplemented
}

// waitForStackStatus checks periodically to see if the heat stack has entered a given status.
// Returns when the status is observed or the timeout is reached.
func (mgr *dodasManager) waitForStackStatus(status string, timeout time.Duration) error {
	// klog.V(2).Infof("Waiting for stack %s status", status)
	// for start := time.Now(); time.Since(start) < timeout; time.Sleep(mgr.waitTimeStep) {
	// 	currentStatus, err := mgr.getStackStatus()
	// 	if err != nil {
	// 		return fmt.Errorf("error waiting for stack status: %v", err)
	// 	}
	// 	if currentStatus == status {
	// 		klog.V(0).Infof("Waited for stack %s status", status)
	// 		return nil
	// 	}
	// }
	return fmt.Errorf("timeout (%v) waiting for stack status %s", timeout, status)
}
