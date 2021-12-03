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

package pipe_test

import (
	"errors"
	"testing"

	"github.com/ganbarodigital/go-ioextra/v2"
	envish "github.com/ganbarodigital/go_envish/v3"
	pipe "github.com/ganbarodigital/go_pipe/v6"
	"github.com/stretchr/testify/assert"
)

func TestNewPipeCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()
	actualResult := unit.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithEmptyStdinStack(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()

	// ----------------------------------------------------------------
	// test the results

	assert.Zero(t, unit.StdinStackLen())
}

func TestNewPipeCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()
	actualResult := unit.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithEmptyStdoutStack(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()

	// ----------------------------------------------------------------
	// test the results

	assert.Zero(t, unit.StdoutStackLen())
}

func TestNewPipeCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()
	actualResult := unit.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeCreatesPipeWithEmptyStderrStack(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()

	// ----------------------------------------------------------------
	// test the results

	assert.Zero(t, unit.StderrStackLen())
}

func TestNewPipeCreatesPipeWithStatusOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := pipe.StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe()
	actualResult := unit.StatusCode()

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

	unit := pipe.NewPipe()
	actualResult := unit.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipeAppliesAnyOptionsWePassIn(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStatusCode := 100
	op1 := func(p *pipe.Pipe) {
		// use a helper method
		pipe.SetStatusCode(p, expectedStatusCode)
	}
	op2 := func(p *pipe.Pipe) {
		p.Env = envish.NewLocalEnv()
	}

	// ----------------------------------------------------------------
	// perform the change

	unit := pipe.NewPipe(op1, op2)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, unit.StatusCode(), expectedStatusCode)
	assert.NotNil(t, unit.Env)
}

func TestPipeDrainStdinToStdoutCopiesStdinToStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"

	unit := pipe.NewPipe()
	unit.SetStdinFromString(expectedResult)

	assert.Equal(t, unit.Stdout.String(), "")

	// ----------------------------------------------------------------
	// perform the change

	unit.DrainStdinToStdout()
	actualResult := unit.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeDrainStdinToStdoutCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeDrainStdinToStdoutCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeDrainStdinToStdoutCreatesNewStdoutIfNecessary(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()
	unit.Stdout = nil

	// ----------------------------------------------------------------
	// perform the change

	unit.DrainStdinToStdout()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, unit.Stdout)
}

func TestPipeErrorCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	actualResult := unit.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, actualResult)
}

func TestPipeOkayCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := unit.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeOkayCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe
	expectedResult := true

	// ----------------------------------------------------------------
	// perform the change

	actualResult := unit.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeOkayReturnsFalseIfTheLastCommandFailed(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()
	expectedResult := false

	op1 := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusOkay, nil
	}
	op2 := func(p *pipe.Pipe) (int, error) {
		return 100, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op1)
	unit.RunCommand(op2)
	actualResult := unit.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeResetBuffersCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetBuffersCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetBuffers()

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
	unit := pipe.NewPipe()
	unit.SetStdinFromString(testDataIn)
	unit.Stdout.WriteString(testDataOut)
	unit.Stderr.WriteString(testDataErr)

	// normally, I'd use assert.Equal() to prove that the pipe has data
	// if we did that here, the reads would empty the pipe, making the
	// rest of the test invalid

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, unit.Stdin.String())
	assert.Empty(t, unit.Stdout.String())
	assert.Empty(t, unit.Stderr.String())
}

func TestPipeResetBuffersEmptiesInternalStdinStdoutStderrStacks(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testDataIn := "this is stdin"
	testDataOut := "this is stdout"
	testDataErr := "this is stderr"

	// we need to start with a pipe that has data
	unit := pipe.NewPipe()
	unit.SetStdinFromString(testDataIn)
	unit.Stdout.WriteString(testDataOut)
	unit.Stderr.WriteString(testDataErr)

	// now, we need to put these onto the internal stacks
	unit.PushStdin(ioextra.NewTextBuffer())
	unit.PushStdout(ioextra.NewTextBuffer())
	unit.PushStderr(ioextra.NewTextBuffer())

	// make sure that all of the stacks now have something on them
	assert.Equal(t, 1, unit.StdinStackLen())
	assert.Equal(t, 1, unit.StdoutStackLen())
	assert.Equal(t, 1, unit.StderrStackLen())

	// normally, I'd use assert.Equal() to prove that the pipe has data
	// if we did that here, the reads would empty the pipe, making the
	// rest of the test invalid

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetBuffers()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, unit.Stdin.String())
	assert.Empty(t, unit.Stdout.String())
	assert.Empty(t, unit.Stderr.String())

	// make sure that all of the io stacks have been emptied
	assert.Zero(t, unit.StdinStackLen())
	assert.Zero(t, unit.StdoutStackLen())
	assert.Zero(t, unit.StderrStackLen())
}

func TestPipeResetErrorCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetError()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetErrorCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetError()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeResetErrorSetsStatusCodeToStatusOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusNotOkay, nil
	}
	unit := pipe.NewPipe()
	unit.RunCommand(op1)

	assert.Equal(t, pipe.StatusNotOkay, unit.StatusCode())
	assert.Error(t, unit.Error())

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, pipe.StatusOkay, unit.StatusCode())
}

func TestPipeResetErrorSetsErrorToNil(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusNotOkay, nil
	}
	unit := pipe.NewPipe()
	unit.RunCommand(op1)

	assert.Equal(t, pipe.StatusNotOkay, unit.StatusCode())
	assert.Error(t, unit.Error())

	// ----------------------------------------------------------------
	// perform the change

	unit.ResetError()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, unit.Error())
}

func TestPipeRunCommandCopesWithNilPointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	expectedResult := pipe.StatusNotOkay
	op := func(p *pipe.Pipe) (int, error) {
		return expectedResult, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	// as long as it doesn't crash, the test has passed
}

func TestPipeRunCommandUpdatesStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	expectedResult := pipe.StatusNotOkay
	op := func(p *pipe.Pipe) (int, error) {
		return expectedResult, errors.New("status not okay")
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, unit.StatusCode())
}

func TestPipeRunCommandUpdatesErr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	expectedResult := errors.New("status not okay")
	op := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusNotOkay, expectedResult
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, unit.Error())
}

func TestPipeRunCommandSetsErrIfStatusCodeNotOkay(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	op := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusNotOkay, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op)

	// ----------------------------------------------------------------
	// test the results

	err := unit.Error()
	assert.NotNil(t, err)
	assert.Error(t, err)
	_, ok := err.(pipe.ErrNonZeroStatusCode)
	assert.True(t, ok)
}

func TestPipeSetNewStdinCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStdinCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStdin()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetStdinFromStringCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetStdinFromString("")

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetStdinFromStringCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetStdinFromString("")

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStdinCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PushStdin(ioextra.NewTextBuffer())

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStdinReplacesThePipeStdin(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdin"
	oldStdin := ioextra.NewTextBuffer()
	oldStdin.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdin = oldStdin

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stdin"
	newStdin := ioextra.NewTextBuffer()
	newStdin.WriteString(testData2)

	unit.PushStdin(newStdin)

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into newStdin
	actualResult := unit.Stdin.String()
	assert.Equal(t, testData2, actualResult)
}

func TestPipePushStdinAddsTheOldStdinToAnInternalStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdin"
	oldStdin := ioextra.NewTextBuffer()
	oldStdin.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdin = oldStdin

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stdin"
	newStdin := ioextra.NewTextBuffer()
	newStdin.WriteString(testData2)

	unit.PushStdin(newStdin)

	// ----------------------------------------------------------------
	// test the results

	// the internal stack should be larger now
	actualResult := unit.StdinStackLen()
	assert.Equal(t, 1, actualResult)
}

func TestPipePopStdinCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdin()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipePopStdinRestoresThePreviousPipeStdin(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdin"
	oldStdin := ioextra.NewTextBuffer()
	oldStdin.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdin = oldStdin

	// before we can test popping the stack, we need to push something
	// onto it
	testData2 := "this is the new stdin"
	newStdin := ioextra.NewTextBuffer()
	newStdin.WriteString(testData2)

	unit.PushStdin(newStdin)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdin()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdin
	actualResult := unit.Stdin.String()
	assert.Equal(t, testData1, actualResult)
	assert.Zero(t, unit.StdinStackLen())
}

