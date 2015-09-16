# Kram - A class based scripting language

[![Build Status](https://semaphoreci.com/api/v1/projects/e793760f-a344-47ce-bbf4-0af68745f97f/491932/shields_badge.svg)](https://semaphoreci.com/zegl/kram)

**Clarification**, this project is currently just a prototype. Don't actually use it anywhere. Things will break.

## Example

```go
IO::Println("Hello, 🌍")

// Everything is a class
150.Sqrt()

class Dog {
	Bark := fn() {
		IO::Println("Woff!")
	}
}

snoopy := new Dog()
snoopy.Bark()
```

## The future of kram

There is a lot of [stuff](https://github.com/kram/kram/labels/Feature) that needs to be implemented before kram is complete. And **you** are very welcome to help! :ok_hand:

An idea of what the language might look like in the future is available at [idea.kr](https://github.com/kram/kram/blob/master/idea.kr).

# License

*kram* is released under a modified 3-clause BSD-license. See [LICENSE](https://github.com/kram/kram/blob/master/LICENSE) for details.
