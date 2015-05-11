package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
	// "github.com/intelsdilabs/pulse/core"
)

var (
	InvalidPayload = errors.New("Payload to convert must be string or []byte")
)

// Represents a validated instance of a workflow for a task.
type scheduleWorkflow struct {
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func FromYaml(payload interface{}) (*WorkflowMap, error) {
	p, err := inStringBytes(payload)
	if err != nil {
		return nil, err
	}

	wmap := new(WorkflowMap)
	err = yaml.Unmarshal(p, wmap)
	if err != nil {
		return nil, err
	}
	return wmap, nil
}

func FromJson(payload interface{}) (*WorkflowMap, error) {
	p, err := inStringBytes(payload)
	if err != nil {
		return nil, err
	}

	wmap := new(WorkflowMap)
	err = json.Unmarshal(p, wmap)
	if err != nil {
		return nil, err
	}
	return wmap, nil
}

func inStringBytes(payload interface{}) ([]byte, error) {
	var p []byte
	switch tp := payload.(type) {
	case string:
		p = []byte(tp)
	case []byte:
		p = tp
	default:
		return p, InvalidPayload
	}
	return p, nil
}

func SampleWorkflowMapJSON() string {
	wf := new(WorkflowMap)

	c1 := &CollectWorkflowMapNode{}
	// pr1 := &ProcessWorkflowMapNode{Name: "learn", Version: 3}
	pu1 := &PublishWorkflowMapNode{Name: "rabbitmq", Version: 5}
	var e error
	// e = pr1.Add(pu1)
	// handleErr(e)
	e = c1.Add(pu1)
	handleErr(e)
	e = c1.AddMetricNamespace("/foo/bar")
	handleErr(e)

	wf.CollectNode = c1
	b, e := json.Marshal(wf)
	if e != nil {
		panic(e)
	}
	return string(b)
}

// A map of a desired workflow that is used to create a scheduleWorkflow
type WorkflowMap struct {
	CollectNode *CollectWorkflowMapNode `json:"collect"yaml:"collect"`
}

func (w *WorkflowMap) String() string {
	var out string
	out += "Workflow\n"
	out += "\tCollect:\n"
	if w.CollectNode != nil {
		out += w.CollectNode.String("\t\t")
	} else {
		out += "\n"
	}

	return out
}

type CollectWorkflowMapNode struct {
	MetricsNamespaces []string `json:"metric_namespaces"yaml:"metric_namespaces"`

	ProcessNodes []ProcessWorkflowMapNode `json:"process"yaml:"process"`
	PublishNodes []PublishWorkflowMapNode `json:"publish"yaml:"publish"`
}

func (c *CollectWorkflowMapNode) String(pad string) string {
	var out string
	out += pad + "Metric Namespaces:\n"
	for _, x := range c.MetricsNamespaces {
		out += pad + "\t\t" + x + "\n"
	}
	out += "\n"
	out += pad + "Process Nodes:\n"
	for _, pr := range c.ProcessNodes {
		out += pr.String(pad)
	}
	out += "\n"
	out += pad + "Publish Nodes:\n"
	for _, pu := range c.PublishNodes {
		out += pu.String(pad)
	}
	out += "\n"
	return out
}

func (c *CollectWorkflowMapNode) Add(node interface{}) error {
	switch x := node.(type) {
	case *ProcessWorkflowMapNode:
		c.ProcessNodes = append(c.ProcessNodes, *x)
	case *PublishWorkflowMapNode:
		c.PublishNodes = append(c.PublishNodes, *x)
	default:
		return errors.New(fmt.Sprintf("cannot add workflow node type (%v) to collect node as child", x))
	}
	return nil
}

func (c *CollectWorkflowMapNode) AddMetricNamespace(ns string) error {
	// TODO regex validation here that this matches /one/two/three format
	c.MetricsNamespaces = append(c.MetricsNamespaces, ns)
	return nil
}

type ProcessWorkflowMapNode struct {
	Name         string                   `json:"plugin_name"yaml:"plugin_name"`
	Version      int                      `json:"plugin_version"yaml:"plugin_version"`
	ProcessNodes []ProcessWorkflowMapNode `json:"process"yaml:"process"`
	PublishNodes []PublishWorkflowMapNode `json:"publish"yaml:"publish"`
	// TODO processor config
	Config interface{} `json:"processor_config"yaml:"processor_config"`
}

func (p *ProcessWorkflowMapNode) String(pad string) string {
	var out string
	out += pad + fmt.Sprintf("Name: %s\n", p.Name)
	out += pad + fmt.Sprintf("Version: %d\n", p.Version)

	out += pad + "Process Nodes:\n"
	for _, pr := range p.ProcessNodes {
		out += pr.String(pad)
	}
	out += pad + "Publsish Nodes:\n"
	for _, pu := range p.PublishNodes {
		out += pu.String(pad)
	}
	return out
}

func (p *ProcessWorkflowMapNode) Add(node interface{}) error {
	switch x := node.(type) {
	case *ProcessWorkflowMapNode:
		p.ProcessNodes = append(p.ProcessNodes, *x)
	case *PublishWorkflowMapNode:
		p.PublishNodes = append(p.PublishNodes, *x)
	default:
		return errors.New(fmt.Sprintf("cannot add workflow node type (%v) to process node as child", x))
	}
	return nil
}

type PublishWorkflowMapNode struct {
	Name    string `json:"plugin_name"yaml:"plugin_name"`
	Version int    `json:"plugin_version"yaml:"plugin_version"`
	// TODO publisher config
	Config interface{} `json:"publisher_config"yaml:"publisher_config"`
}

func (p *PublishWorkflowMapNode) String(pad string) string {
	var out string
	out += pad + fmt.Sprintf("\tName: %s\n", p.Name)
	out += pad + fmt.Sprintf("\tVersion: %d\n", p.Version)

	return out
}
