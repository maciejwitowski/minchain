package chain

import (
	"crypto/ecdsa"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
}

func NewWallet(pk *ecdsa.PrivateKey) *Wallet {
	return &Wallet{privateKey: pk}
}

func (w *Wallet) SignedTransaction(message string) (*Tx, error) {
	from := crypto.PubkeyToAddress(w.privateKey.PublicKey)
	digest := crypto.Keccak256([]byte(message))

	expectSig, err := crypto.Sign(digest, w.privateKey)
	if err != nil {
		return nil, err
	}

	return &Tx{
		From:      from.String(),
		Data:      message,
		Signature: expectSig,
	}, nil
}
