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

func TestNewPipelineCreatesPipeWithEmptyStdin(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stdin.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipeWithEmptyStdout(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipeWithEmptyStderr(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()
	actualResult := pipeline.Pipe.Stderr.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestNewPipelineCreatesPipelineWithNilErrSet(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, pipeline.Err)
}

func TestNewPipelineCreatesPipelineWithZeroStatusCode(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	// ----------------------------------------------------------------
	// perform the change

	pipeline := NewPipeline()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, 0, pipeline.StatusCode)
}

// helper for testing our pipeline behaviour
type testOpResult struct {
	StatusCode int
	Err        error
}

func TestNewPipelineCreatesPipelineWithGivenPipelineOperations(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) { return 0, nil }
	op2 := func(p *Pipe) (int, error) { return 1, nil }

	// we can't compare functions directly in Go, but we can execute them
	// and compare their output
	expectedResult := []testOpResult{{0, nil}, {1, nil}}

	// ----------------------------------------------------------------
	// perform the change

	var actualResult []testOpResult

	pipeline := NewPipeline(op1, op2)
	for _, step := range pipeline.Steps {
		statusCode, err := step(pipeline.Pipe)
		actualResult = append(actualResult, testOpResult{statusCode, err})
	}

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineExecCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestPipelineExecCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline

	// ----------------------------------------------------------------
	// perform the change

	pipeline.Exec()

	// ----------------------------------------------------------------
	// test the results

	// as long as it didn't crash, we're good
}

func TestPipelineExecRunsAllStepsInOrder(t *testing.T) {
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
		// copy what op1 did first
		p.DrainStdin()

		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1, op2)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineExecStopsWhenAStepReportsAnError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedStdout := "hello world\n"
	expectedStderr := "alfred the great\n"
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString(expectedStdout)
		p.Stderr.WriteString(expectedStderr)

		// all done
		return 0, errors.New("stop at step 1")
	}
	op2 := func(p *Pipe) (int, error) {
		// copy what op1 did first
		p.DrainStdin()

		// add our own content
		p.Stdout.WriteString("have a nice day")
		p.Stdout.WriteRune('\n')

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1, op2)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	finalOutput, err := pipeline.String()
	actualStdout := pipeline.Pipe.Stdout.String()

	// ----------------------------------------------------------------
	// test the results

	// pipeline.String() should have returned an error
	assert.NotNil(t, err)

	// pipeline.String() should have returned the contents of our
	// Pipe.Stderr buffer
	assert.Equal(t, expectedStderr, finalOutput)

	// our pipeline's Stdout should still contain what the first step
	// did ... and only the first step
	assert.Equal(t, expectedStdout, actualStdout)
}

func TestPipelineBytesCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := pipeline.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineBytesCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := pipeline.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineBytesReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := pipeline.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineBytesReturnsContentsOfStderrWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")

		// this is the content we want
		p.Stderr.WriteString(expectedResult)

		// all done
		return 0, errors.New("an error occurred")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualBytes, err := pipeline.Bytes()
	actualResult := string(actualBytes)

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineErrorCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline

	// ----------------------------------------------------------------
	// perform the change

	err := pipeline.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestPipelineErrorCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline

	// ----------------------------------------------------------------
	// perform the change

	err := pipeline.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
}

func TestPipelineErrorReturnsErrorWhenOneHappens(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		// all done
		return 0, errors.New("this is an error")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	err := pipeline.Error()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Error(t, err)
}

func TestPipelineOkayCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline

	// ----------------------------------------------------------------
	// perform the change

	success, err := pipeline.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, success)
	assert.Nil(t, err)
}

func TestPipelineOkayCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline

	// ----------------------------------------------------------------
	// perform the change

	success, err := pipeline.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.True(t, success)
	assert.Nil(t, err)
}

func TestPipelineOkayReturnsFalseWhenPipelineErrorHappens(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	op1 := func(p *Pipe) (int, error) {
		// all done
		return NOT_OK, errors.New("this is an error")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	success, err := pipeline.Okay()

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, success)
	assert.Error(t, err)
}

func TestPipelineParseIntCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineParseIntCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline
	expectedResult := 0

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineParseIntConvertsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 100
	op1 := func(p *Pipe) (int, error) {
		p.Stdout.WriteString("100\n")

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineParseIntReturnsZeroWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := 0
	op1 := func(p *Pipe) (int, error) {
		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")
		p.Stderr.WriteString("not a number")

		// all done
		return 0, errors.New("an error occurred")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.ParseInt()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringReturnsContentsOfStderrWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := "hello world\nhave a nice day\n"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stderr.WriteString(expectedResult)

		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")

		// all done
		return 0, errors.New("an eccor occurred")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.String()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringsCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringsCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline
	expectedResult := []string{}

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringsReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{"hello world", "have a nice day"}
	op1 := func(p *Pipe) (int, error) {
		for _, line := range expectedResult {
			p.Stdout.WriteString(line)
			p.Stdout.WriteRune('\n')
		}

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineStringsReturnsContentsOfStderrWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	expectedResult := []string{"hello world", "have a nice day"}
	op1 := func(p *Pipe) (int, error) {
		for _, line := range expectedResult {
			p.Stderr.WriteString(line)
			p.Stderr.WriteRune('\n')
		}

		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")

		// all done
		return 0, errors.New("an error occurred")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.Strings()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineTrimmedStringCopesWithNilPipelinePointer(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline *Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineTrimmedStringCopesWithEmptyPipeline(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	var pipeline Pipeline
	expectedResult := ""

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineTrimmedStringReturnsContentsOfStdoutWhenNoError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "   hello world\nhave a nice day\n\n\n"
	expectedResult := "hello world\nhave a nice day"

	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stdout.WriteString(testData)

		// we don't want to see this in our final output
		p.Stderr.WriteString("we do not want this")

		// all done
		return 0, nil
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}

func TestPipelineTrimmedStringReturnsContentsOfStderrWhenError(t *testing.T) {
	t.Parallel()

	// ----------------------------------------------------------------
	// setup your test

	testData := "   hello world\nhave a nice day\n\n\n"
	expectedResult := "hello world\nhave a nice day"
	op1 := func(p *Pipe) (int, error) {
		// this is the content we want
		p.Stderr.WriteString(testData)

		// we don't want to see this in our final output
		p.Stdout.WriteString("we do not want this")

		// all done
		return 0, errors.New("an error occurred")
	}

	pipeline := NewPipeline(op1)
	pipeline.Exec()

	// ----------------------------------------------------------------
	// perform the change

	actualResult, err := pipeline.TrimmedString()

	// ----------------------------------------------------------------
	// test the results

	assert.NotNil(t, err)
	assert.Equal(t, expectedResult, actualResult)
}
