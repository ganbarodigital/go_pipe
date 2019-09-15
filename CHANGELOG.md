# CHANGELOG

## develop

### New

* Added `Dest`, a `strings.Builder` with useful helper methods
* Added `Source`, an `io.ReadCloser` with useful helper methods
* Added `NewScanReader`, wraps a channel around a scanner and given split function
* Added `ReadBuffer`, an interface to describe our Source/Dest common helper methods
* Added `Pipe`, which represents UNIX-like stdin, stdout and stderr
* Added `PipelineOperation`, a function signature for a task that supports our `Pipe`
* Added `Pipeline`, a UNIX-like pipeline ready to execute