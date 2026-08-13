package extract

import (
	"reflect"
	"testing"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSelectorExtractsReplicationControllerSelector(t *testing.T) {
	rc := &coreV1.ReplicationController{
		Spec: coreV1.ReplicationControllerSpec{
			Selector: map[string]string{
				"app": "web",
			},
		},
	}

	selector, found := Selector(rc)
	if !found {
		t.Fatal("expected selector to be found")
	}

	expected := &metaV1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "web",
		},
	}
	if !reflect.DeepEqual(selector, expected) {
		t.Fatalf("expected selector %#v, got %#v", expected, selector)
	}
}
