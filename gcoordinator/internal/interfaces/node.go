package interfaces

import (
	"context"
	"gcore/gcoordinator/gcoordinatorTypes"
)

// 协调器节点
type CoordinatorNode interface {
	Create(path string, data []byte) (CoordinatorNode, error)
	Open(path string) (CoordinatorNode, error)
	Set(data []byte) error
	Get() ([]byte, error)
	Remove() error
	Refresh() error
	Watch(context context.Context, events int, callback gcoordinatorTypes.EventCallback) error
	GetChildren() ([]string, error)
	GetChildrenNode() (cNodes []CoordinatorNode, err error)
}
