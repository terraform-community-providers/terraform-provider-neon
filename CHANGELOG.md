## 0.1.9

#### Enhancements
* Added `protected` in `neon_branch` and `neon_project.branch`
* Added `neon_connection_uri` datasource
* Added `history_retention`, `logical_replication` & `allowed_ips` in `neon_project`

## 0.1.8

#### Enhancements
* Added support for Postgres 17

## 0.1.7

#### Enhancements
* Added support for `org_id` in `neon_project`

##### Bug Fixes
* Fixes issue with suspend when creating project on free tier

## 0.1.6

#### Bug Fixes
* Fixes issue with compute provisioner now defaulting to k8s-neonvm

## 0.1.5

#### Enhancements
* Added support for Postgres 16

## 0.1.4

#### Bug Fixes
* Fixes issue with password not appearing when reading role

## 0.1.3

#### Enhancements
* Added `neon_branch` resource
* Added `neon_endpoint` resource
* Added `suspend_timeout` in `neon_project.branch.endpoint`

## 0.1.2

#### Bug Fixes
* Fixes issues with specifying autoscaling for compute endpoints

## 0.1.1

#### Bug Fixes
* Updating a project with the same branch name works now

## 0.1.0 (First release)
