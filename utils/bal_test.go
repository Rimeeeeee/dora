package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

type testBALAccountWithoutStorageRoot struct {
	Address        [20]byte
	StorageWrites  []BALSlotWrites
	StorageReads   [][]byte
	BalanceChanges []BALBalanceChange
	NonceChanges   []BALNonceChange
	CodeChanges    []BALCodeChange
}

type testBALAccountWithStorageRoot struct {
	Address        [20]byte
	StorageWrites  []BALSlotWrites
	StorageReads   [][]byte
	BalanceChanges []BALBalanceChange
	NonceChanges   []BALNonceChange
	CodeChanges    []BALCodeChange
	StorageRoot    []byte
}

func TestDecodeBlockAccessListStorageRoot(t *testing.T) {
	var address [20]byte
	address[19] = 1

	withoutRoot, err := rlp.EncodeToBytes([]testBALAccountWithoutStorageRoot{{
		Address: address,
	}})
	if err != nil {
		t.Fatalf("encode without root: %v", err)
	}
	accesses, err := DecodeBlockAccessList(withoutRoot)
	if err != nil {
		t.Fatalf("decode without root: %v", err)
	}
	if len(accesses) != 1 {
		t.Fatalf("decoded %d accesses, want 1", len(accesses))
	}
	if accesses[0].StorageRoot != nil {
		t.Fatalf("storage root = %x, want omitted", accesses[0].StorageRoot.Value)
	}

	emptyRoot, err := rlp.EncodeToBytes([]testBALAccountWithStorageRoot{{
		Address:     address,
		StorageRoot: []byte{},
	}})
	if err != nil {
		t.Fatalf("encode empty root: %v", err)
	}
	accesses, err = DecodeBlockAccessList(emptyRoot)
	if err != nil {
		t.Fatalf("decode empty root: %v", err)
	}
	if accesses[0].StorageRoot == nil {
		t.Fatal("storage root omitted, want present empty marker")
	}
	if len(accesses[0].StorageRoot.Value) != 0 {
		t.Fatalf("storage root = %x, want empty marker", accesses[0].StorageRoot.Value)
	}

	root := make([]byte, 32)
	root[31] = 2
	withRoot, err := rlp.EncodeToBytes([]testBALAccountWithStorageRoot{{
		Address:     address,
		StorageRoot: root,
	}})
	if err != nil {
		t.Fatalf("encode concrete root: %v", err)
	}
	accesses, err = DecodeBlockAccessList(withRoot)
	if err != nil {
		t.Fatalf("decode concrete root: %v", err)
	}
	if accesses[0].StorageRoot == nil {
		t.Fatal("storage root omitted, want concrete root")
	}
	if got := accesses[0].StorageRoot.Value; string(got) != string(root) {
		t.Fatalf("storage root = %x, want %x", got, root)
	}
}
