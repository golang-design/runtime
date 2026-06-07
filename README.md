# runtime [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/runtime)](https://pkg.go.dev/golang.design/x/runtime) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/runtime)![runtime](https://github.com/golang-design/runtime/workflows/runtime/badge.svg?branch=main)

runtime utilities that complement the standard library `runtime` package

```go
import "golang.design/x/runtime"
```

The root package provides:

- `Goid` — the ID of the current goroutine.
- `SetMaxThreads` / `WaitThreads` / `NumThreads` — cap the number of OS
  threads the process uses, and wait until the count converges.

By default `Goid` parses `runtime.Stack` (portable, ~µs). Build with
`-tags goidfast` on amd64/arm64 for a ~ns implementation that reads the id
directly from the runtime goroutine structure (the offset is calibrated at
startup, with a safe fallback to the portable path):

```
go build -tags goidfast ./...
```

Use it alongside the standard `runtime` package; it does not re-export it.

## Subpackages

- [`thread`](./thread) — threading facilities: scheduling calls on a
  specific OS thread, thread-local storage, and a bounded thread pool.
  Previously published as `golang.design/x/thread`.
- [`mainthread`](./mainthread) — schedule functions to run on the main OS
  thread (e.g. for GUI/OpenGL). Previously published as
  `golang.design/x/mainthread`.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).