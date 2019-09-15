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
	"io"
	"io/ioutil"
	"strings"
)

// Source is an input source for our pipe
type Source struct {
	r io.ReadCloser
}

// Read implements the io.Reader interface
func (input *Source) Read(b []byte) (int, error) {
	// do we have anything to read from?
	if input.r == nil {
		return 0, io.EOF
	}

	// read from the input source
	retval, err := input.r.Read(b)

	// what happened?
	if err == io.EOF {
		// we don't want to keep files open when we're done with them,
		// in case we are part of a long-running process
		input.r.Close()
	}

	// all done
	return retval, err
}

// Close tells the input source that we're done reading from it
func (input *Source) Close() error {
	// do we have an input source to close?
	if input.r == nil {
		return nil
	}

	// try and close it
	err := input.r.Close()

	// all done
	return err
}

// String returns all of the data in our buffer as a single (possibly
// multi-line) string
func (input *Source) String() string {
	data, _ := ioutil.ReadAll(input.r)
	return string(data)
}

// NewReader returns a `strings.Reader` for the contents of our buffer
func (input *Source) NewReader() io.Reader {
	return input
}

// ReadLines returns a channel that you can `range` over to get each
// line from our buffer
func (input *Source) ReadLines() <-chan string {
	return NewScanReader(input, bufio.ScanLines)
}

// ReadWords returns a channel that you can `range` over to get each
// word from our buffer
func (input *Source) ReadWords() <-chan string {
	return NewScanReader(input, bufio.ScanWords)
}

// NewSourceFromReader wraps an ordinary io.Reader with our helper methods
func NewSourceFromReader(input io.Reader) *Source {
	return &Source{r: ioutil.NopCloser(input)}
}

// NewSourceFromReadCloser wraps an ordinary io.ReadCloser with our helper
// methods
func NewSourceFromReadCloser(input io.ReadCloser) *Source {
	return &Source{r: input}
}

// NewSourceFromString creates an io.Reader from a string, and wraps it
// with our helper methods
func NewSourceFromString(input string) *Source {
	return NewSourceFromReader(strings.NewReader(input))
}
