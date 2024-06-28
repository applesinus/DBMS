package database

import "DBMS/task0+3"

var database = task0.CreateDB()
var settings = make(map[string]string)
var avaliableCollectionTypes = []string{"BI", "AVL", "RB", "Btree"}

func Database() *task0.Database { return database }

func SetSettings(sets map[string]string) { settings = sets }

func AvaliableCollectionTypes() []string { return avaliableCollectionTypes }
