# usrconfig

[![PkgGoDev](https://pkg.go.dev/badge/github.com/qeubar/usrconfig?tab=doc)](https://pkg.go.dev/github.com/qeubar/usrconfig?tab=doc)
[![Build Status](https://travis-ci.com/qeubar/usrconfig.svg?branch=master)](https://travis-ci.com/qeubar/usrconfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/qeubar/usrconfig)](https://goreportcard.com/report/github.com/qeubar/usrconfig)
[![codecov](https://codecov.io/gh/qeubar/usrconfig/branch/master/graph/badge.svg)](https://codecov.io/gh/qeubar/usrconfig)

usrconfig is a very simple platform independent user config management.

## Usage

```go

import "github.com/quebar/usrconfig"

type MyConfig struct {
    Name  string `yaml:"user_name"`
    Email string `yaml:"user_email"`
}

conf := MyConfig{
    Name: "QeuBar",
    Email: "que@bar.com",
}

usrconfig.Update(conf, "my-app")
```

### Supported platforms
usrconfig is built to work on any platform that supports Go.
