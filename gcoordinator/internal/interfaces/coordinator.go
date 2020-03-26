package interfaces

// 协调器
type Coordinator interface {
	Start() (err error)
	Stop() (err error)
	CreateNode(path string, data []byte, store bool) (CoordinatorNode, error)
	Open(path string) (CoordinatorNode, error)
	Exist(path string) (bool, error)
}