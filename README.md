# go-dslp
![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)
![Programming Language](https://img.shields.io/badge/language-Go-orange.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/iceflowRE/go-dslp/blob/master/LICENSE.md)

[![Github Actions](https://github.com/IceflowRE/go-dslp/workflows/Build/badge.svg)](https://github.com/IceflowRE/go-dslp/actions)
[![Go report card](https://goreportcard.com/badge/github.com/IceflowRE/go-dslp)](https://goreportcard.com/report/github.com/IceflowRE/go-dslp)

---

## Distributed Systems Learning Protocol

The Distributed Systems Learning Protocol (DSLP) realizes the transmission of messages over an already established transport connection, in this case TCP.

---

## Requirements

- Go (>= 1.10)

## Build

- `go build -x -o go-dslp`

## Run options

- --client \<address>

- --server \<port>
    
- --version \<version>

    1.2 | 2.0
    
examples:

    go-dslp --server 28813 --version 2.0
    go-dslp --client localhost:28813 --version 1.2

---

## Web
https://github.com/IceflowRE/go-dslp

## Credits
- Developer
    - Iceflower S
        - iceflower@gmx.de

### Third Party
Nothing.

## License
Copyright 2019-present Iceflower S (iceflower@gmx.de)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
