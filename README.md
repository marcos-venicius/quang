# Quang

`Quang` is meant to be a "query" language for your tools.

`Quang` extends for `Query Language`.

You can use the language to implement filters in your applications without the need of tons of flags to your user,
just provide a simple and consise integrated language and let the user build the filter himself.

So, it is at the end a _"language as a lib"_.

_There is a full [playlist on youtube](https://youtube.com/playlist?list=PL3YefAkg_zCgPTINetXJ7Aatpps2DECcW&si=_jBBPe7s6Ne4B70s) creating this project._

## Api and language details

Does not matter what kind of query you provide as input to the evaluator, it will always return `true` or `false`. If the query is empty, it will always return `true`.

**Data Types**

| name     | supported | format        | description                                                                   |
| -------- | --------- | ------------- | ----------------------------------------------------------------------------- |
| Integers | yes       | `[0-9]+`      | golang 64bit signed integers                                                  |
| Atoms    | yes       | `:[a-zA-Z_]+` | it works like enumerators                                                     |
| String   | yes       | `'.*'`        | you can scape string with `\'`                                                |
| Boolean  | yes       | `true\|false` |                                                                               |
| Nil      | yes       | `nil`         | represents all kinds of empty values ("", nil) (zero is not considered empty) |
| Floats   | yes       | `\d+\.\d*`    | golang 64bit floats                                                           |

**Keywords**

| name     | description   | usage            |
| -------- | ------------- | ---------------- |
| and      | logical and   | `true and true`  |
| or       | logical or    | `true or false`  |
| nil      | null value    | `name eq nil`    |
| true     | boolean true  | `alive eq true`  |
| false    | boolean false | `alive eq false` |

**Operators**

| name     | description                                                                                        | example                 |
| -------- | -------------------------------------------------------------------------------------------------- | ----------------------- |
| eq       | check if `a` is equal to `b`. strict types. (Integers, Strings, Booleans, Nils, Floats, Atoms)     | `a eq b`                |
| ne       | check if `a` is not equal to `b`. strict types. (Integers, Strings, Booleans, Nils, Floats, Atoms) | `a ne b`                |
| lt       | check if `a` is less than `b`. strict types. (Integers, Strings)                                   | `a lt b`                |
| gt       | check if `a` is greater than `b`. strict types. (Integers, Strings)                                | `a gt b`                |
| lte      | check if `a` is less than or equal to `b`. strict types. (Integers, Strings)                       | `a lte b`               |
| gte      | check if `a` is greater than or equal to `b`. strict types. (Integers, Strings)                    | `a gte b`               |
| reg      | check if `a` matches pattern `b`. `b` accepts valid regex. `a` should be a string                  | `a reg b`               |

**Basic syntax**

Pretend we have a list of computers that have the following properties:

- Identifier
- Running
- Cors

So, we could query something like:

```elixir
(running eq true and cors gte 4 and cors lte 10) or (running eq false and identifier reg 'ML-\d+') or identifier eq nil
```
