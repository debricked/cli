# Callgraph Generation

The Debricked CLI can generate static callgraphs for projects to enable reachability analysis for vulnerabilities.

## Language Support
Debricked CLI callgraph generation currently supports Java and Go.

Java callgraph generation requires compiled classes and supports both the `soot`
and `sootup` engines. For command usage and setup details, see the full CLI
documentation:
https://docs.debricked.com/tools-and-integrations/cli/debricked-cli

## Use

To generate a callgraph for your project you can use the direct command:


```shell
debricked callgraph <path>
```

```shell
debricked callgraph --help
```

Or you can enable it in your scan to add reachability analysis:

```shell
debricked scan --callgraph 
```

To analyze the generated callgraph it needs to be uploaded using the scan command, either with the 
callgraph generation flag as above, or with an already generated call graph by omitting the flag.

For more information, see the full CLI documentation:
https://docs.debricked.com/tools-and-integrations/cli/debricked-cli
