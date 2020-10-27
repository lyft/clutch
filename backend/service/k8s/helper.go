package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// GenerateStrategicPatch will return a patch that yields
// the modified object when applied to the original object,
// or an error if either of the two objects is invalid.
func GenerateStrategicPatch(oldBytes, newBytes []byte) ([]byte, error) {
	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldBytes, newBytes, appsv1.Deployment{})
	if err != nil {
		return nil, err
	}
	return patchBytes, nil
}
