# aws-copyami

This resource is responsible for copying an AWS AMI from one source region to multiple target regions within same AWS account. This resource also provides some additional features that are listed below:

* Adding extra tags to target ami(s).
* Copy source ami tags to target ami(s).
* Specify AWS Account Credential using multiple creds provider like profile, roleArn and default.

To start using this resource in proffer template, take a look at the `aws-copyami` resource schema mentioned below:

``` YAML
---
# Schema
name: string       # Required | Desc: Friendly Name of the resource.
type: aws-copyami  # Required | Desc: This value is fixed for aws-copyami resource type.
config: configDict # Required | Desc: Resource configuration.


# Example : How to use resource in proffer template.
---
resources:
  - name: example resource
    type: aws-copyami
    config: configDict
```

### configDict:

It provides the configuration needed for `aws-copyami` resource to work.

``` YAML
---
# Schema
config:
  source: sourceDict # Required | Desc: Source AMI and Account Information.
  target: targetDict # Required | Desc: Target AMI(s) and Account information.
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

> [!NOTE] If both `profile` and `roleArn` property not specified then proffer will get the AWS Creds from AWS Default credential providers like environment vars, default profile, aws config file etc.

### amiFiltersDict:

AMI filter to uniquely identify an AMI in a region.

``` YAML
---
# Schema
amiFilters:
  filerName:filterValue


# Examples
amiFilters:
  name: test-ami
  image-id: ami-123456789012
  tag:Purpose: testing
```

### targetDict:

It provides the information about target ami(s). This dict object defines in which target regions, we want to copy the source ami. It also has some bool flags that can be used to change the behavior of copy operation.

``` YAML
---
# Schema
target:
  regions: [string]            # Required | Desc: List of target AWS regions where we want to copy the source ami.
  copyTagsAcrossRegions: bool  # Optional | Desc: Flag to indicate if we want to copy the source ami tags to target ami(s).
  addExtraTags: tags           # Optional | Desc: AWS EC2 tags to add to target ami(s).
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

It shows how to use `aws-copyami` resource type in proffer template with all configuration possible.

```YAML
---
resources:
- name: Copy AMI To Multiple Regions
  type: aws-copyami
  config:
    source:
      profile: demo-2
      region: us-east-1
      amiFilters:
        name: test-image
        tag:Purpose: Testing
        tag:Department: DevOps
    target:
      regions:
      - ap-northeast-1
      - ap-northeast-2
      - us-west-2
      copyTagsAcrossRegions: true
      addExtraTags:
        CreatedBy: local-testing-tool
        Type: testing

```
