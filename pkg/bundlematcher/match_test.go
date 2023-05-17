package bundlematcher

import (
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	//	"github.com/rancher/fleet/pkg/match"
	//        "github.com/rancher/fleet/pkg/fleetyaml"
	"fmt"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
)

type fleetYAML struct {
	Name   string            `json:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	fleet.BundleSpec
	TargetCustomizations []fleet.BundleTarget `json:"targetCustomizations,omitempty"`
	ImageScans           []fleet.ImageScan    `json:"imageScans,omitempty"`
}

func TestBundleMatchCustomization(t *testing.T) {
	bundleSpecReader, _ := os.Open("./test.yaml")
	bytes, err := io.ReadAll(bundleSpecReader)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fy := &fleetYAML{}
	if err := yaml.Unmarshal(bytes, fy); err != nil {
		fmt.Printf("yaml unmarshal error: %s", err)
	}
	fy.BundleSpec.Targets = append(fy.BundleSpec.Targets, fy.TargetCustomizations...)

	bundle := &fleet.Bundle{
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       fy.BundleSpec,
	}

	bm, err := New(bundle)
        fmt.Printf("%+v\n", bm)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	clusterName := "cl-test"
	cg := map[string]map[string]string{
		"apps": {"a": "true", "b": "true"},
		"features": {"f1": "10", "f2": "10"},
	}
	labels := map[string]string{"a": "true",
                                    "b": "true",
                                    "match": "true",
                                    "append": "true"}
        bundleTarget := bm.MatchTargetCustomizations(clusterName, cg, labels)
	fmt.Printf("Result helm values %+v", bundleTarget.Helm.Values)
}
