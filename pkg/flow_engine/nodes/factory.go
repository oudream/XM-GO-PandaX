package nodes

import (
	"fmt"
)

const (
	NODE_CATEGORY_TRANSFORM = "transform"
	NODE_CATEGORY_EXTERNAL  = "external"
	NODE_CATEGORY_OTHERS    = "others"
)

// Factory is node's factory to create node based on metadata
// factory also manage node's metadta description which can be used by other
// service to present node in web
type Factory interface {
	Name() string
	Category() string
	Labels() []string
	Create(id string, meta Metadata) (Node, error)
}

var (
	// allNodeFactories hold all node's factory
	allNodeFactories map[string]Factory = make(map[string]Factory)

	// allNodeCategories hold node's metadata by category
	allNodeCategories map[string][]map[string]interface{} = make(map[string][]map[string]interface{})
)

// RegisterFactory add a new node factory and classify its category for
// metadata description
func RegisterFactory(f Factory) {
	allNodeFactories[f.Name()] = f

	if allNodeCategories[f.Category()] == nil {
		allNodeCategories[f.Category()] = []map[string]interface{}{}
	}
	allNodeCategories[f.Category()] = append(allNodeCategories[f.Category()], map[string]interface{}{"name": f.Name(), "labels": f.Labels()})
}

// NewNode is the only way to create a new node
func NewNode(nodeType string, id string, meta Metadata) (Node, error) {
	if f, found := allNodeFactories[nodeType]; found {
		return f.Create(id, meta)
	}
	return nil, fmt.Errorf("invalid node type '%s'", nodeType)
}

// GetCategoryNodes return specified category's all nodes
func GetCategoryNodes() map[string][]map[string]interface{} { return allNodeCategories }
