# mgg

Scans all files in your project. If the path doesn't contain '_test',
the prefix flag, ends in '.go' and file contains 'interface {' it
generates mocks for that interface using
[mockgen](https://github.com/golang/mock).

Usage:

```
# USAGE:
#     mgg [OPTIONS]
# 
# OPTIONS:
#     -h, --help	  Prints this message
#     -d, --dir       Directory to generate mocks in [default: 'mocks']
#     -p, --prefix    Prefix to use for mock files [default: 'mock_']
#     -i, --ignore    Paths to ignore when scanning for interfaces [default: ['']]
```

Example:

```
# .
# |--- .git
# |--- README.md
# |--- go.mod
# |--- main.go # has interface
# |--- mocks
#     |--- mock_pkg
#         |--- mock_logger.go
#         |--- mock_pubsub.go
# |--- pkg
#     |--- logger.go # has interface
#     |--- pubsub.go # has interface

$ mgg --dir=mocks --prefix=mock_ --ignore=main.go,pkg/logger.go
Generated 'mocks/mock_pkg/mock_pubsub.go'
```

Requires [mockgen](https://github.com/golang/mock) installed.

### TODO:

* Respect `.gitignore`.
* Support passing flags to `mockgen`
* Create unit tests
