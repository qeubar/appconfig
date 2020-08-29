# usrconfig

[![GoDoc](https://godoc.org/github.com/qeubar/usrconfig?status.svg)](https://godoc.org/github.com/qeubar/usrconfig)
[![Build Status](https://travis-ci.com/qeubar/usrconfig.svg?branch=master)](https://travis-ci.com/qeubar/usrconfig)

usrconfig is a very simple platform independent user config anagement.

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
