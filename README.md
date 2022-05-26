# mgutil

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
```

Requires [mockgen](https://github.com/golang/mock) installed.

### TODO:

* Support flags to set a custom directory and prefix for generated
  mocks.
* Support flag for ignoring certain files/directories when scanning.


