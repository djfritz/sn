# sn

sn (stack noise) -- Make your Go stack traces more readable

[demo.gif]

## Overview

sn is a simple command that deduplicates goroutines in a [Go](https://golang.org) stacktrace, removes what are often unneeded pointers and register contents, and displays the stacktrace in a TUI. 

