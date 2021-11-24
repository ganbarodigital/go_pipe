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
	"testing"

	envish "github.com/ganbarodigital/go_envish/v3"
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
	actualResult := pipe.StatusCode()

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

func TestNewPipeAppliesAnyOptionsWePassIn(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStatusCode := 100
	op1 := func(p *Pipe) {
		p.statusCode = expectedStatusCode
	}
	op2 := func(p *Pipe) {
		p.Env = envish.NewLocalEnv()
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe := NewPipe(op1, op2)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, pipe.StatusCode(), expectedStatusCode)
	assert.NotNil(t, pipe.Env)
}

func TestPipeDrainStdinToStdoutCopiesStdinToStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"

	pipe := NewPipe()
	pipe.SetStdinFromString(expectedResult)

	assert.Equal(t, pipe.Stdout.String(), "")

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdinToStdout()
	actualResult := pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeDrainStdinToStdoutCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeDrainStdinToStdoutCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeDrainStdinToStdoutCreatesNewStdoutIfNecessary(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()
	pipe.Stdout = nil

	// ----------------------------------------------------------------
	// perform the change

	pipe.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, pipe.Stdout)
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

func TestPipeOkayCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeOkayCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeOkayReturnsFalseIfTheLastCommandFailed(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()
	expectedResult := false

	op1 := func(p *Pipe) (int, error) {
		return StatusOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return 100, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op1)
	pipe.RunCommand(op2)
	actualResult := pipe.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeResetBuffersCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetBuffersCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetBuffersEmptiesStdinStdoutStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataIn := "this is stdin"
	testDataOut := "this is stdout"
	testDataErr := "this is stderr"

	// we need to start with a pipe that has data
	pipe := NewPipe()
	pipe.SetStdinFromString(testDataIn)
	pipe.Stdout.WriteString(testDataOut)
	pipe.Stderr.WriteString(testDataErr)

	// normally, I'd use assert.Equal() to prove that the pipe has data
	// if we did that here, the reads would empty the pipe, making the
	// rest of the test invalid

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, pipe.Stdin.String())
	assert.Empty(t, pipe.Stdout.String())
	assert.Empty(t, pipe.Stderr.String())
}

func TestPipeResetErrorCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetError()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetErrorCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetError()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetErrorSetsStatusCodeToStatusOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		return StatusNotOkay, nil
	}
	pipe := NewPipe()
	pipe.RunCommand(op1)

	assert.Equal(t, StatusNotOkay, pipe.StatusCode())
	assert.Error(t, pipe.Error())

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, StatusOkay, pipe.StatusCode())
}

func TestPipeResetErrorSetsErrorToNil(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		return StatusNotOkay, nil
	}
	pipe := NewPipe()
	pipe.RunCommand(op1)

	assert.Equal(t, StatusNotOkay, pipe.StatusCode())
	assert.Error(t, pipe.Error())

	// ----------------------------------------------------------------
	// perform the change

	pipe.ResetError()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipe.Error())
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

	assert.Equal(t, expectedResult, pipe.StatusCode())
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

	assert.Equal(t, expectedResult, pipe.Error())
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

	err := pipe.Error()
	assert.NotNil(t, err)
	assert.Error(t, err)
	_, ok := err.(ErrNonZeroStatusCode)
	assert.True(t, ok)
}

func TestPipeSetNewStdinCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStdinCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetStdinFromStringCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetStdinFromString("")

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetStdinFromStringCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetStdinFromString("")

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStdoutCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStdoutCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStderrCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStderr()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStderrCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe

	// ----------------------------------------------------------------
	// perform the change

	pipe.SetNewStderr()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeStatusCodeCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe
	expectedResult := StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusCodeCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe
	expectedResult := StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	actualResult := pipe.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusCodeReturnsTheLastCommandsStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()
	expectedResult := 100

	op1 := func(p *Pipe) (int, error) {
		return StatusOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return expectedResult, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op1)
	pipe.RunCommand(op2)
	actualResult := pipe.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusErrorCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe *Pipe
	expectedResult := StatusOkay
	var expectedErr error = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, actualErr := pipe.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}

func TestPipeStatusErrorCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipe Pipe
	expectedResult := StatusOkay
	var expectedErr error = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, actualErr := pipe.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}

func TestPipeStatusErrorReturnsTheLastCommandsStatusCodeAndError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	pipe := NewPipe()
	expectedResult := 100
	expectedErr := errors.New("this is an error")

	op1 := func(p *Pipe) (int, error) {
		return StatusOkay, nil
	}
	op2 := func(p *Pipe) (int, error) {
		return expectedResult, expectedErr
	}

	// ----------------------------------------------------------------
	// perform the change

	pipe.RunCommand(op1)
	pipe.RunCommand(op2)
	actualResult, actualErr := pipe.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}
