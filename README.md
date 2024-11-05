<p align="center"><img alt="lutho" src="./frontend/src/assets/cognitive.svg" width="32px" height="32px" align="center" /> LUTHO</p>

[![Go](https://github.com/kubeflame/lutho/actions/workflows/go.yml/badge.svg)](https://github.com/kubeflame/lutho/actions/workflows/go.yml)

# App Preview

https://github.com/user-attachments/assets/0827a961-123f-4087-a05f-b4e148e11a0d

If you want to see more demo videos showing indepth actions, you can find them at https://kubeflame.github.io/#demo
 
# Goal

Create a hybrid application that can help in the management of Kubernetes based resources, including Helm charts.

# About

In the kubernetes world Go is a first class language and is tightly coupled with it, therefore it was an obvious pick. 
But, in order to make it available on multiple platforms a web-based framework was the imminent solution, thus Svelte was the choice for the frontend.

# Installation and usage

The instalation notes for the pre-built binaries and Helm chart can be accessed at https://kubeflame.github.io/#install-notes

## Usage Scenarios

### Scenario 1: Using it locally on a desktop or laptop

The application can be started locally from a binary file, and you can connect to a cluster by providing a kubeconfig location or an access token along with the master URL. Then you will need to access the web interface on the configured port -- http://localhost:3001 by default.

`$ lutho start` (and provide the connection details in the interface)

`$ lutho start --port 1337 --kubeconfig /some/location/config`

### Scenario 2: Inside a Kubernetes cluster

The application is configured and deployed using a Helm chart and the interface is accessed using the configured ingress path. (i.e. - http://mycluster.io/lutho/)

In this scenario you can also constrain the access using RBAC and add another layer of security by protecting the ingress using a tool like OAuth2 Proxy.

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
