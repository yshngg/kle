# kle [![Go Reference](https://pkg.go.dev/badge/github.com/yshngg/kle.svg)](https://pkg.go.dev/github.com/yshngg/kle) [![GitHub Release](https://img.shields.io/github/v/release/yshngg/kle?style=flat-square)](https://github.com/yshngg/kle/releases)

A Kubernetes Leader Election Demo.

```console
$ go run ./ --help
A Kubernetes Leader Election Demo

Usage:
  kle [flags]
  kle [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Version of kle

Flags:
      --add_dir_header                           If true, adds the file directory to the header of the log messages
      --addr string                              The address kel server binds to.
      --alsologtostderr                          log to standard error as well as files (no effect when -logtostderr=true)
      --client-connection-burst int32            Burst to use for interacting with kubernetes apiserver.
      --client-connection-kubeconfig string      File path to kube configuration for interacting with kubernetes apiserver.
      --client-connection-qps float32            QPS to use for interacting with kubernetes apiserver.
      --dry-run                                  Execute kle in dry run mode.
  -h, --help                                     help for kle
      --kubeconfig string                        File with kube configuration. Deprecated, use client-connection-kubeconfig instead.
      --leader-elect                             Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.
      --leader-elect-lease-duration duration     The duration that non-leader candidates will wait after observing a leadership renewal until attempting to acquire leadership of a led but unrenewed leader slot. This is effectively the maximum duration that a leader can be stopped before it is replaced by another candidate. This is only applicable if leader election is enabled. (default 15s)
      --leader-elect-renew-deadline duration     The interval between attempts by the acting master to renew a leadership slot before it stops leading. This must be less than the lease duration. This is only applicable if leader election is enabled. (default 10s)
      --leader-elect-resource-lock string        The type of resource object that is used for locking during leader election. Supported options are 'leases'. (default "leases")
      --leader-elect-resource-name string        The name of resource object that is used for locking during leader election. (default "kle")
      --leader-elect-resource-namespace string   The namespace of resource object that is used for locking during leader election. (default "demo")
      --leader-elect-retry-period duration       The duration the clients should wait between attempting acquisition and renewal of a leadership. This is only applicable if leader election is enabled. (default 2s)
      --log_backtrace_at traceLocation           when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                           If non-empty, write log files in this directory (no effect when -logtostderr=true)
      --log_file string                          If non-empty, use this log file (no effect when -logtostderr=true)
      --log_file_max_size uint                   Defines the maximum size a log file can grow to (no effect when -logtostderr=true). Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                              log to standard error instead of files (default true)
      --one_output                               If true, only write logs to their native severity level (vs also writing to each lower severity level; no effect when -logtostderr=true)
      --skip_headers                             If true, avoid header prefixes in the log messages
      --skip_log_headers                         If true, avoid headers when opening log files (no effect when -logtostderr=true)
      --stderrthreshold severity                 logs at or above this threshold go to stderr when writing to files and stderr (no effect when -logtostderr=true or -alsologtostderr=true) (default 2)
  -v, --v Level                                  number for the log level verbosity
      --vmodule moduleSpec                       comma-separated list of pattern=N settings for file-filtered logging

Use "kle [command] --help" for more information about a command.
```
