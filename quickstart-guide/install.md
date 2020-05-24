# Install Proffer Cli

To install proffer cli, download platform specific binary from [release-page](github.com).

If you intend to access Proffer from the command-line, make sure to place it somewhere on your `PATH` .

After installing Proffer, verify the installation worked by opening a new command prompt or console, and checking that proffer is available:

``` Bash
$ proffer

Proffer is a command-line tool to distribute cloud images in between multiple regions
and with multiple environments. This tool only needs a yml configuration file with name proffer that defines
the image distribution operations and actions. Currently AWS cloud is the only supported cloud provider but
support for other cloud providers can be added via resource plugin.

Usage:
  proffer [command]

Examples:

$ proffer [command] [flags] TEMPLATE
$ proffer validate proffer.yml
$ proffer validate -d proffer.yml
$ proffer apply proffer.yml

Available Commands:
  apply       Apply proffer configuration file.
  completion  Generates shell completion script for specified shell type
  gendoc      Generate proffer markdown documentation
  help        Help about any command
  validate    Validate proffer configuration file.

Flags:
  -d, --debug     Set debug flag to get detailed logging
  -h, --help      help for proffer
  -v, --version   version for proffer

Use "proffer [command] --help" for more information about a command.
```

If you get an error that proffer could not be found, then your `PATH` environment variable was not setup properly. Please go back and ensure that your `PATH` variable contains the directory which has proffer installed.

## Enable Auto Completion

Proffer also provide a command called `completion` through which we can generate the bash/zsh auto-completion scripts.

To set-up the base/zsh completion take a look at [setup auto-completion](../doc/proffer_completion.md) page.
