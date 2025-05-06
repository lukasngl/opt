[![Go Reference](https://pkg.go.dev/badge/github.com/lukasngl/opt.svg)](https://pkg.go.dev/github.com/lukasngl/opt)

# Yet another option type for go (YAOTTFG)

Optionality if often *implicitly* represented by overloading the semantics of
pointers or zeroness.
This package exports the `opt.T` type and aliases for built-in types,
to *explicitly* represent optional values.

## Install

```sh
go get github.com/lukasngl/opt
```

## Features

- **concise**: Seriously, this type isn't much more obtrusive than the common
  pointer method.

  ```go
  type Query struct {
  	OrderBy opt.T[OrderBy] `json:"where"`
  	Limit   opt.Int        `json:"limit"`
  	Offset  opt.Int        `json:"offset"`
  }
  ```

- **idiomatic**: Multiple return values look and feel at home between common
  patterns like error handling, map lookup, or type casting:

  ```go
  value, ok := optional.Unpack()
  if !ok {
    return fmt.Errorf("not present :(")
  }
  ```

- **simple**:
  - **zero schnickschnak**, that you've seen in the functional world
    and would really like to use (sniff), but that feel cumbersome
    and look weird in your go code.
  - **zero depencies**

- **compatible**:
  - **pointers**: Seamlessly integrates with common pointer based optionality,
    via `FromNillable` and `ToNillable`.
  - **zero**: Works with zeroness via `FromZero` and `OrZero`,
    honoring `IsZero() bool` methods when implemented.
  - **json**: Implements `json.Unmarshaller` and `json.Marsaller`,
    with support for the `omitzero` tag introduced in [go1.24].

    Note: the package itself only requires go>=1.18 for generics,
    thus the omitzero tests are in a separate module, that requires go1.24.
  - **sql**: Implements `driver.Valuer` and `driver.Scanner`,
    by delegating to `sql.Null`.
  - ~~**xml**:~~ PRs welcome, did not have a use case yet.

[go1.24]: https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson

## Prior Art

As the title suggest, there are loads of other packages,
but all of them had at least one of the following deal-breakers:

- Really long type names, I mean the typescript folks just need to suffix
  question mark, so why should we suffer?
- A slice based approach, this forces the value onto the heap
  and stores length and capacity, leading to an unnecessary overhead.
- No `IsZero() bool` method for the `omitzero` tag.

## Schnickschnak

Here is the Map and Filter method, in case you need it,
or I get weak and decide to add them later:

```go

func Map[A, B any](o opt.T[A], f func(A) B) T[B] {
	a, present := o.Unwrap()
	if !present {
		return None[B]()
	}

	return Some(f(a))
}

func Filter[A any](o opt.T[A], keep func(A) bool) T[A] {
	value, present := o.Unwrap()
	if !present || !keep(value){
		return None[A]()
	}

	return o
}
```
