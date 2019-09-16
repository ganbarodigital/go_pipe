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
	"errors"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceReadReturnsZeroBytesReadOnNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var source Source

	// ----------------------------------------------------------------
	// perform the change

	var buf []byte
	bytesRead, _ := source.Read(buf)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, bytesRead)
}

func TestSourceReadReturnsEOFOnNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var source Source

	// ----------------------------------------------------------------
	// perform the change

	var buf []byte
	_, err := source.Read(buf)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, io.EOF, err)
}

func TestSourceCloseReturnsNilOnNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var source Source

	// ----------------------------------------------------------------
	// perform the change

	err := source.Close()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestSourceCloseReturnsWrappedErrOnClose(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	_, filename, _, _ := runtime.Caller(0)

	fh, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	source := Source{fh}

	// ----------------------------------------------------------------
	// perform the change

	err1 := source.Close()
	err2 := source.Close()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err1)

	_, ok := err2.(*os.PathError)
	if !ok {
		t.Error("second call to Source.Close() did not return the wrapped error")
	}
}

func TestNewReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	source := Source{ioutil.NopCloser(strings.NewReader(""))}

	// ----------------------------------------------------------------
	// perform the change

	reader := source.NewReader()

	// ----------------------------------------------------------------
	// test the results

	_, ok := reader.(io.Reader)
	if !ok {
		t.Errorf("Source.NewReader() did not return an io.Reader-compatible object\n")
	}
}

func TestReadLines(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	source := Source{ioutil.NopCloser(strings.NewReader("one\ntwo\nthree\n"))}
	expectedOutput := []string{"one", "two", "three"}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := []string{}
	for line := range source.ReadLines() {
		actualOutput = append(actualOutput, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput, "ReadLines() did not produce what we expect")
}

func TestSourceStringReturnsBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedOutput := "hello world\n"
	source := NewSourceFromString(expectedOutput)

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := source.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestSourceStringsReturnsBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	source := NewSourceFromString("hello world\nhave a nice day\n")
	expectedOutput := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := source.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestSourceTrimmedStringReturnsBufferWithWhitespaceRemoved(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := " hello world\n"
	expectedOutput := "hello world"
	source := NewSourceFromString(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := source.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestReadWords(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	source := Source{ioutil.NopCloser(strings.NewReader("one two\nthree\nfour five six\n\nseven\n    eight\n"))}
	expectedOutput := []string{"one", "two", "three", "four", "five", "six", "seven", "eight"}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := []string{}
	for word := range source.ReadWords() {
		actualOutput = append(actualOutput, word)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput, "ReadWords() did not produce what we expect")
}

func TestAutoClose(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	_, filename, _, _ := runtime.Caller(0)

	fh, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	source := Source{fh}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := []string{}
	for lines := range source.ReadLines() {
		actualOutput = append(actualOutput, lines)
	}

	// ----------------------------------------------------------------
	// test the results

	var b []byte
	_, err = fh.Read(b)

	if !errors.Is(err, os.ErrClosed) {
		t.Error("Source did not auto-close the reader\n")
	}
}

func TestNewSourceFromReaderWrapsGivenReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	buf := "hello world\n"
	reader := strings.NewReader(buf)
	source := NewSourceFromReader(reader)
	expectedResult := []string{"hello world"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for line := range source.ReadLines() {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSourceFromReadCloserWrapsGivenReader(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	buf := "hello world\n"
	reader := strings.NewReader(buf)
	readCloser := ioutil.NopCloser(reader)
	source := NewSourceFromReadCloser(readCloser)
	expectedResult := []string{"hello world"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for line := range source.ReadLines() {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewSourceFromStringWrapsGivenString(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	buf := "hello world\n"
	source := NewSourceFromString(buf)
	expectedResult := []string{"hello world"}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []string
	for line := range source.ReadLines() {
		actualResult = append(actualResult, line)
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestSourceImplementsReadBuffer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	source := Source{}
	var i interface{} = &source

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(ReadBuffer)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}
