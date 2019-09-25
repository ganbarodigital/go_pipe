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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvReturnsACopyOfTheProgramsEnvironment(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	os.Setenv(testKey, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	env := NewEnv()
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvGetenvReturnsFromTheEnvNotTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	os.Setenv(testKey, expectedResult)
	env := NewEnv()

	// now remove this from the program's environment
	os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvSetenvDoesNotChangeTheProgramEnv(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	env := NewEnv()

	// make sure this key does not exist in the program environment
	os.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	err := env.Setenv(testKey, expectedResult)
	envResult := os.Getenv(testKey)
	actualResult := env.Getenv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Nil(t, err)
	assert.Empty(t, envResult)
	assert.Equal(t, expectedResult, actualResult)
}

func TestEnvSetenvReturnsErrorForZeroLengthKey(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := ""
	testData := "this is a test"
	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestEnvSetenvReturnsErrorForKeyThatOnlyHasWhitespace(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "     "
	testData := "this is a test"
	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	ok := env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// test the results

	assert.Error(t, ok)
}

func TestEnvClearenvDeletesAllVariables(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is my value"

	env := NewEnv()
	env.Setenv(testKey, testData)

	// ----------------------------------------------------------------
	// perform the change

	env.Clearenv()

	// ----------------------------------------------------------------
	// test the results

	assert.Empty(t, env.Environ())
	assert.Empty(t, env.Getenv(testKey))
}

func TestEnvLookupEnvReturnsTrueIfTheVariableExists(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := "this is my value"

	env := NewEnv()
	env.Setenv(testKey, expectedResult)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestEnvLookupEnvReturnsFalseIfTheVariableDoesNotExist(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	expectedResult := ""

	env := NewEnv()
	env.Unsetenv(testKey)

	// ----------------------------------------------------------------
	// perform the change

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
}

func TestEnvUnsetenvDeletesAVariable(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is a test"
	expectedResult := ""

	env := NewEnv()
	env.Setenv(testKey, testData)

	origLen := env.Length()

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	actualResult, ok := env.LookupEnv(testKey)
	actualEnviron := env.Environ()

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.False(t, ok)
	assert.Equal(t, origLen-1, len(actualEnviron))

	// make sure it isn't in the environ too
	prefix := testKey + "="
	for _, pair := range actualEnviron {
		assert.False(t, strings.HasPrefix(pair, prefix))
	}
}

func TestEnvUnsetenvDoesNotChangeProgramEnviron(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData := "this is a test"
	expectedResult := testData

	os.Setenv(testKey, testData)

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey)

	actualResult, ok := os.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	// still in the environment
	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)

	// but gone from our Env
	assert.Equal(t, "", env.Getenv(testKey))
}

func TestEnvEntriesFromProgramEnvironmentCanBeUpdated(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey := "TestNewEnv"
	testData1 := "this is a test"
	testData2 := "this is another test"
	expectedResult := testData2

	os.Setenv(testKey, testData1)

	env := NewEnv()

	// ----------------------------------------------------------------
	// perform the change

	env.Setenv(testKey, testData2)

	actualResult, ok := env.LookupEnv(testKey)

	// ----------------------------------------------------------------
	// test the results

	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, ok)
}

func TestEnvUpdatedEntriesCanBeUnset(t *testing.T) {
	// ----------------------------------------------------------------
	// setup your test

	testKey1 := "TestNewEnv1"
	testKey2 := "TestNewEnv2"
	testData1 := "this is a test"
	testData2 := "this is another test"

	env := NewEnv()

	env.Setenv(testKey1, testData1)
	env.Setenv(testKey2, testData1)

	env.Setenv(testKey1, testData2)
	testValue, ok := env.LookupEnv(testKey1)
	assert.Equal(t, testData2, testValue)
	assert.True(t, ok)

	// ----------------------------------------------------------------
	// perform the change

	env.Unsetenv(testKey1)

	actualResult, ok := env.LookupEnv(testKey1)

	// ----------------------------------------------------------------
	// test the results

	assert.False(t, ok)
	assert.Empty(t, actualResult)

	// prove the 2nd entry hasn't been lost
	testValue, ok = env.LookupEnv(testKey2)
	assert.True(t, ok)
	assert.Equal(t, testData1, testValue)
}
