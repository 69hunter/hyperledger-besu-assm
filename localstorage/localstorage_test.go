package localstorage

import (
	"reflect"
	"testing"

	"github.com/69hunter/hyperledger-besu-assm/core"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestCreateBootnodes(t *testing.T) {
	a := NewAdapter()

	for i := 0; i < 4; i++ {
		a.AddNodeInfo(core.NodeInfo{
			ValidatorAddress: "0xa8bb6f2548e9b2fbc6b08e98e8ba21a2b9acfe42",
			NodePublicKey:    "0xeddf5adbd8ca23df4adf734ec1f4aa564f8be782e602ba6a26ff3bc775ce4022",
			NodeHost:         "10.0.0.0",
			NodePort:         30303,
		})
	}
	result, _ := a.CreateBootnodes()
	assert(t, result, []string{
		"enode://eddf5adbd8ca23df4adf734ec1f4aa564f8be782e602ba6a26ff3bc775ce4022@10.0.0.0:30303",
		"enode://eddf5adbd8ca23df4adf734ec1f4aa564f8be782e602ba6a26ff3bc775ce4022@10.0.0.0:30303",
		"enode://eddf5adbd8ca23df4adf734ec1f4aa564f8be782e602ba6a26ff3bc775ce4022@10.0.0.0:30303",
		"enode://eddf5adbd8ca23df4adf734ec1f4aa564f8be782e602ba6a26ff3bc775ce4022@10.0.0.0:30303",
	})
}
