# Callgraph Generation

The Debricked CLI can generate static callgraphs for projects to enable reachability analysis for vulnerabilities.

For the only currently supported language, Java, some setup is required for callgraph generation to work
properly. For more information on this see the Language Support section below.

## Language Support
Debricked CLI callgraph generation currently only supports Java, the specific
documentation for the Java callgraph generation can be
found [here](https://github.com/debricked/cli/blob/main/internal/callgraph/language/java11/README.md).

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

For more information see documentation on the specific langauge implementation or see full CLI documentation [here](https://docs.debricked.com/tools-and-integrations/cli/debricked-cli)
