# Stream Exporter

[![CircleCI](https://circleci.com/gh/carlpett/stream_exporter.svg?style=shield)](https://circleci.com/gh/carlpett/stream_exporter)

A [Prometheus](https://prometheus.io) exporter for extracting metrics from streaming sources of text, such as logs written to a socket, or tailing a file.
Extraction into metrics is done using regular expressions, on a per-line basis. Capture groups in the regular expression are used as labels in the metrics. All Prometheus metric types are supported, and fully configurable.


# Inputs and patterns
The exporter works by reading from an input module, and then processing each line read with a number of configured metric patterns. Matching lines increment/set/observe their corresponding metric.

When testing configuration for the stream exporter, the `file` input module has a `dryrun` mode which runs an entire file through the configuration, and then prints the resulting metric output before exiting. See below for more information.

The exporter has multiple input modules, and one must be selected when starting the exporter using the `-input.type` flag. A list of all supported modules on your platform is returned when the `-input.print` flag is given. Details about the modules and their configuration are given below.

Metric patterns are configured with a yaml file, and consist of a listing all the patterns used for extraction. As an example, here is a simple configuration with two metrics:

```yaml
metrics:
  - name: comments_total
    kind: counter
    pattern: "^#.+"
  - name: current_visitors
    kind: gauge
    pattern: "Current number of visitors on (?P<section>\w+) is (?P<value>[0-9]+)"
```

The first one will create a counter `comments_total`, which will be incremented for every row starting with a `#`. This metric will not have any labels.

The second pattern is a gauge, with two capture groups, `section` and `value`. `value` is treated specially, and will be the value set on the gauge, while `section` will become a label on the metric. The output will be a metric like `current_visitors{section="landing-page"} 432`.


# Flags
`stream_exporter` accepts a number of flags, most of which are related to the input modules. Those flags are described separately in the Inputs section below. The flags not configuring a specific input module are these:

Name     | Description | Default value
---------|-------------|--------------------
`-config` | Path to the pattern configuration file | `stream_exporter.yaml`
`-input.type` | What input module to use | _(Required)_
`-input.print` | Prints available input modules and exits | -
`-web.listen-address` | Address on which to expose metrics | `:9178`
`-web.metrics-path` | Path under which the metrics are available | `/metrics`


# Pattern configuration

The pattern configuration is a yaml file consisting of a top level list, `metrics`. The elements in this list have the following structure:

```yaml
  - name: metric_name
    kind: metric_kind
    pattern: regular_expression
    <metric kind specific options>
```

The `name` field will become the name of the metric.

The `kind` field is the type of the metric, valid values are `counter`, `gauge`, `summary` and `histogram`.

The `pattern` field is a [Golang regexp](https://golang.org/pkg/regexp/) matching the lines you want to find. Named capture groups (`(?P<name>CAPTURING_REGEXP)`) are turned into labels. Exception to this is having a label `value` for gauges, summaries and histogram, where the matched number is what is set or observed (depending on the metric kind) for the metric.

## Specific options
In addition to the above, histograms and summaries have a number of options that can be set.

### Histogram
For histograms, there is a single option, `buckets`, that can be configured.

- `buckets`: A list of floating point numbers defining the upper bounds of each bucket 

Example:

```yaml
  - name: my_histogram
    kind: histogram
    pattern: I saw (?P<value>\d+) foos!
    buckets:
      - 1
      - 5
      - 10
      - 20
```

### Summary
Summaries have a number of options, `objectives`, `max_age`, `age_buckets` and `buf_cap`. 

- `objectives`: A dictionary with floating point numbers as both keys and values. Defines the quantile rank estimates and their respective absolute error.
- `max_age`: Defines the duration for which an observation stays relevant for the summary. Given as a string representation of a duration, for example `120s`, `20m` or `3d`.
- `age_buckets`: An integer. defining the number of buckets used to exclude observations that are older than MaxAge from the summary. A higher number has a resource penalty, so only increase it if the higher resolution is really required.
- `buf_cap`: The sample stream buffer size. If there is a need to increase the value, a multiple of 500 is recommended.

Example:

```yaml
  - name: my_summary
    kind: summary
    pattern: Found a foo in (?P<value>[\d\.]+) seconds!
    objectives:
      0.5: 0.5
      0.9: 0.1
      0.99: 0.01
    max_age: 15m
    age_buckets: 3
    buf_cap: 500
```


# Input modules
A number of input modules are available, described here with their parameters.

## File
Reads lines from a file.

The `file` input has two parameters:

- `-input.file.path` (Required) : The path to the file to read from.
- `-input.file.mode` (Optional, default `tail`): `tail` or `dryrun`. In `tail` mode, the file is opened and any new lines written to the file is processed. In `dryrun` mode, the entire file is read and processed. When the end of the file is encountered, the current metrics state is written to standard out, and the exporter exits. Useful for debugging and configuration testing.

## Socket
Opens a socket for reading lines. 

The `socket` input has two parameters:

- `input.socket.family` (Optional, default `tcp`): The address family to open the socket. Valid values depend on platform, but include `tcp` and `udp`. On Linux and certain other systems, `unix` is also available.
- `input.socket.listenaddr` (Required): The listen specification, for example `:10000` to listen on all interfaces on port 10000.

## Syslog
Creates a syslog server which can act as a remote server from the local syslog server.

The `syslog` input has three parameters:

- `input.syslog.listenfamily` (Required): The address family of the server. Valid values depend on platform, but include `tcp` and `udp`. On Linux and certain other systems, `unix` is also available.
- `input.syslog.listenaddr` (Required): The listen specification, for example `:1514` to listen on all interfaces on port 1514.
- `input.syslog.format` (Optional, default `autodetect`): The format of the syslog messages. Defaults to automatic detection, can also be set to `rfc3164`, `rfc5424` or `rfc6587`.

## Named pipe
Creates a named pipe to which lines can be written. Only available on Linux.

The `namedpipe` input has one parameter:

- `input.namedpipe.path` (Required): The path where the pipe should be created. This may require elevated privileges to execute a `mkfifo` syscall.


# Building
The project uses [govendor](https://github.com/kardianos/govendor) for dependency management. To build the exporter, call `govendor build +p`.
