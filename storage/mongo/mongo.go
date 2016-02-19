package mongo

import (
    "sync"
)

type MongoStorageEngine struct {
    // crawldoc interface

}



///////////// Singleton  ////////////////////////
var _instance *MongoStorageEngine = nil
var _init_ctx sync.Once

func NewMongoStorageEngine() *MongoStorageEngine {
    _init_ctx.Do(func(){
        _instance = &MongoStorageEngine{
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