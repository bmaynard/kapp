package app

import (
	"fmt"
	"strings"

	ctlres "github.com/k14s/kapp/pkg/kapp/resources"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type LabeledApp struct {
	labelSelector labels.Selector

	coreClient    kubernetes.Interface
	dynamicClient dynamic.Interface
}

var _ App = &LabeledApp{}

func (a *LabeledApp) Name() string {
	str := a.labelSelector.String()
	if len(str) == 0 {
		return "?"
	}
	return str
}

func (a *LabeledApp) Namespace() string { return "" }

func (a *LabeledApp) LabelSelector() (labels.Selector, error) {
	return a.labelSelector, nil
}

func (a *LabeledApp) CreateOrUpdate(labels map[string]string) error { return nil }
func (a *LabeledApp) Exists() (bool, error)                         { return true, nil }

func (a *LabeledApp) Delete() error {
	labelSelector, err := a.LabelSelector()
	if err != nil {
		return err
	}

	rs, err := ctlres.NewIdentifiedResources(a.coreClient, a.dynamicClient).List(labelSelector)
	if err != nil {
		return fmt.Errorf("Relisting app resources: %s", err)
	}

	if len(rs) > 0 {
		var resourceNames []string
		for _, res := range rs {
			resourceNames = append(resourceNames, res.Description())
		}
		return fmt.Errorf("Expected all resources to be gone, but found: %s", strings.Join(resourceNames, ", "))
	}

	return nil
}

func (a *LabeledApp) Meta() (AppMeta, error) { return AppMeta{}, nil }

func (a *LabeledApp) Changes() ([]Change, error)             { return nil, nil }
func (a *LabeledApp) LastChange() (Change, error)            { return nil, nil }
func (a *LabeledApp) BeginChange(ChangeMeta) (Change, error) { return NoopChange{}, nil }
