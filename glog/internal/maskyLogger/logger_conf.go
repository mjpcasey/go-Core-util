package maskyLogger

import (
	"fmt"
	"strings"
	"sync"
)

//logger的配置类
type LoggerConf struct {
	Name      string
	Levels    [LENGTH]bool
	Appenders []Appender
}

func (self *LoggerConf) SetAppender(appenders ...Appender) {
	self.Appenders = appenders
}

func (self *LoggerConf) SetLevel(level int) {
	for _, l := range LogLevelMap {
		if l >= level {
			self.Levels[l] = true
		} else {
			self.Levels[l] = false
		}
	}
}

func (self *LoggerConf) SetOnlyLevels(levels ...int) {
	for _, l := range levels {
		self.Levels[l] = true
	}
}

//在self 的配置为空时 ，复制from配置
func (self *LoggerConf) copynx(from *LoggerConf) {
	if from == nil {
		return
	}

	if len(self.Appenders) == 0 && len(from.Appenders) != 0 {
		self.Appenders = from.Appenders
	}

	for l, v := range from.Levels {
		self.Levels[l] = v
	}
}

//判断配置是否完整
func (self *LoggerConf) complete() bool {
	return len(self.Appenders) != 0 && len(self.Levels) != 0
}

//日志配置树，用于可以获取可继承的日志配置
//根节点为""
//日志配置在树中的节点与名称有关，如a/b/c则是 ""->a->b->c 的树
type Tree struct {
	Root *node
}

//日志配置数节点
type node struct {
	name     string
	parent   *node
	children map[string]*node
	mutex    sync.RWMutex

	current *LoggerConf
	final   *LoggerConf
}

func NewTree(root *LoggerConf) *Tree {
	if !root.complete() {
		panic("日志Root配置不完整")
	}
	return &Tree{
		Root: newNode("", nil, root),
	}
}
func printNode(tree *node, tabs int) {
	tstr := ""
	for i := 0; i < tabs; i++ {
		tstr += "\t"
	}
	fmt.Println(tstr + tree.name)
	fmt.Println("----------")
	tree.mutex.RLock()
	for _, c := range tree.children {
		printNode(c, tabs+1)
	}
	tree.mutex.RUnlock()
	fmt.Println("----------")
}

// 拷贝一颗树，其中 current,final 配置使用原来的指针对象，从而可以更新Logger使用的配置
func (t *Tree) clone() *Tree {
	newRoot := &node{}
	t.Root.clone(newRoot)
	return &Tree{
		Root: newRoot,
	}
}

//在树中插入一个配置
func (t *Tree) insert(logger *LoggerConf) {
	t.Root.addChild(logger.Name, logger)
}
func (t *Tree) updateConf(logger *LoggerConf) {
	t.Root.updateConf(logger.Name, logger)
}

//通过名称获取一个配置
func (t *Tree) get(name string) *LoggerConf {
	if name == "" {
		return t.Root.current
	}
	child := t.Root.child(name)
	if child != nil {
		return child.current
	}
	return nil
}

//获取name的配置，当name的配置为空时，会继承name上级最接近的非空配置
func (t *Tree) inheritConf(name string) *LoggerConf {
	return t.Root.generate(name).inheritConf()
}

func newNode(name string, parent *node, current *LoggerConf) *node {
	return &node{
		name:     name,
		parent:   parent,
		current:  current,
		children: make(map[string]*node),
	}
}

func (n *node) clone(nNode *node) {
	nNode.name = n.name
	nNode.current = n.current
	nNode.final = n.final
	nNode.children = map[string]*node{}
	n.mutex.RLock()
	for key, child := range n.children {
		sNode := &node{parent: nNode}
		child.clone(sNode)
		nNode.children[key] = sNode
	}
	n.mutex.RUnlock()
}
func (n *node) updateConf(name string, logger *LoggerConf) {
	if name == "" {
		if n.isRoot() {
			n.current = logger
			n.resetFinalConf()
		}
		return
	}
	arr := strings.Split(name, "/")
	var son *node
	n.mutex.Lock()
	if n, ok := n.children[arr[0]]; ok {
		son = n
	} else {
		son = newNode(arr[0], n, nil)
		n.children[arr[0]] = son
	}
	n.mutex.Unlock()

	if len(arr) == 1 {
		son.current = logger
		son.resetFinalConf()
	} else if len(arr) > 1 {
		son.updateConf(strings.Join(arr[1:], "/"), logger)
	}
}

//添加节点的子节点
func (n *node) addChild(name string, logger *LoggerConf) {
	if name == "" {
		return
	}
	arr := strings.Split(name, "/")
	var son *node
	n.mutex.Lock()
	if n, ok := n.children[arr[0]]; ok {
		son = n
	} else {
		son = newNode(arr[0], n, nil)
		n.children[arr[0]] = son
	}
	n.mutex.Unlock()

	if len(arr) == 1 {
		son.current = logger
	} else if len(arr) > 1 {
		son.addChild(strings.Join(arr[1:], "/"), logger)
	}
}

//通过name获取节点的子节点
func (n *node) child(name string) (ret *node) {
	if name == "" {
		return
	}
	arr := strings.Split(name, "/")
	n.mutex.RLock()
	if son, ok := n.children[arr[0]]; ok {
		if len(arr) == 1 {
			ret = son
		} else {
			ret = son.child(strings.Join(arr[1:], "/"))
		}
	}
	n.mutex.RUnlock()
	return
}

// 生成子节点如果不存在并返回
func (n *node) generate(name string) (ret *node) {
	if name == "" {
		if n.name == "" {
			return n
		}
		return
	}
	arr := strings.Split(name, "/")
	n.mutex.Lock()
	son, ok := n.children[arr[0]]
	if !ok {
		son = newNode(arr[0], n, nil)
		n.children[arr[0]] = son
	}
	n.mutex.Unlock()
	if len(arr) == 1 {
		ret = son
	} else {
		ret = son.generate(strings.Join(arr[1:], "/"))
	}
	return ret
}

//判断节点是不是根节点
func (n *node) isRoot() bool {
	return n.name == ""
}

//获取节点的上级中最接近的配置非空的节点
func (n *node) higher() *node {
	if n.parent != nil {
		if n.parent.current != nil {
			return n.parent
		}
		return n.parent.higher()
	}
	return nil
}

//获取当前节点的配置，当配置为空时，会继承name上级最接近的非空配置
func (n *node) inheritConf() *LoggerConf {
	if n.final == nil {
		var cfg = &LoggerConf{}
		var cur = n
		cfg.copynx(cur.current)
		if !cfg.complete() && !cur.isRoot() {
			higher := n.parent.inheritConf()
			cfg.copynx(higher)
		}
		n.final = cfg
	}
	return n.final
}

func (n *node) resetFinalConf() {
	var cfg = &LoggerConf{}
	var cur = n
	cfg.copynx(cur.current)
	if !cfg.complete() && !cur.isRoot() {
		higher := n.parent.inheritConf()
		cfg.copynx(higher)
	}
	if n.final == nil {
		n.final = cfg
	} else {
		n.final.Name = cfg.Name
		n.final.Levels = cfg.Levels
		n.final.Appenders = cfg.Appenders
	}
	n.mutex.Lock()
	for _, c := range n.children {
		c.resetFinalConf()
	}
	n.mutex.Unlock()
}
