package zkCoordinator

import (
	"context"
	"fmt"
	"gcore/gcoordinator/gcoordinatorTypes"
	"gcore/gcoordinator/internal/interfaces"
	"github.com/samuel/go-zookeeper/zk"
)

type node struct {
	running   bool
	path      string
	realpath  string
	stat      *zk.Stat
	zookeeper *coordinator
}

func (n *node) init(path string, zk *coordinator) error {
	exist, stat, err := zk.conn.Exists(zk.config.Root + path)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("path[%s] not exist", path)
	}

	n.path = path
	n.realpath = zk.config.Root + path
	n.stat = stat
	n.zookeeper = zk

	return nil
}

// CreateNode create children node
func (n *node) Create(path string, data []byte) (interfaces.CoordinatorNode, error) {
	return n.zookeeper.CreateNode(n.path+`/`+path, data, false)
}

// Open open a children node
func (n *node) Open(path string) (interfaces.CoordinatorNode, error) {
	return n.zookeeper.Open(n.path + `/` + path)
}

// Set set node data
func (n *node) Set(data []byte) (err error) {
	stat, err := n.zookeeper.conn.Set(n.realpath, data, n.stat.Version)
	if err == nil {
		n.stat = stat
	}
	return
}

// Get get node data
func (n *node) Get() ([]byte, error) {
	data, stat, err := n.zookeeper.conn.Get(n.realpath)
	if err == nil {
		n.stat = stat
	}
	return data, err
}

// GetChildren get node children names string
func (n *node) GetChildren() (cNodes []string, err error) {
	cNodes, _, err = n.zookeeper.conn.Children(n.realpath)
	return
}

// GetChildrenNode get Children node
func (n *node) GetChildrenNode() (cNodes []interfaces.CoordinatorNode, err error) {
	children, _, err := n.zookeeper.conn.Children(n.realpath)
	if err != nil {
		return
	}

	cNodes = make([]interfaces.CoordinatorNode, 0, len(children))

	for _, path := range children {
		child := new(node)
		e := child.init(n.path+`/`+path, n.zookeeper)
		if e != nil {
			logger.Errorf("zk children node init error: %s", e.Error())
			continue
		}
		cNodes = append(cNodes, child)
	}
	return
}

// Remove remove node
func (n *node) Remove() error {
	return n.zookeeper.conn.Delete(n.realpath, n.stat.Version)
}

// Refresh update node stat
func (n *node) Refresh() error {
	_, stat, err := n.zookeeper.conn.Exists(n.realpath)
	if err == nil {
		n.stat = stat
	}
	return err
}

// Watch watch node change
func (n *node) Watch(context context.Context, events int, callback gcoordinatorTypes.EventCallback) (err error) {
	conn := n.zookeeper.conn

	var dataEvents <-chan zk.Event = make(chan zk.Event, 1)
	var existEvents <-chan zk.Event = make(chan zk.Event, 1)
	var childrenEvents <-chan zk.Event = make(chan zk.Event, 1)

	var data []byte
	var children []string
	var exists bool

	if events&gcoordinatorTypes.EventChanged != 0 {
		if data, _, dataEvents, err = conn.GetW(n.realpath); err != nil {
			return err
		}
	}
	if events&gcoordinatorTypes.EventChildrenChanged != 0 {
		if children, _, childrenEvents, err = conn.ChildrenW(n.realpath); err != nil {
			return err
		}
	}
	if events&(gcoordinatorTypes.EventCreated|gcoordinatorTypes.EventDeleted) != 0 {
		if exists, _, existEvents, err = conn.ExistsW(n.realpath); err != nil {
			return err
		}
	}

	logger.Debugf("正在监听[%s]", n.realpath)
	for {
		var evt zk.Event
		select {
		case <-context.Done():
			return
		case evt = <-dataEvents:
			if evt.Type == zk.EventNodeDataChanged {
				event := gcoordinatorTypes.Event{
					Type: gcoordinatorTypes.EventChanged,
					Data: data,
				}
				if err := callback(event); err != nil {
					return err
				}
				if data, _, dataEvents, err = conn.GetW(n.realpath); err != nil {
					return err
				}
			}
		case evt = <-childrenEvents:
			if evt.Type == zk.EventNodeChildrenChanged {
				event := gcoordinatorTypes.Event{
					Type: gcoordinatorTypes.EventChildrenChanged,
					Data: children,
				}
				if err := callback(event); err != nil {
					return err
				}
				if children, _, childrenEvents, err = conn.ChildrenW(n.realpath); err != nil {
					return err
				}
			}
		case evt = <-existEvents:
			if evt.Type == zk.EventNodeCreated {
				if events&gcoordinatorTypes.EventCreated != 0 {
					event := gcoordinatorTypes.Event{
						Type: gcoordinatorTypes.EventCreated,
						Data: exists,
					}
					if err := callback(event); err != nil {
						return err
					}
				}
			} else if evt.Type == zk.EventNodeDeleted {
				if events&gcoordinatorTypes.EventDeleted != 0 {
					event := gcoordinatorTypes.Event{
						Type: gcoordinatorTypes.EventDeleted,
						Data: exists,
					}
					if err := callback(event); err != nil {
						return err
					}
				}
			} else {
				break
			}
			if exists, _, existEvents, err = conn.ExistsW(n.realpath); err != nil {
				return err
			}
		}

		if evt.Err != nil || evt.Type == zk.EventNotWatching {
			// 网络隔离断网重连后，导致监听时间方法推出，尝试重新链接监听
			if events&gcoordinatorTypes.EventChanged != 0 {
				if data, _, dataEvents, err = conn.GetW(n.realpath); err != nil {
					return err
				}
			}
			if events&gcoordinatorTypes.EventChildrenChanged != 0 {
				if children, _, childrenEvents, err = conn.ChildrenW(n.realpath); err != nil {
					return err
				}
			}
			if events&(gcoordinatorTypes.EventCreated|gcoordinatorTypes.EventDeleted) != 0 {
				if exists, _, existEvents, err = conn.ExistsW(n.realpath); err != nil {
					return err
				}
			}
		}
	}
}
