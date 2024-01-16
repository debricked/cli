# Go resolution logic

The way resolution of go lock files works is as follows:

1. Run `go mod graph` in order to create dependency graph
2. Run `go list -mod=readonly -e -m all` to get the contents of the requirements.txt file

The results of the commands above are then combined to form the finished lock file.
