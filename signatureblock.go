// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package simplecoin

import(
    "fmt"
    "bytes"
    "encoding/binary"
)
/**************************************
 * ISign
 *
 * Interface for RCB Signatures
 *
 * The signature block holds the signatures that validate one of the RCBs.
 * Each signature has an index, so if the RCD is a multisig, you can know
 * how to apply the signatures to the addresses in the RCD.
 **************************************/
type ISignatureBlock interface {
	IBlock
	GetSignatures() ([]ISignature)
	AddSignature(sig ISignature)
}

type SignatureBlock struct{
    ISignatureBlock
    signatures []ISignature
}

var _ ISignatureBlock = (*SignatureBlock)(nil)


func (s SignatureBlock) IsEqual(signatureBlock IBlock) bool {
    
    sb, ok := signatureBlock.(ISignatureBlock)
    
    if !ok { return false }
    
    sigs1 := s.GetSignatures()
    sigs2 := sb.GetSignatures()
    if len(sigs1) != len(sigs2) {return false}
    for i,sig := range sigs1 {
        if !sig.IsEqual(sigs2[i]) {return false}
    }
    
    return true
}

func (s SignatureBlock) AddSignature(sig ISignature) {
    s.signatures = append(s.signatures, sig)
}

func (s SignatureBlock)GetDBHash() IHash {
    return Sha([]byte("SignatureBlock"))
}

func (s SignatureBlock)GetNewInstance() IBlock {
    return new(SignatureBlock)
}

func (s SignatureBlock)GetSignatures() ([]ISignature) {
    if(s.signatures == nil) {
        s.signatures = make([]ISignature,0,1)
    }
    return s.signatures 
}


func (a SignatureBlock) MarshalBinary() ([]byte, error) {
    var out bytes.Buffer
    
    binary.Write(&out, binary.BigEndian, uint16(len(a.signatures))) 
    for _, sig := range a.GetSignatures() {
       
        data, err := sig.MarshalBinary()
        if err != nil { return nil, fmt.Errorf("Signature failed to Marshal in RCD_1") }
        out.Write(data)
    }
    
    return out.Bytes(), nil
}

func (s SignatureBlock) MarshalText() ([]byte, error) {
    var out bytes.Buffer

    out.WriteString("Signature Block: ")
    WriteNumber16(&out, uint16(len(s.signatures)))
    out.WriteString("\n")
    for _, sig := range s.signatures {
        
        out.WriteString(" signature: ")
        txt, err := sig.MarshalText()
        if err != nil { return nil,err }
        out.Write(txt)
        out.WriteString("\n ")
        
    }
    
    return out.Bytes(), nil
}

func (s *SignatureBlock) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
    
    numSignatures, data := binary.BigEndian.Uint16(data[0:2]), data[2:]
    s.signatures = make([]ISignature,numSignatures)
    for i:=uint16(0);i<numSignatures;i++ {
        s.signatures[i] = new(Signature)
        data,err = s.signatures[i].UnmarshalBinaryData(data)
        if err != nil {
            return nil, fmt.Errorf("Failure to unmarshal Signature")
        }
    }
    
    return data, nil
}