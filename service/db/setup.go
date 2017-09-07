package db

import (
	"gopkg.in/mgo.v2"
)

type Db struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewDb(session *mgo.Session, dbName string, collectionName string) *Db {
	db := Db{}
	db.session = session
	db.dbName = dbName
	db.collectionName = collectionName
	return &db
}

func (d *Db) EnsureSessionActive() error {
	if err := d.session.Ping(); err != nil {
		return err
	}
	return nil
}

func (d *Db) DropCollection() error {
	err := d.session.DB(d.dbName).C(d.collectionName).DropCollection()
	if err != nil {
		return err
	}
	err = d.session.DB(d.dbName).C(d.collectionName).Create(&mgo.CollectionInfo{ForceIdIndex: true})
	if err != nil {
		return err
	}
	return nil
}

func (d *Db) EnsureIndex(fields []string) error {
	index := mgo.Index{
		Key:        fields,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := d.session.DB(d.dbName).C(d.collectionName).EnsureIndex(index)
	if err != nil {
		return err
	}
	return nil
}
