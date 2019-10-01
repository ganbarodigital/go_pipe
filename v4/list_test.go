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

	"github.com/stretchr/testify/assert"
)

func TestNewListCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()
	actualResult := list.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewListCreatesSequenceWithNilErrSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, list.Error())
}

func TestNewListCreatesSequenceWithZeroStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	list := NewList()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, list.StatusCode())
}

func TestListControllerCopesWithNilSequencePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence *Sequence
	controller := ListController(sequence)

	// ----------------------------------------------------------------
	// perform the change

	controller()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestListControllerCopesWithEmptySequence(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var sequence Sequence
	controller := ListController(&sequence)
	sequence.Controller = controller

	// ----------------------------------------------------------------
	// perform the change

	sequence.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestListExecRunsAllStepsInOrder(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("hello world")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}
	op2 := func(p *Pipe) (int, error) {
		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	list := NewList(op1, op2)
	list.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := list.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestListExecDoesNotStopWhenAStepReportsAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStdout := "hello world\nhave a nice day\n"
	expectedStderr := "alfred the great\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("hello world")
		p.Stdout.WriteRune('\n')
		p.Stderr.WriteString(expectedStderr)

		// all done
		return StatusNotOkay, errors.New("stop at step 1")
	}
	op2 := func(p *Pipe) (int, error) {
		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	list := NewList(op1, op2)
	list.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualStdout, err := list.String()
	actualStderr := list.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	// list.String() should not have returned an error
	assert.Nil(t, err)
	assert.Equal(t, StatusOkay, list.StatusCode())

	// list.String() should have returned the output from both
	// of our test operations
	assert.Equal(t, expectedStdout, actualStdout)

	// our list's Stderr should not have been reset along the way
	assert.Equal(t, expectedStderr, actualStderr)
}

func TestListExecSetsErrWhenOpReturnsNonZeroStatusCodeAndNilErr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		// fail, but without an error to say why
		return StatusNotOkay, nil
	}

	list := NewList(op1)

	// ----------------------------------------------------------------
	// perform the change

	list.Exec()

	// ----------------------------------------------------------------
	// test the results

	// list.Err should have been set by Exec()
	assert.NotNil(t, list.Error())
	_, ok := list.Error().(ErrNonZeroStatusCode)
	assert.True(t, ok)
}
