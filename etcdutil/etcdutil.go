package etcdutil

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/sssvip/goutil/logutil"
	"time"
)

type ETCDWrapper struct {
	client *clientv3.Client
}

func NewETCDWrapper(endPoint string) *ETCDWrapper {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endPoint},
		DialTimeout: 5 * time.Second})
	if err != nil {
		logutil.Error.Println(err.Error())
		return nil
	}
	return &ETCDWrapper{client: cli}
}

func (w *ETCDWrapper) Client() (client *clientv3.Client) {
	return w.client
}

func (w *ETCDWrapper) Put(key, value string) {
	_, err := w.client.Put(context.Background(), key, value)
	if err != nil {
		logutil.Error.Println(err)
	}
	return
}

func (w *ETCDWrapper) Get(key string) (value string) {
	v, e := w.client.Get(context.Background(), key, clientv3.WithFromKey())
	if e != nil {
		logutil.Error.Println(e)
		return ""
	}
	for _, kv := range v.Kvs {
		return string(kv.Value)
	}
	return
}

func (w *ETCDWrapper) ListenChanged(key string, execFunc func(value string)) {
	watchCh := w.client.Watch(context.TODO(), key, clientv3.WithKeysOnly())
	for res := range watchCh {
		for _, event := range res.Events {
			execFunc(string(event.Kv.Value))
		}
	}
}

func (w *ETCDWrapper) Close() {
	if w.client != nil {
		e := w.client.Close()
		if e != nil {
			logutil.Error.Println(e)
		}
	}
}
