package dbmanager

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MongoDb连接对象
const DefaultIdCounter = "sys.IdCounter"

type collectSession struct {
	collection *mgo.Collection
	session    *mgo.Session
}

// MongoWrapper Mongo查询代理接口定义
type MongoWrapper interface {
	GetQuery(colName string, cond interface{}) (*QueryWrapper, error)
	Find(colName string, selector *MongoSelector, result interface{}) error
	Insert(colName string, docs ...interface{}) error
	Remove(colName string, selector interface{}) error
	FindAndCount(colName string, selector interface{}) (n int, err error)
	FindOne(colName string, selector interface{}, result interface{}) (err error)
	Exec(colName string, exec func(*mgo.Collection) error) error
}

// mongodb 查询结构
type MongoSelector struct {
	Cond   interface{}
	Select interface{}
	Limit  int
	Sort   []string
	Skip   int
}

// QueryWrapper 封装mgo游标
type QueryWrapper struct {
	*mgo.Query
	colSession *collectSession
}

// query的close方法
func (q *QueryWrapper) Close() {
	q.colSession.session.Close()
}

type mongoWrapper struct {
	config struct {
		Name       string   `json:"name"`
		Username   string   `json:"username"`
		Password   string   `json:"password"`
		Addrs      []string `json:"addrs"`
		Database   string   `json:"database"`
		ReplicaSet string   `json:"replicaSet"`
		Mechanism  string   `json:"mechanism"`
		IdCounter  string   `json:"idCounter"`

		TimeoutSec      int `json:"timeoutSec"`      // 初始化连接时间限制
		PoolTimeoutSec  int `json:"poolTimeoutSec"`  // 等待连接池返回可用连接的时间限制
		ReadTimeoutSec  int `json:"readTimeoutSec"`  // 读操作时间限制
		WriteTimeoutSec int `json:"writeTimeoutSec"` // 写操作时间限制
	}

	session *mgo.Session
}

// GetQuery 提供获取query的方法，但是使用后需要close
func (m *mongoWrapper) GetQuery(col string, cond interface{}) (queryWrapper *QueryWrapper, err error) {
	colSession := m.getCollectionSession(col)

	queryWrapper = &QueryWrapper{
		colSession: colSession,
		Query:      colSession.collection.Find(cond),
	}

	return
}

// GetNextId 获取自增id
func (m *mongoWrapper) GetNextId(name string) (id int) {
	cs := m.getCollectionSession(DefaultIdCounter)

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	doc := struct{ Id int }{}

	if _, err := cs.collection.Find(bson.M{"_id": name}).Apply(change, &doc); err != nil {
		logger.Errorf("获取自增ID[collection=%s]失败: %s", name, err)
	}

	cs.session.Close()
	return doc.Id
}

// 创建 collection 方法
func (m *mongoWrapper) Create(colName string, info *mgo.CollectionInfo) (err error) {
	s := m.getCollectionSession(colName)

	err = s.collection.Create(info)

	s.session.Close()
	return
}

func (m *mongoWrapper) Find(colName string, selector *MongoSelector, result interface{}) (err error) {
	cs := m.getCollectionSession(colName)

	err = cs.collection.Find(selector.Cond).Select(selector.Select).Sort(selector.Sort...).Skip(selector.Skip).Limit(selector.Limit).All(result)

	cs.session.Close()

	return
}

func (m *mongoWrapper) Insert(colName string, docs ...interface{}) (err error) {
	cs := m.getCollectionSession(colName)

	err = cs.collection.Insert(docs...)

	cs.session.Close()

	return
}

func (m *mongoWrapper) Update(colName string, selector interface{}, update interface{}) (err error) {
	cs := m.getCollectionSession(colName)

	err = cs.collection.Update(selector, update)

	cs.session.Close()

	return
}

func (m *mongoWrapper) Remove(colName string, selector interface{}) (err error) {
	cs := m.getCollectionSession(colName)

	err = cs.collection.Remove(selector)

	cs.session.Close()

	return
}

func (m *mongoWrapper) FindAndCount(colName string, selector interface{}) (n int, err error) {
	cs := m.getCollectionSession(colName)

	n, err = cs.collection.Find(selector).Count()

	cs.session.Close()

	return
}

func (m *mongoWrapper) FindOne(colName string, selector interface{}, result interface{}) (err error) {
	cs := m.getCollectionSession(colName)

	err = cs.collection.Find(selector).One(result)

	cs.session.Close()

	return
}

func (m *mongoWrapper) getCollectionSession(col string) (cs *collectSession) {
	session := m.session.Copy()

	cs = &collectSession{
		session:    session,
		collection: session.DB(m.config.Database).C(col),
	}

	return
}

func (m *mongoWrapper) Exec(colName string, exec func(*mgo.Collection) error) (err error) {
	cs := m.getCollectionSession(colName)

	err = exec(cs.collection)

	cs.session.Close()
	return
}
