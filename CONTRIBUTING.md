# Rubrik SDK for Go Development Guide

Contributions via GitHub pull requests are gladly accepted from their original author. Along with any pull requests, please state that the contribution is your original work and that you license the work to the project under the project's open source license. Whether or not you state this explicitly, by submitting any copyrighted material via pull request, email, or other means you agree to license the material under the project's open source license and warrant that you have the legal authority to do so.

## Common Environment Setup - Microsoft Windows

1. Create a new directory for development work and change to that directory
```
md dev
cd dev
```
2. Change your `GOPATH` environment variable to the current directory
```
set "GOPATH=%cd%"
```
3. Create the necessary directory structure
```
md src\github.com\rubrikinc
```
4. Change to the newly created directory
```
cd src\github.com\rubrikinc
```
5. Clone the Rubrik SDK for Go repository
```
git clone https://github.com/rubrikinc/rubrik-sdk-for-go.git
```
6. Change to the repository root directory
```
cd rubrik-sdk-for-go
```
7. Switch to the devel branch
```
git checkout devel
```

## Common Environment Setup - macOS and \*nix

1. Create a new directory for development work and change to that directory
```
md dev
cd dev
```
2. Change your GOPATH environment variable to the current directory
```
export GOPATH=`pwd`
```
3. Create the necessary directory structure
```
mkdir -p src/github.com/rubrikinc
```
4. Change to the newly created directory
```
cd src/github.com/rubrikinc
```
5. Clone the Rubrik SDK for Go repository
```
git clone https://github.com/rubrikinc/rubrik-sdk-for-go.git
```
6. Change to the repository root directory
```
cd rubrik-sdk-for-go
```
7. Switch to the devel branch
```
git checkout devel
```

## New Function Development
The `/rubrik-sdk-for-go/rubrikcdm` directory contains all functions for the SDK.

At a high level the directory contains the following:

* `client.go` -  Base API Functions (get, post, etc.) that should only be touched for bug fixes.
* `cloud.go` - Cloud related functions
* `cluster.go` - Functions involving the configuration of the Rubrik cluster itself (think Day 0 configurations)
* `data_management.go` - Functions related to Data Protection tasks (ex. On-demand snapshots)
* `examples_test.go` - Example code utilizing the SDK

When adding a new function it ideally should be categorized to fit into one of the above files. Each function should meet the following requirements:

* Each function must be idempotent. Before making any configuration changes (post, patch, delete) you should first check to see if that change is necessary. If it's not you must return a message formatted as `No change required. {message}`. For example, the `AssignSLA()` function first checks to see if the Rubrik object is already assigned to the provided SLA domain.
* A corresponding example created in `examples_test.go`.

Once a new function has been added you will then submit a new Pull Request which will be reviewed before merging into the devel branch.
* A corresponding example created in `examples_test.go`.

Once a new function has been added you will then submit a new Pull Request which will be reviewed before merging into the `devel` branch.

