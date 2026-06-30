package utils

import (
	"github.com/ethereum/go-ethereum/rlp"
)

// EIP-7928 Block Access List (BAL) types for RLP decoding.
// These use variable-length []byte for uint256 fields (balance, storage key/value)
// to match the EIP-7928 spec, which uses standard RLP variable-length encoding.
type (
	BALStorageRoot struct {
		Value []byte
	}
	BALStorageWrite struct {
		TxIdx      uint16
		ValueAfter []byte
	}
	BALSlotWrites struct {
		Slot     []byte
		Accesses []BALStorageWrite
	}
	BALBalanceChange struct {
		TxIdx   uint16
		Balance []byte
	}
	BALNonceChange struct {
		TxIdx uint16
		Nonce uint64
	}
	BALCodeChange struct {
		TxIndex uint16
		Code    []byte
	}
	BALAccountAccess struct {
		Address        [20]byte
		StorageWrites  []BALSlotWrites
		StorageReads   [][]byte
		BalanceChanges []BALBalanceChange
		NonceChanges   []BALNonceChange
		CodeChanges    []BALCodeChange
		StorageRoot    *BALStorageRoot
	}
)

// DecodeRLP supports both the original EIP-7928 account tuple and the
// Heza/Bogota extension that appends an optional storage_root byte string.
func (a *BALAccountAccess) DecodeRLP(s *rlp.Stream) error {
	if _, err := s.List(); err != nil {
		return err
	}
	if err := s.Decode(&a.Address); err != nil {
		return err
	}
	if err := s.Decode(&a.StorageWrites); err != nil {
		return err
	}
	if err := s.Decode(&a.StorageReads); err != nil {
		return err
	}
	if err := s.Decode(&a.BalanceChanges); err != nil {
		return err
	}
	if err := s.Decode(&a.NonceChanges); err != nil {
		return err
	}
	if err := s.Decode(&a.CodeChanges); err != nil {
		return err
	}
	if s.MoreDataInList() {
		storageRoot, err := s.Bytes()
		if err != nil {
			return err
		}
		a.StorageRoot = &BALStorageRoot{Value: storageRoot}
	}
	return s.ListEnd()
}

// DecodeBlockAccessList decodes an EIP-7928 block access list from RLP bytes.
func DecodeBlockAccessList(data []byte) ([]BALAccountAccess, error) {
	var accesses []BALAccountAccess
	if err := rlp.DecodeBytes(data, &accesses); err != nil {
		return nil, err
	}
	return accesses, nil
}
