// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package simplecoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	//    "github.com/agl/ed25519"
)

/**************************************
 * ISign
 *
 * Interface for RCB Signatures
 *
 * Data structure to support signatures, including multisig.
 **************************************/

type ISignature interface {
	IBlock
	GetIndex() int       // Used with multisig to know where to apply the signature
	SetIndex(int)        // Set the index
    SetSignature(i int, sig []byte) error // Set or update the signature
    GetSignature(i int) ([SIGNATURE_LENGTH]byte,error)
}

// The default signature doesn't care about indexing.  We will extend this
// signature for multisig
type Signature struct {
	ISignature
	signature [SIGNATURE_LENGTH]byte // The signature
}

var _ ISignature = (*Signature)(nil)

func (w1 Signature)GetDBHash() IHash {
    return Sha([]byte("Signature"))
}

func (w1 Signature)GetNewInstance() IBlock {
    return new(Signature)
}

// Checks that the signatures are the same.  
func (s1 Signature) IsEqual(sig IBlock) bool {
	s2, ok := sig.(*Signature)
	if !ok || // Not the right kind of IBlock
		s1.signature != s2.signature { // Not the right rcd
		return false
	}
	return true
}

// Index is ignored.  We only have one signature
func (s *Signature) SetSignature(i int, sig []byte) error {
	if len(sig) != SIGNATURE_LENGTH {
		return fmt.Errorf("Bad signature.  Should not happen")
	}
	copy(s.signature[:], sig)
    return nil
}

func (s *Signature) GetSignature(i int) ([SIGNATURE_LENGTH]byte, error) {
    return s.signature,nil
}


func (s Signature) MarshalBinary() ([]byte, error) {
	var out bytes.Buffer

	out.Write(s.signature[:])

	return out.Bytes(), nil
}

func (s Signature) MarshalText() ([]byte, error) {
	var out bytes.Buffer

	out.WriteString(" signature: ")
	out.WriteString(hex.EncodeToString(s.signature[:]))
	out.WriteString("\n")

	return out.Bytes(), nil
}

func (s *Signature) UnmarshalBinaryData(data []byte) ([]byte, error) {
    copy(s.signature[:], data[:SIGNATURE_LENGTH])
    return data[SIGNATURE_LENGTH:], nil 
}

