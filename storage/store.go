package storage

import (
	"time"
	"encoding/json"
	"github.com/coreos/etcd/client"
	log "github.com/Sirupsen/logrus"
	"context"
)

type Storage interface {
	Add(key string, obj interface{}) error
	Update(key string, obj interface{}) error
	Delete(key string) error
	List() []interface{}
	Get(key string) (item interface{}, err error)
}

// EtcdStorage implement interface Storage.
//var _ Storage = &EtcdStorage{}

type EtcdStorage struct {
	client client.Client
	endpoints string
	KeysAPI client.KeysAPI
}

func NewEtcdStorage(endpoints string) *EtcdStorage {
	cfg := client.Config {
		Endpoints: 		 []string{endpoints},
		Transport: 		 client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Error: cannot connect to etcd: %v", err)
	}
	return &EtcdStorage{
		client: etcdClient,
		endpoints: endpoints,
	}
}

func (e *EtcdStorage) Add(key string, obj interface{}) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = e.KeysAPI.Create(context.Background(), key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStorage) Update(key string, obj interface{}) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = e.KeysAPI.Set(context.Background(), key, string(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStorage) UpdateExpire(key string, obj interface{}, ttl time.Duration) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = e.KeysAPI.Set(context.Background(), key, string(value), &client.SetOptions{
		TTL: ttl,
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStorage) Delete(key string) error {
	_, err := e.KeysAPI.Delete(context.Background(), key, nil)
	return err
}

func (e *EtcdStorage) Get(key string) (interface{}, error) {
	resp, err := e.KeysAPI.Get(context.Background(), key, nil)
	if err != nil {
		return nil, err
	}
	var obj interface{}
	err = json.Unmarshal([]byte(resp.Node.Value), &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (e *EtcdStorage) Watch(key string, watchHandler func(resp *client.Response)) {
	watcher := e.KeysAPI.Watcher(key, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Errorf("Error watch storage: %v", err)
			break
		}

		watchHandler(res)
	}
}