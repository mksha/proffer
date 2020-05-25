# Copy An Image

After installing the Proffer, lets copy our first image. Our first image will be an AWS AMI.

## Pre-requisites

* An AWS Account, if you don't have one, [create](https://aws.amazon.com/free/) one.
* An AMI to copy.

## Create Template

Configuration file used to define which AMI we want to copy , from which source region , to which target region. Format of the template is simple YAML.

Let's create a proffer template file `example.yml` and declare the state.

``` YAML
---
resources:
- name: Copy AMI To Diff AWS Regions
  type: aws-copyami
  config:
    source:
      # Source AMI region.
      region: us-east-1
      amiFilters:
        # Source AMI id and name.
        image-id: ami-3481274ede9e4a3
        name: test-image
    target:
      # Target regions where we want to copy the source AMI.
      regions:
      - ap-northeast-1
      - ap-northeast-2
      - us-west-2
```

This is a basic template that is ready to go.

> [!NOTE] When applying the template we need to make sure we have AWS creds available on your machine, in this case we are not specifying any cred provider name so proffer will get AWS Creds from AWS Env Vars. Make sure, yous creds have permission to copy source ami. For more info take a look at [aws-docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/CopyingAMIs.html).

Template has one top-level section `resources` that is required in each template. `resources` section contains a list of resources. Resource is a component of proffer that is responsible to do a specific operation. In proffer we have different types of resources, and each resource has an unique purpose to perform.

In this case, in our template we are using a resource of type `aws-copyami`. This resources type is responsible to copy an AWS AMI from one region to multiple regions.

Each resource has a fixed set of properties. Some are optional, some are required. Detailed documentation for each resource can be found at [Proffer Resources](../resources/README.md).

Our current template has resource of type `aws-copyami`. This resource has some configuration properties like source and target.
`source` specifies the source of AWS AMI and `target` specifies the target where we want to copy the source ami.

Before applying the template, let's validate the template by running `proffer validate example.yml`. This command checks the syntax as well as the configuration values to verify they look valid. The output should look similar to below, because the template should be valid. If there are any errors, this command will tell you.

```Bash
$ proffer validate example.yml
validate-syntax | Template syntax is valid.
validate-config | Template config is valid.
```

Once we have valid template, let's apply this template and copy the source image to target regions.

Before applying the template and copy the AMI to target regions, we need to make sure the system from which we are going to apply the template has valid AWS Account Credentials. In this case as AWS ENv Vars.

After , that let's apply the template by running `proffer apply example.yml` command.

```bash
$ proffer apply example.yml

start-validation| Validating template before applying...
validate-syntax | Template syntax is valid.
validate-config | Template config is valid.

start-apply | Applying template config...
aws-copyami | Resource : Copy AMI To Diff AWS Regions  Status: Started
aws-copyami | 
aws-copyami | Started Copying AMI In Account: 12345678901 Region: ap-northeast-1 ...
aws-copyami | Started Copying AMI In Account: 12345678901 Region: ap-northeast-2 ...
aws-copyami | Started Copying AMI In Account: 12345678901 Region: us-west-2 ...
aws-copyami | Copied AMI In Account: 12345678901 In Region: ap-northeast-1 , New AMI Id Is: ami-0347a3dc51f46491d
aws-copyami | Copied AMI In Account: 12345678901 In Region: ap-northeast-2 , New AMI Id Is: ami-0dd435a3959fb57e4
aws-copyami | Copied AMI In Account: 12345678901 In Region: us-west-2 , New AMI Id Is: ami-09ff2a7d34a6bc60c
aws-copyami | 
aws-copyami | Resource : Copy AMI To Diff AWS Regions  Status: Succeeded

```

Proffer will copy the source AMI in target region in parallel.

This resource is also idempotent, so if we apply the same configuration again then we it will not do anything.

```Bash
$ proffer apply example.yml

start-validation| Validating template before applying...
validate-syntax | Template syntax is valid.
validate-config | Template config is valid.

start-apply | Applying template config...
aws-copyami | Resource : Copy AMI To Diff AWS Regions  Status: Started
aws-copyami | 
aws-copyami | AMI test-image Already Exist In Account 12345678901 In Region ap-northeast-1
aws-copyami | AMI test-image Already Exist In Account 12345678901 In Region ap-northeast-2
aws-copyami | AMI test-image Already Exist In Account 12345678901 In Region us-west-2
aws-copyami | 
aws-copyami | Resource : Copy AMI To Diff AWS Regions  Status: Succeeded

```
