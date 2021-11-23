// pipe is a library to help you write UNIX-like pipelines of operations
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package pipe

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

// Dest is an output source for our pipe
type Dest struct {
	bytes.Buffer
}

// NewDest creates a new output buffer
func NewDest() *Dest {
	retval := Dest{}

	// all done
	return &retval
}

// Close tells the input source that we're done reading from it
//
// As our input source is a memory-backed buffer, this is a no-op
func (d *Dest) Close() error {
	return nil
}

// NewReader returns an `io.Reader` for the contents of our buffer
func (d *Dest) NewReader() io.Reader {
	return d
}

// NewSource returns a `Source` for reading the contents of our buffer
func (d *Dest) NewSource() *Source {
	return NewSourceFromString(d.String())
}

// ParseInt returns the data in our buffer as an integer.
//
// If the buffer contains anything other than a valid number, an error
// is returned.
func (d *Dest) ParseInt() (int, error) {
	text := d.TrimmedString()
	return strconv.Atoi(text)
}

// ReadLines returns a channel that you can `range` over to get each
// line from our buffer
func (d *Dest) ReadLines() <-chan string {
	return NewScanReader(d.NewReader(), bufio.ScanLines)
}

// ReadWords returns a channel that you can `range` over to get each
// word from our buffer
func (d *Dest) ReadWords() <-chan string {
	return NewScanReader(d.NewReader(), bufio.ScanWords)
}

// Strings returns all of the data in our buffer as an array of
// strings, one line per array entry
func (d *Dest) Strings() []string {
	retval := []string{}
	for line := range d.ReadLines() {
		retval = append(retval, line)
	}

	return retval
}

// TrimmedString returns all of the data in our buffer as a string,
// with any leading or trailing whitespace removed.
func (d *Dest) TrimmedString() string {
	return strings.TrimSpace(d.String())
}
