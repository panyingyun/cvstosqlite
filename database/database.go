package database

import (
	"github.com/panyingyun/cvstosqlite/model"

	log "github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

//DBEngin
type DBEngine struct {
	dbengine *xorm.Engine
}

//Create DBEngin
func NewDBEngine(dburl string) (*DBEngine, error) {
	db, err := xorm.NewEngine("sqlite3", dburl)
	if err != nil {
		return nil, err
	}

	db.ShowSQL(false)
	log.Infof("connect to database(%v) server OK!", dburl)

	db.Sync2(new(model.Node))

	engine := DBEngine{
		dbengine: db,
	}

	return &engine, nil
}

//Insert All Node Data
func (db *DBEngine) InsertAllNodeData(nodes []model.Node) error {
	_, err := db.dbengine.Insert(&nodes)
	if err != nil {
		log.Errorf("InsertAllNodeData  fail %v!", err)
	}
	return err
}
