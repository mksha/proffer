# Proffer

[![codecov](https://codecov.io/gh/mohit-kumar-sharma/proffer/branch/master/graph/badge.svg?token=YFU0AS3HEJ)](https://codecov.io/gh/mohit-kumar-sharma/proffer)
![Tests](https://github.com/mohit-kumar-sharma/proffer/workflows/Tests/badge.svg)

Proffer is a cross platform command line tool to copy and share the cloud images across multiple regions and accounts. It is a lightweight tool that can be used on any major platform.

## Supported Cloud Providers

* AWS

Support for the other cloud providers can be added via resource plugin.

## How Proffer works

Proffer command takes a template called `proffer.yml` written in yaml format and apply the resources defined in template. Each proffer template has a top-level section called `resources` that is list of proffer resources.
Each resource then have their own properties like type, keys and etc. To find all available proffer resources, visit [Available Proffer Resources](resources/README.md) page.

## Quick Start

To quickly start with proffer , you can follow the [quick-start-guide](quickstart-guide/main.md).


## Resources Available In Proffer

Resource is a component in proffer. Each resources is responsible to perform a particular set of operations. Proffer has different kinds of resources.For more details, check [Available Proffer Resources](resources/README.md).


## Access Environment Variables

To access the environment variables within proffer template , we can use below format:

```Yaml
Home: {{ env "HOME" }}
```

If we want to set default value of a environment variable if its not set then we can use below format:

```Yaml
Home: {{ env "HOME" | default "default home dir path" }}
```


## License

Proffer is released under the Apache License, Version 2.0.
