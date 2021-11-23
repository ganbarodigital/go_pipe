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

	dest := NewDest()
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

func TestDestNewSourceReturnsSourceForBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
	dest.WriteString("hello world\nhave a nice day")

	expectedResult := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	source := dest.NewSource()
	scanChn := NewScanReader(source, bufio.ScanLines)
	var actualResult []string
	for line := range scanChn {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestDestParseIntReturnsValueOnSuccess(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := " 100 \n"
	expectedOutput := 100

	dest := NewDest()
	dest.WriteString(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualOutput, err := dest.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestDestReadLinesIteratesOverBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
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

func TestDestReadLinesEmptiesTheBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
	dest.WriteString("hello world\nhave a nice day")

	expectedResult := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for line := range dest.ReadLines() {
		actualResult = append(actualResult, line)
	}

	extraOutput := []string{}
	for line := range dest.ReadLines() {
		extraOutput = append(extraOutput, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Empty(t, extraOutput)
}

func TestDestReadWordsIteratesOverBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
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

func TestDestStringReturnsBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedOutput := "hello world\n"
	dest := NewDest()
	dest.WriteString(expectedOutput)

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestDestStringsReturnsBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
	dest.WriteString("hello world\nhave a nice day\n")
	expectedOutput := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestDestTrimmedStringReturnsBufferWithWhitespaceRemoved(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := NewDest()
	dest.WriteString(" hello world\n")
	expectedOutput := "hello world"

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestDestImplementsReadBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := Dest{}
	var i interface{} = &dest

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Input)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestDestImplementsOutput(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := Dest{}
	var i interface{} = &dest

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(Output)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

func TestDestImplementsReadWriteBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := Dest{}
	var i interface{} = &dest

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(ReadWriteBuffer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}
