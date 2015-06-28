# Kram - A class based scripting language

[![Build Status](https://travis-ci.org/kram/kram.svg?branch=master)](https://travis-ci.org/kram/kram)

**Clarification**, this project is currently just a prototype. Don't actually use it anywhere. Things will break.

## Example

```go
IO.Print("Hello World!")

var age = 100

if age > 90 {
    IO.Print("You're old! :)")
}

class Magic {
    Say() {
        IO.Print("Yoho!")
    }
}

var three = 1 + 2 * 3 - 16.Sqrt() // 3
```

Are you interested? Make a visit to the [Getting Started](https://github.com/kram/kram/wiki) page!

## Features

### Variables

```go
var str = "Hi, there!"
```

### If-cases

```go
if first > second {
    IO.Println("first is bigger than second")
} else {
    IO.Println("first is tiny!")
}
```

### Classes

```go
class Magician {
    
    // Instance variables
    var name

    Name(name) {
        self.name = name
    }
    
    // Uppercase -> Public method
    Say {
        IO.Println("My name is " + self.name)
    }
    
    // Static methods are not a part of the class instance
    static Woho {
        IO.Println("Woho!")
    }
}

var Harry = new Magician("Harry")
Harry.Say() // My name is Harry

Harry.Woho() // Woho!
Magician.Woho() // Woho!
```

#### Default parameters

```go
class Magician {
    Setup(name = "Harry") {
        // Stuff here
    }
}
```

#### Named arguments

```go
class Magician {
    Setup(name = "Harry", age = 14) {
        // Stuff here
    }
}

var magic = new Magician()
magic.Setup(age = 11, name = "Ron") // ("Ron", 11)
magic.Setup(age = 7) // ("Harry", 7)
```

### Pretty numbers

You can seperate numbers by spaces (as many or as few as you like) to increase readability of the sourcecode.

```go
IO.Println(20 000) // Prints "20000"
```

### Everything is a class

In kram, everything is a class, that means that you can do stuff like

```go
150.Sqrt()
```

## The future of kram

There is a lot of [stuff](https://github.com/kram/kram/labels/Feature) that needs to be implemented before kram is complete. And **you** are very welcome to help! :ok_hand:

# License

*kram* is released under a modified 3-clause BSD-license. See [LICENSE](https://github.com/kram/kram/blob/master/LICENSE) for details.
