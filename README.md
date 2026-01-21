# go-serix
![](https://img.shields.io/badge/status-Work%20In%20Progress-8A2BE2)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-serix)
[![test](https://github.com/cybergarage/go-serix/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-serix/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-serix.svg)](https://pkg.go.dev/github.com/cybergarage/go-serix)
[![codecov](https://codecov.io/gh/cybergarage/go-serix/graph/badge.svg?token=GOLCBMUVB1)](https://codecov.io/gh/cybergarage/go-serix)

## Overview

`go-serix` is a Go library that provides an idiomatic, schema-driven API for encoding and decoding structured data. Serix is a coined term used in this project to describe a compact, schema-first serialization approach for representing and exchanging structured data. The library is designed to be lightweight and easy to integrate into existing Go services, with a focus on clarity, maintainability, and performance.

## Features

- Stable high-level interfaces: keep your application-facing API the same
- Pluggable serialization formats: switch the underlying wire format without changing the upper-layer code
- Composable compression: combine your chosen serialization format with different compression strategies
- Integration-friendly design for production Go services

The name "Serix" is derived from "Serialize" and "X" (standing for extensibility), reflecting the project's goal of providing a flexible and extensible serialization framework.
