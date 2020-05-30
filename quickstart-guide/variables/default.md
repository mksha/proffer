# Default variables:

To define default variables, we need to define the variables in proffer template file itself.

Follow the below steps to define the default variables:

1. Go to the proffer template file
2. Define the default variables using below syntax

``` yaml
# proffer.yml
---

vars:
  defaultVar1: test
  defaultVar2: test2
  defaultVar3:
    x: 1
  defaultVar4:
    - 1
    - 2

resources:
 ...
```

3. Use the defied variable

``` yaml
# proffer.yml

x: {{var `defaultVar1` }}
```


## Limitations:

All the variables defined in the template configuration are constant variable. they can't contain
template evaluation or any evaluation statements like `env`, `default`.
