# Dynamic variables

To define a dynamic variables, we need to define the variables in separate variables file.

Follow the below steps to define the dynamic variables:

1. Create variable file. For example `data.yml`.
2. Define the dynamic variables using below syntax

``` yaml
dynamicVar1: test
dynamicVar2: test2
dynamicVar3:
  x: 1
dynamicVar4:
  - 1
  - 2
dynamicVar5:
  home: {{ env "HOME" | default "my default value" }}

```

3. Use the defied variable in proffer template

``` yaml
x: {{var `dynamicVar1` }}
```


## Advantage

All the variables defined in the variable file are runtime variable. they can contain
template evaluation or any evaluation statements like `env`, `dynamic`.
