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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestNewReaderReturnsReaderForBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var dest Dest
	dest.WriteString("hello world\nhave a nice day")

	expectedResult := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	reader := dest.NewReader()
	scanChn := NewScanReader(reader, bufio.ScanLines)
	var actualResult []string
	for line := range scanChn {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestDestReadLinesIteratesOverBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var dest Dest
	dest.WriteString("hello world\nhave a nice day")

	expectedResult := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for line := range dest.ReadLines() {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestDestReadWordsIteratesOverBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var dest Dest
	dest.WriteString("hello world\nhave a nice day")

	expectedResult := []string{"hello", "world", "have", "a", "nice", "day"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for word := range dest.ReadWords() {
		actualResult = append(actualResult, word)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}
