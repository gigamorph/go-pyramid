# go-pyramid
Convert images to pyramidal TIFF.

It depends entirely on shell invoked programs for its operations.
The `explore-cgo` branch tries to incorporate C libraries but is practically abandoned at the moment.

## Running as Standalone

```bash
go run main/pyramid/pyramid.go [<options>] <infile> <outfile>
```

### Options

* -m
* -c
* -q
