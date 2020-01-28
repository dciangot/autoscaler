package dodas

import (
	"fmt"
	"sync"
	"time"

	apiv1 "k8s.io/api/core/v1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type dodasNodeGroup struct {
	dodasManager dodasManager
	id           string

	clusterUpdateMutex *sync.Mutex

	minSize int
	maxSize int
	// Stored as a pointer so that when autoscaler copies the nodegroup it can still update the target size
	targetSize *int

	// TODO --> use dodas CRD
	nodesToDelete      []*apiv1.Node
	nodesToDeleteMutex sync.Mutex

	waitTimeStep time.Duration
}

// waitForClusterStatus checks periodically to see if the cluster has entered a given status.
// Returns when the status is observed or the timeout is reached.
func (ng *dodasNodeGroup) waitForClusterStatus(status string, timeout time.Duration) error {
	// klog.V(2).Infof("Waiting for cluster %s status", status)
	// for start := time.Now(); time.Since(start) < timeout; time.Sleep(ng.waitTimeStep) {
	// 	clusterStatus, err := ng.magnumManager.getClusterStatus()
	// 	if err != nil {
	// 		return fmt.Errorf("error waiting for %s status: %v", status, err)
	// 	}
	// 	if clusterStatus == status {
	// 		klog.V(0).Infof("Waited for cluster %s status", status)
	// 		return nil
	// 	}
	// }
	return fmt.Errorf("timeout (%v) waiting for %s status", timeout, status)
}

// IncreaseSize increases the number of nodes by replacing the cluster's node_count.
//
// Takes precautions so that the cluster is not modified while in an UPDATE_IN_PROGRESS state.
// Blocks until the cluster has reached UPDATE_COMPLETE.
func (ng *dodasNodeGroup) IncreaseSize(delta int) error {

	return nil
}

// deleteNodes deletes a set of nodes chosen by the autoscaler.
//
// The process of deletion depends on the implementation of magnumManager,
// but this function handles what should be common between all implementations:
//   - simultaneous but separate calls from the autoscaler are batched together
//   - does not allow scaling while the cluster is already in an UPDATE_IN_PROGRESS state
//   - after scaling down, blocks until the cluster has reached UPDATE_COMPLETE
func (ng *dodasNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	return nil
}

// DecreaseTargetSize decreases the cluster node_count in magnum.
func (ng *dodasNodeGroup) DecreaseTargetSize(delta int) error {
	// 	if delta >= 0 {
	// 		return fmt.Errorf("size decrease must be negative")
	// 	}
	// 	klog.V(0).Infof("Decreasing target size by %d, %d->%d", delta, *ng.targetSize, *ng.targetSize+delta)
	// 	*ng.targetSize += delta
	// 	return ng.magnumManager.updateNodeCount(ng.id, *ng.targetSize)
	return nil
}

// Id returns the node group ID
func (ng *dodasNodeGroup) Id() string {
	return ng.id
}

// Debug returns a string formatted with the node group's min, max and target sizes.
func (ng *dodasNodeGroup) Debug() string {
	return fmt.Sprintf("%s min=%d max=%d target=%d", ng.id, ng.minSize, ng.maxSize, *ng.targetSize)
}

// Nodes returns a list of nodes that belong to this node group.
func (ng *dodasNodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	// nodes, err := ng.magnumManager.getNodes(ng.id)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get nodes: %v", err)
	// }
	// var instances []cloudprovider.Instance
	// for _, node := range nodes {
	// 	instances = append(instances, cloudprovider.Instance{Id: node})
	// }
	//return instances, nil
	return []cloudprovider.Instance{}, nil
}

// TemplateNodeInfo returns a node template for this node group.
func (ng *dodasNodeGroup) TemplateNodeInfo() (*schedulernodeinfo.NodeInfo, error) {
	return ng.dodasManager.templateNodeInfo(ng.id)
}

// Exist returns if this node group exists.
// Currently always returns true.
func (ng *dodasNodeGroup) Exist() bool {
	return true
}

// Create creates the node group on the cloud provider side.
func (ng *dodasNodeGroup) Create() (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrAlreadyExist
}

// Delete deletes the node group on the cloud provider side.
func (ng *dodasNodeGroup) Delete() error {
	return cloudprovider.ErrNotImplemented
}

// Autoprovisioned returns if the nodegroup is autoprovisioned.
func (ng *dodasNodeGroup) Autoprovisioned() bool {
	return false
}

// MaxSize returns the maximum allowed size of the node group.
func (ng *dodasNodeGroup) MaxSize() int {
	return ng.maxSize
}

// MinSize returns the minimum allowed size of the node group.
func (ng *dodasNodeGroup) MinSize() int {
	return ng.minSize
}

// TargetSize returns the target size of the node group.
func (ng *dodasNodeGroup) TargetSize() (int, error) {
	return *ng.targetSize, nil
}
