# Rubrik SDK for Go Development Guide

## Common Environment Setup

1. Clone the Python SDK repository

      `git clone https://github.com/rubrikinc/rubrik-sdk-for-python.git`

2. Change directory to the repository root dir

     `cd rubrik-sdk-for-python`

3. Switch to the `devel` branch

     `git checkout devel`

4. Create a virtual environment

5. Activate the virtual environment

6. Install the SDK from Source

      `go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm`

## New Function Development
The `/rubrik-sdk-for-go/rubrik_cdm directory` contains all functions for the SDK.

At a high level the directory contains the following:

* `client.go` -  Base API Functions (get, post, etc.) that should only be touched for bug fixes.
* `cloud.go` - Cloud related functions
* `cluster.go` - Functions involving the configuration of the Rubrik cluster itself (think Day 0 configurations)
* `data_management.go` - Functions related to Data Protection tasks (ex. On-demand snapshots)
* `examples_test.go` - Example code utilizing the SDK

When adding a new function, it would ideally fit into one of these files. Each function should have the following:

* Each function must be idempotent. Before making any configuration changes (post, patch, delete) you should first check to see if that change is necessary. If it's not you must return a message formatted as `No change required. {message}:`  For example, the `AssignSLA` function first checks to see if the Rubrik object is already assigned to the provided SLA domain.
* A corresponding example created in `examples_test.go`.

Once a new function has been added you will then submit a new Pull Request which will be reviewed before merging into the `devel` branch.
