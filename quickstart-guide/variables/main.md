# Variables

## Define dynamic and default variables

Before using the variables in proffer template file, we ned to define those variables. We can define those variables at below places:

* [In a separate variables file](dynamic.md)
* [In same proffer template config file](default.md)

If the variable is dynamic, means can change the values for different environments and cases then we can define that variables in different variables file and then pass that variable file to proffer cli using `--var-file` option.

If the variable is constant then we can just define its values in proffer template file itself and use it directly without passing as var file.

## Use variables in proffer configuration file

Once we have defined the variables in either variable file or template file, then we can use them using below syntax:

``` yaml
x : {{ var `var_name` }}
```

## Limitations

Nesting of variables are not allowed on both variable and template file.

Ex: Invalid variable definition.

```
x:
  s:
    f: 2
```
