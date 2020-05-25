# Share An Image

After installing the Proffer, lets share our first image. Our first image will be an AWS AMI.

## Pre-requisites

* At least 2 AWS Accounts, if you don't have them, [create](https://aws.amazon.com/free/) them.
* An AMI to share.

## Create Template

Configuration file used to define which AMI we want to share , with which account , to which target region. Format of the template is simple YAML.

Let's create a proffer template file `example.yml` and declare the state.

```YAML
---
resources:
- name: Share AMI With Other Accounts
  type: aws-shareami
  config:
    source:
      # AWS profile that will provide creds for demo-2 account.
      profile: demo-2
      amiFilters:
      # Source AMI name and id.
        name: test-image
        image-id: ami-123456789012
    # Target accounts with which we want to share the AMI.
    target:
      accountRegionMappingList:
          # Target Account 1
        - accountId: 871209123409
          # Target regions for Account-1 account.
          regions:
            - ap-northeast-1

        # Target Account 2
        - accountId: 120923873465
          # Target regions for Account-2 account.
          regions:
            - ap-northeast-2

```

This is a basic template that is ready to go.

> [!NOTE] When applying the template we need to make sure we have AWS creds available on your machine, in this case we are specifying aws profiles. For source account we are using `demo-2` and for target account we are using `demo-1` and `demo-2` aws profiles. Make sure, yous creds are valid. For more info take a look at [aws-docs](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/sharingamis-explicit.html).

Template has one top-level section `resources` that is required in each template. `resources` section contains a list of resources. Resource is a component of proffer that is responsible to do a specific operation. In proffer we have different types of resources, and each resource has an unique purpose to perform.

In this case, in our template we are using a resource of type `aws-shareami`. This resources type is responsible for sharing an AWS AMI from one account to multiple accounts and regions.

Our current template has resource of type `aws-copyami`. This resource has some configuration properties like source and target.
`source` specifies the source of AWS AMI and `target` specifies the target account with which we want to share the source ami.

Each resource has a fixed set of properties. Some are optional, some are required. Detailed documentation for each resource can be found at [Proffer Resources](../resources/README.md).

Before applying the template, let's validate the template by running `proffer validate example.yml`. This command checks the syntax as well as the configuration values to verify they look valid. The output should look similar to below, because the template should be valid. If there are any errors, this command will tell you.

```Bash
$ proffer validate example.yml
validate-syntax | Template syntax is valid.
validate-config | Template config is valid.
```

Once we have valid template, let's apply this template and share the source image to target accounts and regions.

Before applying the template and copy the AMI to target regions, we need to make sure the system from which we are going to apply the template has valid AWS Account Credentials. In this case we are using AWS profiles `demo-2` for `demo2` source aws account and `demo-1`, `demo-3` for `demo-1`, `demo-3` target aws accounts.

After , that let's apply the template by running `proffer apply example.yml` command.

```bash
$ proffer apply example.yml
```

Proffer will share the source AMI with target accounts and regions in parallel.
