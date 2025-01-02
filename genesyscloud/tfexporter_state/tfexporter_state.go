package tfexporter_state

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"
	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	tfVersion "github.com/hashicorp/terraform/version"

	"github.com/hashicorp/go-uuid"
)

/*
Export as v4 Terraform State
Much of this code is based on the internal Terraform State V4 format found here https://github.com/hashicorp/terraform/blob/main/internal/states/statefile/version4.go
*/

type TerraformStateV4 struct {
	Version          int                               `json:"version"`
	TerraformVersion string                            `json:"terraform_version"`
	Serial           uint64                            `json:"serial"`
	Lineage          string                            `json:"lineage"`
	Resources        []TerraformStateResourceV4        `json:"resources"`
	OutputValues     map[string]TerraformStateOutputV4 `json:"outputs,omitempty"`
}

type TerraformStateResourceV4 struct {
	Module    string                     `json:"module,omitempty"`
	Mode      string                     `json:"mode"`
	Type      string                     `json:"type"`
	Name      string                     `json:"name"`
	Provider  string                     `json:"provider"`
	Instances []TerraformStateInstanceV4 `json:"instances"`
}

type TerraformStateInstanceV4 struct {
	IndexKey                interface{}            `json:"index_key,omitempty"`
	Status                  string                 `json:"status,omitempty"`
	Deposed                 string                 `json:"deposed,omitempty"`
	SchemaVersion           uint64                 `json:"schema_version"`
	AttributeSensitivePaths []string               `json:"sensitive_attributes"`
	Attributes              map[string]interface{} `json:"attributes"`
	Private                 string                 `json:"private,omitempty"`
	Dependencies            []string               `json:"dependencies"`
	CreateBeforeDestroy     bool                   `json:"create_before_destroy,omitempty"`
}

type TerraformStateOutputV4 struct {
	Value     interface{} `json:"value"`
	Type      interface{} `json:"type,omitempty"`
	Sensitive bool        `json:"sensitive,omitempty"`
}

type sortResourcesV4 []TerraformStateResourceV4

func (sr sortResourcesV4) Len() int      { return len(sr) }
func (sr sortResourcesV4) Swap(i, j int) { sr[i], sr[j] = sr[j], sr[i] }
func (sr sortResourcesV4) Less(i, j int) bool {
	switch {
	case sr[i].Module != sr[j].Module:
		return sr[i].Module < sr[j].Module
	case sr[i].Mode != sr[j].Mode:
		return sr[i].Mode < sr[j].Mode
	case sr[i].Type != sr[j].Type:
		return sr[i].Type < sr[j].Type
	case sr[i].Name != sr[j].Name:
		return sr[i].Name < sr[j].Name
	default:
		return false
	}
}

type sortInstancesV4 []TerraformStateInstanceV4

func (si sortInstancesV4) Len() int      { return len(si) }
func (si sortInstancesV4) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si sortInstancesV4) Less(i, j int) bool {
	ki := si[i].IndexKey
	kj := si[j].IndexKey
	if ki != kj {
		if (ki == nil) != (kj == nil) {
			return ki == nil
		}
		if kii, isInt := ki.(int); isInt {
			if kji, isInt := kj.(int); isInt {
				return kii < kji
			}
			return true
		}
		if kis, isStr := ki.(string); isStr {
			if kjs, isStr := kj.(string); isStr {
				return kis < kjs
			}
			return true
		}
	}
	if si[i].Deposed != si[j].Deposed {
		return si[i].Deposed < si[j].Deposed
	}
	return false
}

// normalize makes some in-place changes to normalize the way items are
// stored to ensure that two functionally-equivalent states will be stored
// identically.
func (s *TerraformStateV4) normalize() {
	sort.Stable(sortResourcesV4(s.Resources))
	for _, rs := range s.Resources {
		sort.Stable(sortInstancesV4(rs.Instances))
	}
}

// sets default values for Attributes based on Cty Types
func getDefaultValueFromCtyType(t cty.Type) cty.Value {
	switch {
	case t.Equals(cty.String):

		return cty.StringVal("")
	case t.Equals(cty.Number):
		return cty.NumberIntVal(0)
	case t.Equals(cty.Bool):
		return cty.False
	case t.IsListType() || t.IsSetType():
		return cty.ListValEmpty(t.ElementType())
	case t.IsMapType():
		return cty.MapValEmpty(t.ElementType())
	case t.IsObjectType():
		attrs := make(map[string]cty.Value)
		for k, at := range t.AttributeTypes() {
			attrs[k] = getDefaultValueFromCtyType(at)
		}
		return cty.ObjectVal(attrs)
	default:
		return cty.NullVal(t)
	}
}

/*
Export state is used to indicate whether an export is being done.  If the export state is set to true, then this should be
a signal that any resources being exported should be reading their data from each resource's internal cache rather then the API.
*/
var (
	stateMutex  sync.Mutex
	exportState bool
	once        sync.Once
)

// ActivateExporterState will be used to indicate that caching should be used to process requests.
// We are setting this as an environment variable so we can experiment with it, without creating an attribute
// on the resource
func ActivateExporterState() {
	once.Do(func() {
		log.Printf("Exporter State is active")
		exportState = true
	})
}

func IsExporterActive() bool {
	return exportState
}

func GenerateTerraformStateV4(resources []resourceExporter.ResourceInfo, providerSource string) (*TerraformStateV4, error) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	lineage, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate lineage: %v", err)
	}

	state := &TerraformStateV4{
		Version:          4,
		TerraformVersion: tfVersion.Version,
		Serial:           0,
		Lineage:          lineage,
		Resources:        make([]TerraformStateResourceV4, 0),
	}

	for _, res := range resources {
		stateResource := TerraformStateResourceV4{
			Mode:     res.Mode,
			Type:     res.Type,
			Name:     res.Name,
			Provider: fmt.Sprintf("provider[\"%s\"]", providerSource),
			Instances: []TerraformStateInstanceV4{
				{
					Attributes:              make(map[string]interface{}),
					SchemaVersion:           1,
					Dependencies:            make([]string, 0),
					AttributeSensitivePaths: make([]string, 0),
				},
			},
		}

		for key, value := range res.StateAttributes {
			// If no value, set to appropriate default value depending on type
			if value == nil {
				ctyType := res.CtyType.AttributeTypes()[key]
				ctyValue := getDefaultValueFromCtyType(ctyType)
				defaultBytes, err := ctyjson.Marshal(ctyValue, ctyType)
				if err != nil {
					return nil, fmt.Errorf("Failed to marshal default value: %v", err)
				}
				var extractedValue interface{}
				err = json.Unmarshal(defaultBytes, &extractedValue)
				if err != nil {
					return nil, fmt.Errorf("Failed to unmarshal default value: %v", err)
				}
				stateResource.Instances[0].Attributes[key] = extractedValue
			} else {
				stateResource.Instances[0].Attributes[key] = value
			}
		}

		state.Resources = append(state.Resources, stateResource)
	}

	// Sort the state so that it is deterministic
	state.normalize()

	return state, nil
}
