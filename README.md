# protoparts

[![Go Reference](https://pkg.go.dev/badge/github.com/obeattie/protoparts.svg)](https://pkg.go.dev/github.com/obeattie/protoparts)

Facilities for splitting apart binary Protocol Buffer messages into individual fields, and combining the parts again into whole (or partial) Protocol Buffer messages.

The motivation is to enable Protocol Buffers to be used as a storage format, where fields need to be addressable on-disk individually to enable [projection](https://en.wikipedia.org/wiki/Projection_(relational_algebra)), [selection](https://en.wikipedia.org/wiki/Selection_(relational_algebra)), and so on.
