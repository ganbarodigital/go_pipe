# CHANGELOG

## develop

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