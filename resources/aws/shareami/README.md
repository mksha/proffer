# aws-shareami

This resource is responsible for sharing an AWS AMI from one source region to multiple target accounts and regions across multiple AWS account. This resource also provides some additional features that are listed below:

* Adding extra tags to target ami(s).
* Copy source ami tags to target ami(s).
* Specify AWS Account Credential using multiple creds provider like profile, roleArn and default.
* Specify AWS common target regions for different target accounts.
* Give `CreateVolumePermission` for shared ami to target account.

To start using this resource in proffer template, take a look at the `aws-shareami` resource schema mentioned below:

``` YAML
---
# Schema
name: string        # Required | Desc: Friendly Name of the resource.
type: aws-shareami  # Required | Desc: This value is fixed for aws-shareami resource type.
config: configDict  # Required | Desc: Resource configuration.


# Example : How to use resource in proffer template.
---
resources:
  - name: example resource
    type: aws-shareami
    config: configDict
```

### configDict:

It provides the configuration needed for `aws-shareami` resource to work.

``` YAML
---
# Schema
config:
  source: sourceDict # Required | Desc: Source AMI and Account Information.
  target: targetDict # Required | Desc: Target AMI(s) and Account(s) information.
```

### sourceDict:

It provides the information regarding source ami. This dict object includes information like source ami region, ami-filters to use, how to get aws creds to access source ami.

``` YAML
---
# Schema
source:
  profile: string             # Optional | Desc: AWS Profile to get creds for source ami account.
  roleArn: string             # Optional | Desc: AWS Role ARN to get creds for source ami account
  region: string              # Required | Desc: Source AMI region.
  amiFilters: amiFiltersDict  # Required | Desc: AMI filters to uniquely identify the source ami.
```

> **NOTE:**
    If both `profile` and `roleArn` property not specified then proffer will get the AWS Creds from AWS Default credential providers like environment vars, default profile, aws config file etc.

### amiFiltersDict:

AMI filters to uniquely identify source AMI(s) in different regions.

``` YAML
---
# Schema
amiFilters:
  filerName:filterValue


# Examples
amiFilters:
  name: test-ami
  tag:Purpose: testing
```

### targetDict:

It provides the information about target ami(s) and target AWS accounts(s). This dict object defines with with target accounts we want to share the source ami, in which target regions, we want to share the source ami. It also has some bool flags that can be used to change the behavior of share operation.

``` YAML
---
# Schema
target:
    accountRegionMappingList: [accountRegionMapping]  # Required | Desc: List of accountRegionMapping to specify with which accounts and region we want to share the source ami.
    copyTagsAcrossAccounts: bool                      # Optional | Desc: Flag to indicate if we want to copy the source ami tags to target ami(s) across target accounts.
    addCreateVolumePermission: bool                   # Optional | Desc: Flag to indicate if we want to give `CreateVolumePermission` to target accounts for shared source ami.
    commonRegions: [string]                           # Optional | Desc: List of common target AWS regions with which we want to share the source ami.
```

### accountRegionMapping:

This dict object specifies the information regarding target AWS accounts , regions and how we want to share the ami. This dict object has some boolean flags that can be used to control the ami sharing behavior.

``` Yaml
---
# Schema
accountId: integer                # Required | Desc: Target AWS account id.
accountAlias: string              # Optional | Desc: Target AWS account alias. It is recommended, bec we will get better inventory report. Otherwise accountAlias will be null in report.
profile: string                   # Optional | Desc: AWS Profile to get creds for target aws account. Needed if `copyTagsAcrossAccounts` flag is true and `roleArn` key is not set.
roleArn: string                   # Optional | Desc: AWS Role ARN to get creds for target aws account. Needed if `copyTagsAcrossAccounts` flag is true and `profile` key is not set.
regions: [string]                 # Required | Desc: List of target AWS account regions with which we want to share the source ami.
copyTagsAcrossAccounts: bool      # Optional | Desc: Flag to indicate if we want to copy the source ami tags to target ami across target account. If this flag is true then make sure either `profile` or `roleArn` key is specified for target aws account creds.
addCreateVolumePermission: bool   # Optional | Desc: Flag to indicate if we want to give `CreateVolumePermission` to target account for shared source ami.
addExtraTags: tags                # Optional | Desc: AWS EC2 tags to add to target ami.
```

### tags:

AWS EC2 tags to add to the EC2 resources like AMI, etc. In this case these tags are added to the target ami(s).

``` YAML
---
# Schema
tags:
  string: string


# Examples
tags:
  work: test
  job: test
```

## Complete Example:

It shows how to use `aws-shareami` resource type in proffer template with all configuration possible.

```YAML
---
resources:
- name: Share AMI With Other Accounts
  type: aws-shareami
  config:
    source:
      profile: demo-2
      amiFilters:
        name: test-image
        tag:Purpose: Testing
        tag:Department: DevOps
    target:
      accountRegionMappingList:
        - accountId: 121616226324
          accountAlias: demo-1
          profile: demo-1
          regions:
            - ap-northeast-1
          copyTagsAcrossAccounts: true
          addExtraTags:
            CreatedBy: SharedByDemo1
            Type: AMITesting
            Home: {{ env "HOME" | default "default value" }}
        - accountId: 121266418583
          accountAlias: demo-3
          regions:
            - ap-northeast-2
      addCreateVolumePermission: true
      commonRegions:
      - us-east-1
      - us-west-2

```
