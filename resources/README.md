# Resources

In Proffer, we have a component called resource. Each resource is responsible to perform a particular sets of tasks. Each resources has its type. We define the resources in `resources` list in proffer template file.

Available Resource Types:

* [aws-copyami](aws/copyami/README.md)
    - Provider: AWS
    - Use case: This resource type can be used to copy AMI in different regions of an AWS Account.
    - [Allowed Properties](aws/copyami/README.md)

* [aws-shareami](aws/shareami/README.md)
    - Provider: AWS
    - Use case: This resource type can be used to share the AMI with different AWS Accounts and regions.
    - [Allowed Properties](aws/shareami/README.md)
