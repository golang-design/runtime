# runtime [![PkgGoDev](https://pkg.go.dev/badge/golang.design/x/runtime)](https://pkg.go.dev/golang.design/x/runtime) ![](https://changkun.de/urlstat?mode=github&repo=golang-design/runtime)![runtime](https://github.com/golang-design/runtime/workflows/runtime/badge.svg?branch=main)

an extend to Go `runtime` package

```diff
-import "runtime"
+import "golang.design/x/runtime"
```

## Subpackages

- [`thread`](./thread) — threading facilities: scheduling calls on a
  specific OS thread, thread-local storage, and a bounded thread pool.
  Previously published as `golang.design/x/thread`.
- [`mainthread`](./mainthread) — schedule functions to run on the main OS
  thread (e.g. for GUI/OpenGL). Previously published as
  `golang.design/x/mainthread`.

## License

MIT | &copy; 2021 The golang.design Initiative Authors, written by [Changkun Ou](https://changkun.de).