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
	"os"
	"testing"

	envish "github.com/ganbarodigital/go_envish"
	"github.com/stretchr/testify/assert"
)

func TestNewPipeCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe()
	actualResult := pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe()
	actualResult := pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe()
	actualResult := pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithStatusOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe()
	actualResult := pipe.StatusCode

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithNilError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var expectedResult error

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe()
	actualResult := pipe.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeNextMakesStdoutTheNextStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world"

	pipe := NewPipe()
	pipe.Stdout.WriteString("hello world")

	// prove that pipe.Stdin is empty before we call pipe.Next()
	assert.Equal(t, pipe.Stdin.String(), "")

	// ----------------------------------------------------------------
	// perform the change

	pipe.Next()
	actualResult := pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeNextUpdatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	pipe := NewPipe()
	pipe.Stdout.WriteString("hello world")

	// make sure that Stdout does have some content first
	assert.Equal(t, pipe.Stdout.String(), "hello world")
	// make sure that a call to Stdout.String() doesn't empty the buffer!
	assert.Equal(t, pipe.Stdout.String(), "hello world")

	// ----------------------------------------------------------------
	// perform the change

	pipe.Next()
	actualResult := pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeNextUpdatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	pipe := NewPipe()
	pipe.Stderr.WriteString("hello world")

	// make sure that Stderr does have some content first
	assert.Equal(t, pipe.Stderr.String(), "hello world")
	// make sure that a call to Stderr.String() doesn't empty the buffer!
	assert.Equal(t, pipe.Stderr.String(), "hello world")

	// ----------------------------------------------------------------
	// perform the change

	pipe.Next()
	actualResult := pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeNextCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.Next()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeNextCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.Next()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.Reset()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.Reset()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetEmptiesStdinStdoutStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataIn := "this is stdin"
	testDataOut := "this is stdout"
	testDataErr := "this is stderr"

	// we need to start with a pipe that has data
	pipe := NewPipe()
	pipe.Stdout.WriteString(testDataIn)
	pipe.Next()
	pipe.Stdout.WriteString(testDataOut)
	pipe.Stderr.WriteString(testDataErr)

	// normally, I'd use assert.Equal() to prove that the pipe has data
	// if we did that here, the reads would empty the pipe, making the
	// rest of the test invalid

	// ----------------------------------------------------------------
	// perform the change

	pipe.Reset()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, pipe.Stdin.String())
	assert.Empty(t, pipe.Stdout.String())
	assert.Empty(t, pipe.Stderr.String())
}

func TestPipeDrainCopiesStdinToStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"

	pipe := NewPipe()
	pipe.Stdout.WriteString(expectedResult)
	pipe.Next()

	assert.Equal(t, pipe.Stdout.String(), "")

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdin()
	actualResult := pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeDrainStdinCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeDrainStdinCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeErrorCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, actualResult)
}

func TestPipeExpandCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	expectedResult := os.Getenv("HOME")

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Expand("${HOME}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeExpandCopesWithEmptyStruct(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	expectedResult := os.Getenv("HOME")

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Expand("${HOME}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeExpandUsesTemporaryEnvironmentIfWeHaveOne(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "this is not a real HOME folder"

	pipe := NewPipe()
	pipe.Env = envish.NewEnv()
	pipe.Env.Setenv("HOME", expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Expand("${HOME}")

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeRunCommandCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	expectedResult := StatusNotOkay
	op := func(p *Pipe) (int, error) {
		return expectedResult, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	// as long as it doesn't crash, the test has passed
}

func TestPipeRunCommandUpdatesStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()

	expectedResult := StatusNotOkay
	op := func(p *Pipe) (int, error) {
		return expectedResult, errors.New("status not okay")
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, pipe.StatusCode)
}

func TestPipeRunCommandUpdatesErr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()

	expectedResult := errors.New("status not okay")
	op := func(p *Pipe) (int, error) {
		return StatusNotOkay, expectedResult
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, pipe.Err)
}

func TestPipeRunCommandSetsErrIfStatusCodeNotOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()

	op := func(p *Pipe) (int, error) {
		return StatusNotOkay, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, pipe.Err)
	assert.Error(t, pipe.Err)
	_, ok := pipe.Err.(ErrNonZeroStatusCode)
	assert.True(t, ok)
}
