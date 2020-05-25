# Troubleshooting

To troubleshoot the errors came during template validation or apply phase, we can enable the logging with debug level, so we can take a look what is going on in background of template validation and apply phase.

To enable the debug level logs, we just need to make sure we pass either `--debug` or `-d` flag to proffer command.

Example:

``` Bash
$ proffer validate --debug example.com
    or
$ proffer validate -d example.com
```

Similarly, we can use the debug flag in apply command also.

``` Bash
$ proffer apply --debug example.com
    or
$ proffer apply -d example.com
```


## How to read the logs more effectively to understand what info they represent

Each message printed by proffer has two parts and they are separated by `|`.

* Prefix :

    It represents the operation type that is currently running. For example

    ``` diff
    + validate-syntax | Template syntax is valid.
    ```

    validate-syntax is representing that, we are running validate-syntax operation.

* Data :
    It represents the message send by the current running operation in a particular log level/ color.
    Data in green color: Information message.
    Data in Bright Green color: Success message.
    Data in Bright Yellow color: Warning/Notice message.
    Data in Red color: Error message.
    Data in Bright Red color: Fatal/Panic message. With this kind of data , program will exit with status code 1.

    ```diff
    - aws-copyami | Invalid AWS AMI ID [ami-07898754ede9e4a342] passed in [config.source.amiFilters] property of Resource: [Copy AMI To Diff AWS Regions]
    ```

    in above logs, aws-copyami tells that we are applying aws-copyami operation and the Data part in Bright Red Color tells that we got an fatal error that caused program to broke with given error message.


## Template debugging:

When we run the proffer `validate` or `apply` command , it generates a `output.yml` file in the same location from where we are running the proffer. We can take a look at this file to check if the template resolved the dynamic information (env vars) or not.
