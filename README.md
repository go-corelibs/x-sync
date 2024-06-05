[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-corelibs/x-sync)
[![codecov](https://codecov.io/gh/go-corelibs/x-sync/graph/badge.svg?token=2hZkB2epa4)](https://codecov.io/gh/go-corelibs/x-sync)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-corelibs/x-sync)](https://goreportcard.com/report/github.com/go-corelibs/x-sync)

# x-sync - experimental x-sync utilities

x-sync is likely not what you're looking for. This project exists to mainly
have a reference for things one may think are great optimizations of standard
go things but did not turn out to be the case.

# General Notes

1. the built-in `append` function is really fast as it is but if there's a
   need for controlling the slice scaling process, this package has an
   `AppendScale` function which may be useful
2. working with sync.Pool is a very specific task, usually unique to each
   project's needs, however if there's a need for a generic and convenient
   sync.Pool that doesn't need boilerplate, take at look at this package's
   Pool type and the NewStringBuilderPool convenience function
3. some benchmarks have been added and they basically demonstrate that the
   standard Go `append` function is really tough to beat, though admittedly
   the case scenario being benchmarked may not be the best implementation
   for comparative benchmarking purposes
4. don't optimize too early in the development cycle, just get things done
   until it's obviously working and then consider this package for testing
   really simply optimizations before delving into the more appropriate
   project-specific minimum viable product performance tuning

# Installation

``` shell
> go get github.com/go-corelibs/x-sync@latest
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

```
Copyright 2024 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
