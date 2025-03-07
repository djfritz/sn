# Stack Noise

sn (stack noise) -- Make your Go stack traces more readable

<img alt="Demo" src="https://github.com/djfritz/sn/blob/main/demo.gif?raw=true" width="600" />

## Installing

If you're using `sn`, you probably have Go installed:

```
go install github.com/djfritz/sn@latest
```

## Overview

sn is a simple command that deduplicates goroutines in a [Go](https://golang.org) stacktrace, removes what are often unneeded pointers and register contents, and displays the stacktrace in a TUI. 

