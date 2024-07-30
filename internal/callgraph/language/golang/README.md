# Go CallGraph Generation

Building callgraphs for Go is generally easy, but there are some caveats.
If your project depends on building non-Go components then you need to manually build the project before building the callgraph.
Once the project is built you can generate the callgraph with the following command:
```shell
debricked callgraph . 
```

And then upload it and scan it using:


```shell
debricked scan .
```

# Additional Information

The callgraph generation depends only on internal functionality in the Go standard library, for more information about this and the implementation see: 
https://cs.opensource.google/go/x/tools/+/refs/tags/v0.19.0:cmd/callgraph/main.go

As always, callgraph cannot be expected to include all possible calls in your program and not all included calls are guaranteed to be reachable.
