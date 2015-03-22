# Gus - A class based scripting language

**Clarification**, this project is currently just a prototype. Don't actually use it anywhere. Things will break.

## Example

```dart
IO.Print("Hello World!")

var Age = 100

if Age > 90 {
    IO.Print("You're old! :)")
}

class Magic {
    Say() {
        IO.Print("Yoho!")
    }
}

var yolo = 1 + 2 * 3 - 4 // 3
```

## Features

### Variables

```dart
var Str = "Hi, there!"
```

### If-cases

```dart
if A > B {
    IO.Println("A is bigger than B")
} else {
    IO.Println("A is tiny!")
}
```

### Classes

```dart
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

### Pretty numbers

You can seperate numbers by spaces (as many or as few as you like) to increase readability of the sourcecode.

```dart
IO.Println(20 000) // Prints "20000"
```

## The future of Gus

There is a lot of [stuff](https://github.com/zegl/Gus/labels/Feature) that needs to be implemented before Gus is complete. And **you** are very welcome to help! :ok_hand:
