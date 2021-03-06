// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package simplecoin

import (
	// "fmt"
	"encoding"
)

type IBlock interface {
	encoding.BinaryMarshaler   // Easy to support this, just drop the slice.
	encoding.BinaryUnmarshaler // And once in Binary, it must come back.
	encoding.TextMarshaler     // Using this mostly for debugging

	// We need the progress through the slice, so we really can't use the stock spec
	// for the UnmarshalBinary() method from encode.  We define our own method that
	// makes the code easier to read and way more efficent.
	UnmarshalBinaryData(data []byte) ([]byte, error) 
    
	IsEqual(IBlock)  bool   // Check if this block is the same as itself
	GetDBHash()      IHash  // Identifies the object 
	GetNewInstance() IBlock // Get a new instance of this object
}
