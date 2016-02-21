package mongo

import (
    "sync"
    . "mustard/crawl/proto"
    "gopkg.in/mgo.v2"
    LOG "mustard/base/log"
    "mustard/base/time_util"
)

const (
    kDefaultDB = "mustard"
    kDefaultCollection = "record"
)

type MongoStorageEngine struct {
    // crawldoc interface
    connected bool
    collection *mgo.Collection
    session *mgo.Session
    servers,db,table string
}
func (m *MongoStorageEngine)connect(db, table string) {
    if !m.connected {
        for true {
            var err error
            m.session,err = mgo.Dial(m.servers)
            if err != nil {
                LOG.Errorf("Connect MongoDb Error %s : %s", m.servers,err.Error())
                time_util.Sleep(2)
            } else {
                break
            }
        }
        m.session.SetMode(mgo.Monotonic,true)
        LOG.VLog(2).Debugf("Connect MongoDB to %s", m.servers)
    }
    m.connected = true
    if m.table != table || m.db != db || m.collection == nil{
        m.table,m.db = table,db
        m.collection = m.session.DB(m.db).C(m.table)
        LOG.VLog(2).Debugf("MongoDB Use %s : %s", m.db,m.table)
    }
}
func (m *MongoStorageEngine)Init(servers string) {
    /*
    input params:
        servers: 192.168.0.3:27017,192.168.0.3:27018
        db: database, no need create
        table:  collection
    */
    m.servers,m.db,m.table = servers,kDefaultDB,kDefaultCollection
    m.connect(m.db,m.table)
}
func (m *MongoStorageEngine)WithDb(db string) *MongoStorageEngine{
    m.connect(db, m.table)
    return m
}
func (m *MongoStorageEngine)WithTable(table string) *MongoStorageEngine{
    m.connect(m.db,table)
    return m
}
func (m *MongoStorageEngine)Save(doc *CrawlDoc) error {
    m.connect(m.db,m.table)
    err := m.collection.Insert(doc)
    if err != nil {
        m.connected = false
        LOG.Errorf("Insert Error[%s] (%s/%s:%s), url:%s", err.Error(),
                    m.servers,m.db,m.table,doc.Url)
        return err
    }
    return nil
}
func (m *MongoStorageEngine)SaveBatchSavage(docs []*CrawlDoc) (int,error) {
    for _,doc := range docs {
        m.Save(doc)
    }
    return len(docs),nil
}
func (m *MongoStorageEngine)SaveBatch(docs []*CrawlDoc) (int,error) {
    // convert slice to ...
    container := make([]interface{},len(docs))
    for i,v := range docs {
        container[i] = interface{}(v)
    }
    err := m.collection.Insert(container...)
    if err != nil {
        m.connected = false
        LOG.Errorf("Insert Error[%s] (%s/%s:%s)", err.Error(),
            m.servers,m.db,m.table)
        return 0,err
    }
    return len(docs),nil
}
// TODO traverse batch interface, use callback maybe better


///////////// Singleton  ////////////////////////
var _instance *MongoStorageEngine = nil
var _init_ctx sync.Once

func NewMongoStorageEngine() *MongoStorageEngine {
    _init_ctx.Do(func(){
        _instance = &MongoStorageEngine{
            connected:false,
        }
    })
    return _instance
}

func NewMongoStorageEngine2() *MongoStorageEngine {
    if _instance == nil {
        // Not thread safe....
        _instance = &MongoStorageEngine{}
    }
    return _instance
}