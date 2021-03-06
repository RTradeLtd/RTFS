package rtfs_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/RTradeLtd/krab/v4"
	"github.com/RTradeLtd/rtfs/v2"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	ci "github.com/libp2p/go-libp2p-core/crypto"
)

func TestKeystoreManager(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()
	kb, err := krab.NewKeystore(dssync.MutexWrap(datastore.NewMapDatastore()), "password123")
	if err != nil {
		t.Fatal(err)
	}
	km, err := rtfs.NewKeystoreManager(kb)
	if err != nil {
		t.Fatal(err)
	}

	keys, err := km.ListKeyIdentifiers()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range keys {
		_, err := km.GetPrivateKeyByName(v)
		if err != nil {
			t.Error(err)
			continue
		}
		present, err := km.CheckIfKeyExists(v)
		if err != nil {
			t.Error(err)
			continue
		}
		if !present {
			t.Error("key not present when it should be")
			continue
		}
	}

	present, err := km.CheckIfKeyExists("thiskeyshouldreallynotexistwithsucharandomname")
	if err == nil {
		t.Fatal("no error, key was found when it shouldn't have been")
	}
	if present {
		t.Fatal("key found when it should'nt have been")
	}
	// DO NOT USE 1024 IN PRODUCTION, >= 2048
	pk, _, err := ci.GenerateKeyPair(ci.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	hexed := hex.EncodeToString(b)
	err = km.SavePrivateKey(hexed, pk)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetKey(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()

	var (
		k1 = "b6ec4a647a7738ef8eea3b21777ecf41630d6d0ac79dc36739d81e927f910a65"
		k2 = "test1"
	)
	kb, err := krab.NewKeystore(dssync.MutexWrap(datastore.NewMapDatastore()), "password123")
	if err != nil {
		t.Fatal(err)
	}
	km, err := rtfs.NewKeystoreManager(kb)
	if err != nil {
		t.Fatal(err)
	}
	// Create keys
	km.CreateAndSaveKey(k1, 1, 1)
	km.CreateAndSaveKey(k2, 1, 1)
	// Check if keys exist
	present, err := km.CheckIfKeyExists(k1)
	if err != nil {
		t.Fatal(err)
	}
	if !present {
		t.Error("key not present when it should be")
	}
	pk1, err := km.GetPrivateKeyByName(k1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", pk1.GetPublic())

	present, err = km.CheckIfKeyExists(k2)
	if err != nil {
		t.Fatal(err)
	}

	if !present {
		t.Fatal("key not present when it should be")
	}

	pk2, err := km.GetPrivateKeyByName(k2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", pk2.GetPublic())
}
