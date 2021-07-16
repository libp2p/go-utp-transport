# go-utp-transport

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](https://webchat.freenode.net/?channels=%23libp2p)
[![GoDoc](https://godoc.org/github.com/libp2p/go-utp-transport?status.svg)](https://godoc.org/github.com/libp2p/go-utp-transport)
[![Coverage Status](https://coveralls.io/repos/github/libp2p/go-utp-transport?branch=master)](https://coveralls.io/github/libp2p/go-utp-transport?branch=master)
[![Travis CI](https://travis-ci.org/libp2p/go-utp-transport?branch=master)](https://travis-ci.org/libp2p/go-utp-transport)

> A libp2p transport implementation for utp.

## ⚠️ Unmaintained ⚠️

This library is not currently maintained and isn't kept up-to-date with the latest go-libp2p releases. If you'd like to revive it, please fork it, get it working again, and start a discussion on https://discuss.libp2p.io/.

However, if you're just looking for a UDP-based libp2p transport, try [QUIC](https://github.com/libp2p/go-libp2p-quic-transport/). It's better in just about every way (1RTT handshake, built-in stream multiplexing with no head of line blocking, etc.).

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Contribute](#contribute)
- [License](#license)

## Install

```sh
go get github.com/libp2p/go-utp-transport
```

## API

Check out the [GoDocs](https://godoc.org/github.com/libp2p/go-utp-transport).

## Contribute

PRs are welcome!

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © Jeromy Johnson
