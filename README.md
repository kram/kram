# Kram - A class based scripting language

[![Build Status](https://semaphoreci.com/api/v1/projects/4dd50ed2-e0bd-4694-924e-9d9a14fcae01/487253/shields_badge.svg)](https://semaphoreci.com/zegl/kram-cpp)

**Clarification**, this project is currently just a prototype. Don't actually use it anywhere. Things will break.

## Example

```go
IO::Println("Hello World!")

age := 100

if age > 90 {
    IO::Println("You're old! :)")
}

get_three := fn() {
    1 + 2 * 3 - 16.Sqrt()
}

IO::Println(get_three()) // Prints 3

```

## Features

### Variables

```go
str := "Hi, there!"
```

### If-cases

```go
if first > second {
    IO::Println("first is bigger than second")
} else {
    IO::Println("first is tiny!")
}
```

### Everything is a class

In kram, everything is a class, that means that you can do stuff like

```go
150.Sqrt()
```

## The future of kram

There is a lot of [stuff](https://github.com/kram/kram/labels/Feature) that needs to be implemented before kram is complete. And **you** are very welcome to help! :ok_hand:

An idea of what the language might look like in the future is available at [idea.kr](https://github.com/kram/kram/blob/master/idea.kr).

# License

*kram* is released under a modified 3-clause BSD-license. See [LICENSE](https://github.com/kram/kram/blob/master/LICENSE) for details.