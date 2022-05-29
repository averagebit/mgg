# mgg

Scans all files in your project. If the path doesn't contain '_test',
ends in '.go' and file contains 'interface {' it generates a mock for
that interface using [mockgen](https://github.com/golang/mock).

Mocks are generated in the current working directory under the 'mocks'
folder and use the original file path with each section prefixed with
'mock_'.

Example:

```
# Source:
#   app/app.go
#
# Generates:
#   mocks/mock_app/mock_app.go
# 
# USAGE:
#     mgg [FLAGS]
# 
# FLAGS:
#     -d, --dir       directory to generate mock files in
#     -r, --remove    remove old mock files before generating
#     -p, --prefix    prefix to use for mock files shorthand

$ mgg --remove --dir=mocks --prefix=mock_
```

Requires [mockgen](https://github.com/golang/mock) installed.

### TODO:

* Support flag for ignoring certain files/directories when scanning.
  Respect .gitignore
* Support flag for generating only updated files. git diff
* Support flag for dry runs
* Support passing flags to `mockgen`
