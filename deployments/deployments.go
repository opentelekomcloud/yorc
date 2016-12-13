package deployments

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"novaforge.bull.com/starlings-janus/janus/helper/consulutil"
	"path"
	"strings"
)

type deploymentNotFound struct {
	deploymentId string
}

func (d deploymentNotFound) Error() string {
	return fmt.Sprintf("Deployment with id %q not found", d.deploymentId)
}

// IsDeploymentNotFoundError checks if an error is a deployment not found error
func IsDeploymentNotFoundError(err error) bool {
	_, ok := err.(deploymentNotFound)
	return ok
}

// DeploymentStatusFromString returns a DeploymentStatus from its textual representation.
//
// If ignoreCase is 'true' the given status is uppercased to match the generated status strings.
// If the given status does not match any known status an error is returned
func DeploymentStatusFromString(status string, ignoreCase bool) (DeploymentStatus, error) {
	if ignoreCase {
		status = strings.ToUpper(status)
	}
	for i := startOfDepStatusConst + 1; i < endOfDepStatusConst; i++ {
		if DeploymentStatus(i).String() == status {
			return DeploymentStatus(i), nil
		}
	}
	return INITIAL, fmt.Errorf("Invalid deployment status %q", status)
}

// GetDeploymentStatus returns a DeploymentStatus for a given deploymentId
//
// If the given deploymentId doesn't refer to an existing deployment an error is returned. This error could be checked with
//  IsDeploymentNotFoundError(err error) bool
//
// For example:
//  if status, err := GetDeploymentStatus(kv, deploymentId); err != nil {
//  	if IsDeploymentNotFoundError(err) {
//  		// Do something in case of deployment not found
//  	}
//  }
func GetDeploymentStatus(kv *api.KV, deploymentId string) (DeploymentStatus, error) {
	kvp, _, err := kv.Get(path.Join(consulutil.DeploymentKVPrefix, deploymentId, "status"), nil)
	if err != nil {
		return INITIAL, err
	}
	if kvp == nil || len(kvp.Value) == 0 {
		return INITIAL, deploymentNotFound{deploymentId: deploymentId}
	}
	return DeploymentStatusFromString(string(kvp.Value), true)
}

// DoesDeploymentExists checks if a given deploymentId refer to an existing deployment
func DoesDeploymentExists(kv *api.KV, deploymentId string) (bool, error) {
	if _, err := GetDeploymentStatus(kv, deploymentId); err != nil {
		if IsDeploymentNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
