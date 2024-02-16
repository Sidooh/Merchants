package cache

import (
	"github.com/jellydator/ttlcache/v3"
	"os"
	"testing"
	"time"
)

var instance ICache[string, string]

func TestMain(m *testing.M) {

	instance = New[string, string]()

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	item := instance.Set("key", "value", ttlcache.DefaultTTL)

	if item == nil {
		t.Errorf("item = nil; want item")
	}

	item = instance.Set("key", "value", 5*time.Second)
}

func TestGet(t *testing.T) {
	nonExistentItem := instance.Get("key1")

	if nonExistentItem != nil {
		t.Errorf("nonExistentItem is %v; want nil", nonExistentItem)
	}
}

func TestExpiry(t *testing.T) {
	item := instance.Set("key", "value", 1*time.Second)

	if item == nil {
		t.Errorf("item = nil; want item")
	}

	exists := instance.Get("key")

	if exists == nil {
		t.Errorf("exists = %v; want item", exists)
	}

	time.Sleep(1 * time.Second)

	nonExistentItem := instance.Get("key")

	if nonExistentItem != nil {
		t.Errorf("nonExistentItem = %v; want nil", nonExistentItem)
	}
}
