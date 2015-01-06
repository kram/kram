# Gus - A class based scripting language

TL:DR; Very little is actually implemented.

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

## Future (?)

### For, List, Range

```dart
for abc in ["a", "b", "c"] {
    IO.Print(abc)
}

for num in 1..100 {
    IO.Print(num)
}
```

### Functions

```dart
var myFunction = new Fn(a, b, c) {
    IO.Println(a, b, c)
}
```

### Classes

```dart
class Magician {
    
    // Lowercase -> Private variable
    name = "Muggle"

    New (name) {
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
