package storage
import (
    . "mustard/storage/mongo"
)

/*
    Storage API Layer
*/

var STORAGE_ENGINE_IMPL = NewMongoStorageEngine()