func TestPipePopStdinDoesNothingWhenTheInternalStackIsEmpty(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can prove that popping an empty stack
	// makes no changes
	testData1 := "this is the old stdin"
	oldStdin := ioextra.NewTextBuffer()
	oldStdin.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdin = oldStdin

	// make sure the stack is empty
	assert.Zero(t, unit.StdinStackLen())

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdin()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdin
	actualResult := unit.Stdin.String()
	assert.Equal(t, testData1, actualResult)
}

func TestPipeStdinStackLenCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.StdinStackLen()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipeStdinStackLenReturnsLengthOfInternalStdinStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	testData1 := "this is the second stdin; it replaces the one created by NewPipe()"
	secondStdin := ioextra.NewTextBuffer()
	secondStdin.WriteString(testData1)

	unit.PushStdin(secondStdin)

	testData2 := "this is the third stdin; it does not count"
	thirdStdin := ioextra.NewTextBuffer()
	thirdStdin.WriteString(testData2)

	unit.PushStdin(thirdStdin)

	expectedResult := 2

	// ----------------------------------------------------------------
	// perform the change

	// this is how package users will get the stack length
	actualResult := unit.StdinStackLen()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeSetNewStdoutCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStdoutCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStdout()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStdoutCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PushStdout(ioextra.NewTextBuffer())

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStdoutReplacesThePipeStdout(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stdout"
	newStdout := ioextra.NewTextBuffer()
	newStdout.WriteString(testData2)

	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into newStdout
	actualResult := unit.Stdout.String()
	assert.Equal(t, testData2, actualResult)
}

func TestPipePushStdoutReplacesThePipeStderrIfNeeded(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStdout := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout
	unit.Stderr = oldStdout

	// ----------------------------------------------------------------
	// perform the change

	newStdout := ioextra.NewTextBuffer()

	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// test the results

	// the pipe's Stdout and Stderr should be the same
	assert.Equal(t, unit.Stdout, newStdout)
	assert.Equal(t, unit.Stderr, newStdout)
}

func TestPipePushStdoutAddsTheOldStdoutToAnInternalStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stdout"
	newStdout := ioextra.NewTextBuffer()
	newStdout.WriteString(testData2)

	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// test the results

	// the internal stack should be larger now
	actualResult := unit.StdoutStackLen()
	assert.Equal(t, 1, actualResult)
}

func TestPipePopStdoutCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdout()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipePopStdoutRestoresThePreviousPipeStdout(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// before we can test popping the stack, we need to push something
	// onto it
	testData2 := "this is the new stdout"
	newStdout := ioextra.NewTextBuffer()
	newStdout.WriteString(testData2)

	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdout()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdout
	actualResult := unit.Stdout.String()
	assert.Equal(t, testData1, actualResult)
	assert.Zero(t, unit.StdoutStackLen())
}

func TestPipePopStdoutRestoresThePreviousPipeStderrIfNecessary(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStdout := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stderr = oldStdout
	unit.Stdout = oldStdout

	newStdout := ioextra.NewTextBuffer()
	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdout()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, unit.Stdout, oldStdout)
	assert.Equal(t, unit.Stderr, oldStdout)
}

func TestPipePopStdoutDoesNothingWhenTheInternalStackIsEmpty(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can prove that popping an empty stack
	// makes no changes
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// make sure the stack is empty
	assert.Zero(t, unit.StdoutStackLen())

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdout()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdout
	actualResult := unit.Stdout.String()
	assert.Equal(t, testData1, actualResult)
}

func TestPipePopStdoutOnlyCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdoutOnly()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipePopStdoutOnlyRestoresThePreviousPipeStdout(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// before we can test popping the stack, we need to push something
	// onto it
	testData2 := "this is the new stdout"
	newStdout := ioextra.NewTextBuffer()
	newStdout.WriteString(testData2)

	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdoutOnly()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdout
	actualResult := unit.Stdout.String()
	assert.Equal(t, testData1, actualResult)
	assert.Zero(t, unit.StdoutStackLen())
}

