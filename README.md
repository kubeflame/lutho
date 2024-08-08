<p align="center"><img alt="lutho" src="./frontend/src/assets/cognitive.svg" width="32px" height="32px" align="center" /> LUTHO</p>

[![Go](https://github.com/kubeflame/lutho/actions/workflows/go.yml/badge.svg)](https://github.com/kubeflame/lutho/actions/workflows/go.yml)

# Goal

Create a hybrid application that can help in the management of Kubernetes based resources, including Helm charts.

# About

In the kubernetes world Go is a first class language and is tightly coupled with it, therefore it was an obvious pick. 
But, in order to make it available on multiple platforms a web-based framework was the imminent solution, thus Svelte was the choice for the frontend.

# Installation and usage

The instalation notes for the pre-built binaries can be accessed at https://kubeflame.github.io/#install-notes

## Building it from source

### Requirements
 - go
 - nodejs (yarn)

### One line command to execute
```bash
go generate ./... && go build -ldflags="-s -w -X 'main.AppVersion=<provide-a-version>'"
```

## Starting the application

```bash
./lutho -h
```

```
NAME:
   lutho - Manage Kubernetes based resources in a different way

USAGE:
   lutho [global options] command [command options]

VERSION:
   dev

AUTHOR:
   KubeFlame <https://github.com/kubeflame>

COMMANDS:
   start  Starts the application

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### There are some start options available for providing the kubernetes configuration

```bash
./lutho start -h
```

```
NAME:
   lutho start - Starts the application

USAGE:
   lutho start [command options]

DESCRIPTION:
   This command will start the application

OPTIONS:
   --port value          the port where the UI can be accessed at (default: "3001")
   --kubeconfig value    kubeconfig file path (default: "/Users/alex/.kube/config")
   --access-token value  kubernetes cluster access token
   --master-url value    kubernetes master URL
   --in-cluster          set this flag if the app is inside the kubernetes cluster (default: false)
   --help, -h            show help
```

Once the application is started `lutho start <options>`, by default it is available at `localhost:3001`.
