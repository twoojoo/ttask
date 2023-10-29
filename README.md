## TTask

A stream processing library for Go, heavily inspired by [alyxstream](https://github.com/smartpricing/alyxstream).

> **!!!** This module is **not idiomatic Go**. Due to the lack of generics on struct methods (*go 1.21.3*), I was forced to use a weird pattern to replicate a sort of fluent syntax while mantaining full type safety, hence the name of the module. If a future version of go support this feature, a new version of this module may be written.

