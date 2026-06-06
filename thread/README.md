# thread [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/runtime/thread)](https://pkg.go.dev/golang.design/x/runtime/thread)

Package thread provides threading facilities, such as scheduling
calls on a specific thread, local storage, etc.

> This package was previously published as `golang.design/x/thread`. It now
> lives in the `golang.design/x/runtime` module; update imports accordingly.

```go
import "golang.design/x/runtime/thread"
```

## Quick Start

```go
th := thread.New()
defer th.Terminate()

// Run on the thread and block until it returns.
th.Call(func() {
    // ... runs on the created thread ...
})

// Schedule on the thread without waiting.
th.Go(func() {
    // ... runs on the created thread ...
})

// Run on the thread and return a typed value.
n := thread.Eval(th, func() int {
    return 42
})
```

## License

MIT &copy; 2020 - 2026 The golang.design Initiative