func TestPipePopStdoutOnlyDoesNotRestoreThePreviousPipeStderr(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStdout := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stderr = oldStdout
	unit.Stdout = oldStdout

	newStdout := ioextra.NewTextBuffer()
	unit.PushStdout(newStdout)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdoutOnly()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, unit.Stdout, oldStdout)
	assert.Equal(t, unit.Stderr, newStdout)
}

func TestPipePopStdoutOnlyDoesNothingWhenTheInternalStackIsEmpty(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can prove that popping an empty stack
	// makes no changes
	testData1 := "this is the old stdout"
	oldStdout := ioextra.NewTextBuffer()
	oldStdout.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stdout = oldStdout

	// make sure the stack is empty
	assert.Zero(t, unit.StdoutStackLen())

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStdoutOnly()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStdout
	actualResult := unit.Stdout.String()
	assert.Equal(t, testData1, actualResult)
}

func TestPipeStdoutStackLenCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.StdoutStackLen()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipeStdoutStackLenReturnsLengthOfInternalStdinStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	testData1 := "this is the second stdout; it replaces the one created by NewPipe()"
	secondStdout := ioextra.NewTextBuffer()
	secondStdout.WriteString(testData1)

	unit.PushStdout(secondStdout)

	testData2 := "this is the third stdout; it does not count"
	thirdStdout := ioextra.NewTextBuffer()
	thirdStdout.WriteString(testData2)

	unit.PushStdout(thirdStdout)

	expectedResult := 2

	// ----------------------------------------------------------------
	// perform the change

	// this is how package users will get the stack length
	actualResult := unit.StdoutStackLen()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeSetNewStderrCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStderr()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipeSetNewStderrCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.SetNewStderr()

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStderrCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PushStderr(ioextra.NewTextBuffer())

	// ----------------------------------------------------------------
	// test the results
	//
	// as long as the code doesn't segfault, it works!
}

func TestPipePushStderrReplacesThePipeStderr(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stderr"
	newStderr := ioextra.NewTextBuffer()
	newStderr.WriteString(testData2)

	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into newStderr
	actualResult := unit.Stderr.String()
	assert.Equal(t, testData2, actualResult)
}

func TestPipePushStderrReplacesThePipeStdoutIfNeeded(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStderr := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stdout = oldStderr
	unit.Stderr = oldStderr

	// ----------------------------------------------------------------
	// perform the change

	newStderr := ioextra.NewTextBuffer()

	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// test the results

	// the pipe's Stdout and Stderr should be the same
	assert.Equal(t, unit.Stdout, newStderr)
	assert.Equal(t, unit.Stderr, newStderr)
}

func TestPipePushStderrAddsTheOldStderrToAnInternalStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// ----------------------------------------------------------------
	// perform the change

	testData2 := "this is the new stderr"
	newStderr := ioextra.NewTextBuffer()
	newStderr.WriteString(testData2)

	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// test the results

	// the internal stack should be larger now
	actualResult := unit.StderrStackLen()
	assert.Equal(t, 1, actualResult)
}

func TestPipePopStderrCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderr()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipePopStderrRestoresThePreviousPipeStderr(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// before we can test popping the stack, we need to push something
	// onto it
	testData2 := "this is the new stderr"
	newStderr := ioextra.NewTextBuffer()
	newStderr.WriteString(testData2)

	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderr()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStderr
	actualResult := unit.Stderr.String()
	assert.Equal(t, testData1, actualResult)
	assert.Zero(t, unit.StderrStackLen())
}

func TestPipePopStderrRestoresThePreviousPipeStdoutIfNecessary(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStderr := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr
	unit.Stdout = oldStderr

	newStderr := ioextra.NewTextBuffer()
	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderr()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, unit.Stdout, oldStderr)
	assert.Equal(t, unit.Stderr, oldStderr)
}

func TestPipePopStderrDoesNothingWhenTheInternalStackIsEmpty(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can prove that popping an empty stack
	// makes no changes
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// make sure the stack is empty
	assert.Zero(t, unit.StderrStackLen())

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderr()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStderr
	actualResult := unit.Stderr.String()
	assert.Equal(t, testData1, actualResult)
}

func TestPipePopStderrOnlyCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderrOnly()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipePopStderrOnlyRestoresThePreviousPipeStderr(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can try and tell the two buffers
	// apart later on
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// before we can test popping the stack, we need to push something
	// onto it
	testData2 := "this is the new stderr"
	newStderr := ioextra.NewTextBuffer()
	newStderr.WriteString(testData2)

	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderrOnly()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStderr
	actualResult := unit.Stderr.String()
	assert.Equal(t, testData1, actualResult)
	assert.Zero(t, unit.StderrStackLen())
}

func TestPipePopStderrOnlyNeverRestoresThePreviousPipeStdout(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	oldStderr := ioextra.NewTextBuffer()

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr
	unit.Stdout = oldStderr

	newStderr := ioextra.NewTextBuffer()
	unit.PushStderr(newStderr)

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderrOnly()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, unit.Stdout, newStderr)
	assert.Equal(t, unit.Stderr, oldStderr)
}

func TestPipePopStderrOnlyDoesNothingWhenTheInternalStackIsEmpty(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	// we need some data, so that we can prove that popping an empty stack
	// makes no changes
	testData1 := "this is the old stderr"
	oldStderr := ioextra.NewTextBuffer()
	oldStderr.WriteString(testData1)

	unit := pipe.NewPipe()
	unit.Stderr = oldStderr

	// make sure the stack is empty
	assert.Zero(t, unit.StderrStackLen())

	// ----------------------------------------------------------------
	// perform the change

	unit.PopStderrOnly()

	// ----------------------------------------------------------------
	// test the results

	// we should get back the data that was written into oldStderr
	actualResult := unit.Stderr.String()
	assert.Equal(t, testData1, actualResult)
}
func TestPipeStderrStackLenCopesWithNilPointer(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe

	// ----------------------------------------------------------------
	// perform the change

	unit.StderrStackLen()

	// ----------------------------------------------------------------
	// test the results

	// as long as the code does not segfault, it works!
}

func TestPipeStderrStackLenReturnsLengthOfInternalStderrStack(t *testing.T) {

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()

	testData1 := "this is the second stderr; it replaces the one created by NewPipe()"
	secondStderr := ioextra.NewTextBuffer()
	secondStderr.WriteString(testData1)

	unit.PushStderr(secondStderr)

	testData2 := "this is the third stderr; it does not count"
	thirdStderr := ioextra.NewTextBuffer()
	thirdStderr.WriteString(testData2)

	unit.PushStderr(thirdStderr)

	expectedResult := 2

	// ----------------------------------------------------------------
	// perform the change

	// this is how package users will get the stack length
	actualResult := unit.StderrStackLen()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusCodeCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe
	expectedResult := pipe.StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	actualResult := unit.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusCodeCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe
	expectedResult := pipe.StatusOkay

	// ----------------------------------------------------------------
	// perform the change

	actualResult := unit.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusCodeReturnsTheLastCommandsStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()
	expectedResult := 100

	op1 := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusOkay, nil
	}
	op2 := func(p *pipe.Pipe) (int, error) {
		return expectedResult, nil
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op1)
	unit.RunCommand(op2)
	actualResult := unit.StatusCode()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipeStatusErrorCopesWithNilPipePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit *pipe.Pipe
	expectedResult := pipe.StatusOkay
	var expectedErr error = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, actualErr := unit.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}

func TestPipeStatusErrorCopesWithEmptyPipe(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var unit pipe.Pipe
	expectedResult := pipe.StatusOkay
	var expectedErr error = nil

	// ----------------------------------------------------------------
	// perform the change

	actualResult, actualErr := unit.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}

func TestPipeStatusErrorReturnsTheLastCommandsStatusCodeAndError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	unit := pipe.NewPipe()
	expectedResult := 100
	expectedErr := errors.New("this is an error")

	op1 := func(p *pipe.Pipe) (int, error) {
		return pipe.StatusOkay, nil
	}
	op2 := func(p *pipe.Pipe) (int, error) {
		return expectedResult, expectedErr
	}

	// ----------------------------------------------------------------
	// perform the change

	unit.RunCommand(op1)
	unit.RunCommand(op2)
	actualResult, actualErr := unit.StatusError()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedErr, actualErr)
}
