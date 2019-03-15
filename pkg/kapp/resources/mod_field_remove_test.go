package resources_test

import (
	"testing"

	ctlres "github.com/k14s/kapp/pkg/kapp/resources"
)

func TestModFieldRemove(t *testing.T) {
	exs := []modFieldRemoveExample{
		{
			Description: "deleting leaf key that exists",
			Res: `
metadata:
  labels:
    label-key: label-val`,
			Expected: `
metadata:
  labels: {}`,
			Path: ctlres.NewPathFromStrings([]string{"metadata", "labels", "label-key"}),
		},
		{
			Description: "deleting leaf key that does not exist",
			Res: `
metadata:
  labels: {}`,
			Expected: `
metadata:
  labels: {}`,
			Path: ctlres.NewPathFromStrings([]string{"metadata", "labels", "label-key"}),
		},
		{
			Description: "deleting parent key that exists",
			Res: `
metadata:
  labels:
    label-key: label-val`,
			Expected: `
metadata: {}`,
			Path: ctlres.NewPathFromStrings([]string{"metadata", "labels"}),
		},
		{
			Description: "deleting parent key that does not exist",
			Res: `
metadata:
  labels:
    label-key: label-val`,
			Expected: `
metadata:
  labels:
    label-key: label-val`,
			Path: ctlres.NewPathFromStrings([]string{"metadata", "not-labels"}),
		},
		{
			Description: "deleting leaf key that exists under array",
			Res: `
metadata:
  labels:
  - label-key: label-val`,
			Expected: `
metadata:
  labels:
  - {}`,
			Path: ctlres.Path{
				ctlres.NewPathPartFromString("metadata"),
				ctlres.NewPathPartFromString("labels"),
				ctlres.NewPathPartFromIndex(0),
				ctlres.NewPathPartFromString("label-key"),
			},
		},
		{
			Description: "deleting leaf key that does not exist under array",
			Res: `
metadata:
  labels:
  - label-key: label-val`,
			Expected: `
metadata:
  labels:
  - label-key: label-val`,
			Path: ctlres.Path{
				ctlres.NewPathPartFromString("metadata"),
				ctlres.NewPathPartFromString("labels"),
				ctlres.NewPathPartFromIndex(0),
				ctlres.NewPathPartFromString("not-label-key"),
			},
		},
		{
			Description: "deleting multiple keys that do or do not exist under array",
			Res: `
metadata:
  labels:
  - label-key: label-val
  - not-label-key: label-val`,
			Expected: `
metadata:
  labels:
  - {}
  - not-label-key: label-val`,
			Path: ctlres.Path{
				ctlres.NewPathPartFromString("metadata"),
				ctlres.NewPathPartFromString("labels"),
				ctlres.NewPathPartFromIndexAll(),
				ctlres.NewPathPartFromString("label-key"),
			},
		},
	}

	for _, ex := range exs {
		ex.Check(t)
	}
}

type modFieldRemoveExample struct {
	Description string
	Res         string
	Path        ctlres.Path
	Expected    string
}

func (e modFieldRemoveExample) Check(t *testing.T) {
	res, err := ctlres.NewResourceFromBytes([]byte(e.Res))
	if err != nil {
		t.Fatalf("Expected no err, but was %s", err)
	}

	err = ctlres.FieldRemoveMod{
		ResourceMatcher: ctlres.AllResourceMatcher{},
		Path:            e.Path,
	}.Apply(res)
	if err != nil {
		t.Fatalf("Expected no err, but was %s", err)
	}

	resultBs, err := res.AsYAMLBytes()
	if err != nil {
		t.Fatalf("Expected no err, but was %s", err)
	}

	expectEqualsStripped(t, e.Description, string(resultBs), e.Expected)
}