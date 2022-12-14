package datastore

import (
	"errors"

	"github.com/darkcat013/pr-datastore/utils"
)

var ds = make(map[string][]byte)

func Get(id string) ([]byte, error) {
	utils.Log.Infof("Datastore GET | Get data with id: %s", id)

	if value, ok := ds[id]; ok {
		return value, nil
	}

	utils.Log.Infof("Datastore GET | Id not found: %s", id)
	return nil, errors.New("id not found: " + id)
}

func Insert(value []byte) (string, error) {
	utils.Log.Infof("Datastore INSERT | Insert new data")

	id := utils.GetNewId(value)

	if _, ok := ds[id]; !ok {
		ds[id] = value
		utils.Log.Infof("Datastore INSERT | Inserted new data at id: %s", id)
		return id, nil
	}

	utils.Log.Infof("Datastore INSERT | Id already exists %s", id)
	return "", errors.New("id already exists: " + id)
}

func InsertAtId(id string, value []byte) error {
	utils.Log.Infof("Datastore INSERT AT ID | Insert new data at id %s", id)

	if _, ok := ds[id]; !ok {
		ds[id] = value
		return nil
	}

	utils.Log.Infof("Datastore INSERT | Id already exists %s", id)
	return errors.New("id already exists: " + id)
}

func Update(id string, value []byte) error {
	utils.Log.Infof("Datastore UPDATE | Update data with id: %s", id)

	if _, ok := ds[id]; ok {
		ds[id] = value
		return nil
	}

	utils.Log.Infof("Datastore UPDATE | Id not found: %s", id)
	return errors.New("id not found: " + id)
}

func Delete(id string) error {
	utils.Log.Infof("Datastore DELETE | Delete data with id: %s", id)

	if _, ok := ds[id]; ok {
		delete(ds, id)
		return nil
	}

	utils.Log.Infof("Datastore DELETE | Id not found: %s", id)
	return errors.New("id not found: " + id)
}

func GetAllKeys() []string {
	utils.Log.Infof("Datastore GetAllKeys | Get all keys from datastore")

	keys := make([]string, 0, len(ds))

	for k := range ds {
		keys = append(keys, k)
	}

	return keys
}
