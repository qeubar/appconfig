# appconfig

[![GoDoc](https://godoc.org/github.com/qeubar/appconfig?status.svg)](https://godoc.org/github.com/qeubar/appconfig)
[![Build Status](https://travis-ci.com/qeubar/appconfig.svg?branch=master)](https://travis-ci.com/qeubar/appconfig)

appconfig is a very simple platform independent config file management.

## Usage

```go

import "github.com/quebar/appconfig"

type MyConfig struct {
    Name  string `yaml:"user_name"`
    Email string `yaml:"user_email"`
}

conf := MyConfig{
    Name: "QeuBar",
    Email: "que@bar.com",
}

appconfig.Update(conf, "my-app")
```

### Supported platforms
appconfig is built to work on any platform that supports Go.
