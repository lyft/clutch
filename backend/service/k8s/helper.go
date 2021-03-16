package k8s

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// GenerateStrategicPatch will return a patch that yields
// the modified object when applied to the original object,
// or an error if either of the two objects is invalid.
func GenerateStrategicPatch(original, modified runtime.Object, dataStruct interface{}) ([]byte, error) {
	oldBytes, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	newBytes, err := json.Marshal(modified)
	if err != nil {
		return nil, err
	}

	return strategicpatch.CreateTwoWayMergePatch(oldBytes, newBytes, dataStruct)
}
