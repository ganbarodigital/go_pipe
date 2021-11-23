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
	"os"
	"strconv"
	"strings"
)

// FileDest is an output destination for Pipe
type FileDest struct {
	f *os.File
}

// ===========================================================================
//
// Constructors
//
// ---------------------------------------------------------------------------

// NewFileDest creates a new output destination that reads from / writes to
// and underlying file
func NewFileDest(f *os.File) *FileDest {
	retval := FileDest{
		f: f,
	}

	// all done
	return &retval
}

// NewSource returns a `Source` that contains a copy of our data
func (d *FileDest) NewSource() *Source {
	return NewSourceFromString(d.String())
}

// ===========================================================================
//
// ReadBuffer interface
//
// ---------------------------------------------------------------------------

// Read reads up to len(b) bytes from the File. It returns the number of
// bytes read and any error encountered. At the end of file, Read returns
// (0, io.EOF).
func (d *FileDest) Read(b []byte) (int, error) {
	return d.f.Read(b)
}

// Close tells the underlying file that we're done reading from it.
//
// All subsequent attempts to read from us will result in errors and
// empty data sets.
func (d *FileDest) Close() error {
	return d.f.Close()
}

// NewReader returns an `io.Reader` for the contents of our underlying
// file
func (d *FileDest) NewReader() io.Reader {
	return d
}

// ParseInt returns the data in our underlying file as an integer.
//
// If the file contains anything other than a valid number, an error
// is returned.
func (d *FileDest) ParseInt() (int, error) {
	text := d.TrimmedString()
	return strconv.Atoi(text)
}

// ReadLines returns a channel that you can `range` over to get each
// line from our underlying file
func (d *FileDest) ReadLines() <-chan string {
	return NewScanReader(d.NewReader(), bufio.ScanLines)
}

// ReadWords returns a channel that you can `range` over to get each
// word from our underlying file
func (d *FileDest) ReadWords() <-chan string {
	return NewScanReader(d.NewReader(), bufio.ScanWords)
}

// String returns all of the data in our underlying file as a single
// (possibly multi-line) string
func (d *FileDest) String() string {
	// make sure we are at the start of the buffer
	d.f.Seek(0, 0)

	retval, err := ioutil.ReadAll(d.f)
	if err != nil {
		return ""
	}

	return string(retval)
}

// Strings returns all of the data in our underlying file as an array of
// strings, one line per array entry
func (d *FileDest) Strings() []string {
	retval := []string{}
	for line := range d.ReadLines() {
		retval = append(retval, line)
	}

	return retval
}

// TrimmedString returns all of the data in our underlying file as a string,
// with any leading or trailing whitespace removed.
func (d *FileDest) TrimmedString() string {
	return strings.TrimSpace(d.String())
}

// ===========================================================================
//
// WriteBuffer interface
//
// ---------------------------------------------------------------------------

// Write writes len(p) bytes from p to the underlying file. It returns
// the number of bytes written from p (0 <= n <= len(p)) and any error
// encountered that caused the write to stop early.
//
// Write returns a non-nil error if it returns n < len(p).
//
// Write does not modify the slice data, even temporarily.
func (d *FileDest) Write(p []byte) (int, error) {
	return d.f.Write(p)
}

// WriteByte writes a single byte to the underlying file. It returns any
// error encountered that caused the write to fail.
func (d *FileDest) WriteByte(c byte) error {
	_, err := d.f.Write([]byte{c})
	return err
}

// WriteRune writes a single rune (a unicode character) to the underlying
// file. It returns the number of types written, and any error encountered
// that caused the write to file.
func (d *FileDest) WriteRune(r rune) (int, error) {
	return d.WriteString(string(r))
}

// WriteString writes the contents of the string to the file. It returns
// the number of bytes written to the file, and any error encountered
// that caused the write to stop early.
func (d *FileDest) WriteString(s string) (int, error) {
	return d.f.WriteString(s)
}
