package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"projekat/model"
)

const (
	//namespace/namespaceid
	namespaceKey = "namespace/%s"

	//application/namespaceId/applicationId
	applicationKey = "application/%s/%s"

	//dataspace/applicationid/dataspaceid
	dataSpaceKey = "dataspace/%s/%s"

	//dataspaceitem/path/dataspaceitemid
	dataSpaceItemKey = "dataspaceitem/%s/%s"

	//hardlink/namespaceid/applicationid/dataspaceitemid
	hardlinkKey = "hardlink/%s/%s"
)

func key_ns(id, template string) string {
	return fmt.Sprintf(template, id)
}

func key(id1, id2, template string) string {
	return fmt.Sprintf(template, id1, id2)
}

type DB struct {
	Kv     clientv3.KV
	Client *clientv3.Client
}

func New(endpoints []string, timeout time.Duration) (model.Store, error) {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: timeout,
		Endpoints:   endpoints,
	})

	if err != nil {
		return nil, err
	}

	return &DB{
		Kv:     clientv3.NewKV(cli),
		Client: cli,
	}, nil
}

func (db *DB) Close() { db.Client.Close() }

func (db *DB) PutApp(app *model.Application) error {
	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()

	key := key(app.ParentNamespaceId, app.ApplicationId, applicationKey)

	jsonData, err := json.Marshal(app)
	if err != nil {
		return err
	}

	_, err = db.Kv.Put(ctx, key, string(jsonData))

	return err

}

func (db *DB) PutDataSpace(applicationID string, dataSpace *model.DataSpace) error {

	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()

	key := key(applicationID, dataSpace.DataSpaceId, dataSpaceKey)

	jsonData, err := json.Marshal(dataSpace)
	if err != nil {
		return err
	}

	_, err = db.Kv.Put(ctx, key, string(jsonData))

	return err
}

func (db *DB) PutDataSpaceItem(dataSpaceItem *model.DataSpaceItem) error {

	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	path := fmt.Sprintf("%s/%s", dataSpaceItem.Path, dataSpaceItem.Name)
	key := key(path, dataSpaceItem.DataSpaceItemId, dataSpaceItemKey)

	jsonData, err := json.Marshal(dataSpaceItem)
	if err != nil {
		return err
	}

	_, err = db.Kv.Put(ctx, key, string(jsonData))

	return err
}

func (db *DB) PutHardlink(hardlink *model.Hardlink) error {

	ctx, cncl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cncl()
	key := key(hardlink.ApplicationID, hardlink.DataSpaceItemID, hardlinkKey)

	jsonData, err := json.Marshal(hardlink)
	if err != nil {
		return err
	}

	_, err = db.Kv.Put(ctx, key, string(jsonData))

	return err
}

func (db *DB) GetApp(namespaceID, appID string) (*model.Application, error) {
	key := key(namespaceID, appID, applicationKey)
	resp, err := db.Kv.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		log.Fatalf("No data found for the key")
	}

	var app model.Application
	err = json.Unmarshal(resp.Kvs[0].Value, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
