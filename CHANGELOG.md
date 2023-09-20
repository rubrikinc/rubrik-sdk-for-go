# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Types of changes

- **Added** for new features.
- **Changed** for changes in existing functionality.
- **Deprecated** for soon-to-be removed features.
- **Removed** for now removed features.
- **Fixed** for any bug fixes.
- **Security** in case of vulnerabilities.


## v1.0.3

### Added

- `ConnectAPIToken()` ([Issue 40](https://github.com/rubrikinc/rubrik-sdk-for-go/issues/40))
- `RecoverFileDownload()` 


### Changed

- `ConnectEnv()` will not look for a configured `rubrik_cdm_token` environment variable

## v1.0.4

### Fixed
- Prevent potential HTTP connection issues from occurring on long-running jobs

## v1.1.0

### Added

- `BootstrapCcesAws()` will bootstrap a Cloud Cluster Elastic Storage in AWS
- `BootstrapCcesAzure()` will bootstrap a Cloud Cluster Elastic Storage in Azure

### Fixed

- `Bootstrap()` fixed NTP code so that NTP will be updated properly on newer Rubrik clusters. This is a breaking change for Rubrik clusters older than v5.0.

## 1.1.1

### Fixed

- Fixed `module declares its path as: rubrikcdm but was required as: github.com/rubrikinc/rubrik-sdk-for-go` error when calling the SDK from github.

## 1.2

### Added

- Added support for bootstrapping CCES on AWS and Azure with the immutability flag. 

### Fixed

- Fixed the examples file for encryption.
- Fixed the .gitignore file to ignore test files.


## 1.2.1

- Added retry on connection errors during bootstrap 
## 1.3.0

- Upgraded to Go 1.21