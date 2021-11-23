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
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTestFile is a helper function. It gives us a file that we can
// test against
func createTestFile(content string) *os.File {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pipe-*")
	if err != nil {
		log.Fatal(err)
	}

	// clean up after ourselves
	defer os.Remove(tmpFile.Name())

	// write the content into our new file
	tmpFile.WriteString(content)
	tmpFile.Seek(0, 0)

	// all done
	return tmpFile
}

// ================================================================
//
// Constructors
//
// ----------------------------------------------------------------

func TestNewFileDestCreatesAFileDest(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"

	// ----------------------------------------------------------------
	// perform the change

	dest := NewFileDest(
		createTestFile(testData),
	)

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, dest)
}

func TestFileDestNewSourceReturnsSourceForUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

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

// ================================================================
//
// Interface compatibility
//
// ----------------------------------------------------------------

func TestFileDestImplementsReadBuffer(t *testing.T) {
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

func TestFileDestImplementsOutput(t *testing.T) {
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

func TestFileDestImplementsInputOutput(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	dest := Dest{}
	var i interface{} = &dest

	// ----------------------------------------------------------------
	// perform the change

	_, ok := i.(InputOutput)

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, ok)
}

// ================================================================
//
// ReadBuffer interface
//
// ----------------------------------------------------------------

func TestFileDestReadFetchesBytesFromTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

	expectedResult := []byte(testData)
	expectedLen := len(testData)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := make([]byte, expectedLen)
	actualLen, err := dest.Read(actualResult)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedLen, actualLen)
	assert.Equal(t, expectedResult, actualResult)
}

func TestFileDestCloseClosesTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

	// ----------------------------------------------------------------
	// perform the change

	dest.Close()

	// we should not be able to read from this buffer any more
	actualResult := make([]byte, len(testData))
	actualLen, err := dest.Read(actualResult)

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, 0, actualLen)
}

func TestFileDestNewReaderReturnsReaderForUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

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

func TestFileDestParseIntReturnsValueOnSuccess(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := " 100 \n"
	dest := NewFileDest(
		createTestFile(testData),
	)

	expectedOutput := 100

	// ----------------------------------------------------------------
	// perform the change

	actualOutput, err := dest.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFileDestReadLinesIteratesOverUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

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

func TestFileDestReadLinesEmptiesTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

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

func TestFileDestReadWordsIteratesOverTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

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

func TestFileDestStringReturnsTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

	expectedOutput := testData

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFileDestStringReturnsEmptyStringIfTheUnderlyingFileHasBeenClosed(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)
	dest.Close()

	expectedOutput := ""

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFileDestStringsReturnsTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(testData),
	)

	expectedOutput := []string{"hello world", "have a nice day"}

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFileDestTrimmedStringReturnsUnderlyingFileWithWhitespaceRemoved(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := " hello world\nhave a nice day\n "
	dest := NewFileDest(
		createTestFile(testData),
	)

	// NOTE: Golang treats trailing '\n' characters as white space :(
	expectedOutput := "hello world\nhave a nice day"

	// ----------------------------------------------------------------
	// perform the change

	actualOutput := dest.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedOutput, actualOutput)
}

// ================================================================
//
// WriteBuffer interface
//
// ----------------------------------------------------------------

func TestFileDestWriteWritesToTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(""),
	)

	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	dest.Write([]byte(testData))

	// ----------------------------------------------------------------
	// test the results

	actualResult := dest.String()

	assert.Equal(t, expectedResult, actualResult)
}

func TestFileDestWriteByteWritesToTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "h"
	dest := NewFileDest(
		createTestFile(""),
	)

	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	dest.WriteByte(testData[0])

	// ----------------------------------------------------------------
	// test the results

	actualResult := dest.String()

	assert.Equal(t, expectedResult, actualResult)
}

func TestFileDestWriteRuneWritesToTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := rune('🙂')
	dest := NewFileDest(
		createTestFile(""),
	)

	expectedResult := string(testData)

	// ----------------------------------------------------------------
	// perform the change

	dest.WriteRune(testData)

	// ----------------------------------------------------------------
	// test the results

	actualResult := dest.String()

	assert.Equal(t, expectedResult, actualResult)
}

func TestFileDestWriteStringWritesToTheUnderlyingFile(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "hello world\nhave a nice day"
	dest := NewFileDest(
		createTestFile(""),
	)

	expectedResult := testData

	// ----------------------------------------------------------------
	// perform the change

	dest.WriteString(testData)

	// ----------------------------------------------------------------
	// test the results

	actualResult := dest.String()

	assert.Equal(t, expectedResult, actualResult)
}
