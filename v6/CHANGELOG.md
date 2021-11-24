# CHANGELOG

## develop

The main focus for v6.0.0 has been to extract reusable I/O concepts out into a separate `ioextra` package.

### Refactoring

* `ReadBuffer` has been replaced by `ioextra.TextReader`
* `WriteBuffer` has been replaced by `ioextra.TextWriter`
* `Source` has been replaced by `ioextra.TextBuffer`
* `Dest` has been replaced by `ioextra.TextBuffer`
* `NewScanReader` has been replaced by `ioextra.NewTextScanner`

### New

* Added `AttachOsStdin`
* Added `AttachOsStdout`
* Added `AttachOsStderr`

### Fixes

* `Pipe.Stdin` is now an `ioextra.TextReader`
* `Pipe.Stdout` is now an `ioextra.TextWriter`
* `Pipe.Stderr` is now an `ioextra.TextWriter`

## v5.2.0

Released Sunday, 3rd November 2019.

### New

* Added `Pipe.WriteBuffer` interface

### Fixes

* `NewScanReader()` now panics if it is given a nil pointer
* Reads on a `Dest` struct are now destructive

## v5.1.0

Released Sunday, 3rd November 2019.

### New

* Added `Pipe.Flags` property

## v5.0.0

Released Wednesday, 30th October 2019.

### Backwards-Compatibility Breaks

* Dropped `Pipe.Vars`
* `Pipe.Env` now defaults to the program's environment
  - you can replace it with an `OverlayEnv` from Envish if you'd like
* Dropped `Pipe.Getvar()`
  - call `Pipe.Env.Getenv()` instead
* Dropped `Pipe.Expand()`
  - call `Pipe.Env.Expand()` instead

### Dependencies

- bumped to go_envish v3.0.1

## v4.3.0

Released Tuesday, 8th October 2019.

### New

* `Pipe.Vars.Expand()` now supports UNIX shell special variables (`$#`, `$*` et al)

## v4.2.0

Released Monday, 7th October 2019.

### New

* Added `Pipe.Vars` to store local variables
* Added `Pipe.Getvar()` to retrieve a variable from:
  - the pipe's local variable store
  - the pipe's local environment store
  - the program's environment
  in that order.

## v4.1.0

Released Sunday, 6th October 2019.

### New

* Added `Pipe.Okay()`
* Added `Pipe.StatusError()` to avoid calling `Pipe.StatusCode()` and `Pipe.Error()` together.

## v4.0.0

Released Sunday, 6th October 2019.

### Breaking Changes

We're making some changes to further improve support for building UNIX-like shell behaviour in Golang packages and apps.

- Status codes and errors have moved into the `Pipe`
- The local Env has moved into the `Pipe`
- Sequence, Pipeline and List have been moved into the `go_scriptish` package
- `Pipe.DrainStdin()` is now `Pipe.DrainStdinToStdout()`
- Removed `Pipe.Next()`
  - this behaviour belongs in packages that use `go_pipe`
- `Pipe.StatusCode` is now a method, not an exported struct member
- `Pipe.Err` is no longer an exported struct member; use `Pipe.Error()` instead
- `Pipe.Reset()` is now `Pipe.ResetBuffers()`

### New

* `NewPipe()` now accepts option functions
* Added `Pipe.ResetError()`
* Added `Pipe.SetNewStdin()`
* Added `Pipe.SetStdinFromString()`
* Added `Pipe.SetNewStdout()`
* Added `Pipe.SetNewStderr()`

## v3.1.0

Released Sunday, 29th September 2019.

### New

The `Env` structure added in v3.0.0 has now been spun out into a separate [Envish](https://github.com/ganbarodigital/go_envish) package.

## v3.0.0

Released Wednesday, 25th September 2019.

### Breaking Changes

We're making some changes to improve compatibility with UNIX shell terminology and behaviours.

- `PipelineOperation` is now `Command`
- `Sequence` represents a set of Commands to be executed
- `NewPipeline()` now returns a `Sequence` (used to return a `*Pipeline`)
- `NewList()` creates a `Sequence` that executes like a UNIX shell list
- `ErrPipelineNonZeroStatusCode` is now `ErrNonZeroStatusCode`

### New

- Added `Env` for a local environment
- Added `ErrEmptyKey` error
- Added local environment support to Sequences
- Added `Sequence.Expand()` to mimic `os.Expand()`

## v2.0.1

Released Wednesday, 25th September 2019.

### Fixes

* Golang v2+ module compatibility fixes (sigh)

## v2.0.0

Released Wednesday, 25th September 2019.

### Breaking Changes

* Pipeline capture methods now *always* return the pipeline's Stdout's contents, regardless of the pipeline's error status.
  - this change improves compatibility with UNIX behaviour

## v1.7.0

Released Wednesday, 25th September 2019.

### New

* Added `ErrPipelineNonZeroStatusCode`
* `Pipeline.Exec_()` now sets `Pipeline.Err` if a step fails with a non-zero status code, but no error of its own

## v1.6.0

Released Wednesday, 25th September 2019.

### New

* Added `Pipe.Reset()`

## v1.5.1

Released Tuesday, 17th September 2019.

### Fixes

* `Pipeline.Okay()` now returns `(bool, error)`

## v1.5.0

Released Tuesday, 17th September 2019.

### New

We've been adding some more convenience methods, to help [scriptish](https://github.com/ganbarodigital/go_scriptish).

* Added `Pipeline.Okay()`

## v1.4.0

Released Monday, 16th September 2019.

### New

We've been adding some more convenience methods, to help [scriptish](https://github.com/ganbarodigital/go_scriptish).

* Added `Pipeline.Error()`

## v1.3.0

Released Monday, 16th September 2019.

### New

We've been adding some more convenience methods, to help [scriptish](https://github.com/ganbarodigital/go_scriptish).

* Added `Pipeline.TrimmedString()`

## v1.2.0

Released Monday, 16th September 2019.

### New

We've been adding some more convenience methods, to help [scriptish](https://github.com/ganbarodigital/go_scriptish).

* Added `Source.ParseInt()`
* Added `Source.TrimmedString()`
* Added `Dest.ParseInt()`
* Added `Dest.TrimmedString()`
* Added `Pipeline.ParseInt()`
* Added `ReadBuffer.ParseInt()`
* Added `ReadBuffer.TrimmedString()`

## v1.1.0

Released Sunday, 15th September 2019.

### New

* Added `Pipeline.Exec_()` to help with embedding Pipeline in other structs.

## v1.0.2

Released Sunday, 15th September 2019.

### Fixes

* `Pipeline.Exec()` returns pointer to self, allows method chaining

## v1.0.1

Released Sunday, 15th September 2019.

### Fixes

* Fixed typo in Go module name

## v1.0.0

Released Sunday, 15th September 2019.

### New

* Added `Dest`, a `strings.Builder` with useful helper methods
* Added `Source`, an `io.ReadCloser` with useful helper methods
* Added `NewScanReader`, wraps a channel around a scanner and given split function
* Added `ReadBuffer`, an interface to describe our Source/Dest common helper methods
* Added `Pipe`, which represents UNIX-like stdin, stdout and stderr
* Added `PipelineOperation`, a function signature for a task that supports our `Pipe`
* Added `Pipeline`, a UNIX-like pipeline ready to execute
