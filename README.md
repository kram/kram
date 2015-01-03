# Variables
var a = 123 // new Number(123)
var b = "fniss" // new String("fniss")

# Defined variables
null // new Null
false // new Boolean(false)
true // new Boolean(false)

IO.Print("Yayaya")

# Custom classes
class MyClass {

    # Public method
    Public {
        return 123
    }

    # Private method
    private {
        return 123
    }
}

# Using a class
var me = new MyClass
me.Method()

# Functions
var func = new Fn {
    IO.Print("I am a function")
}

## Functions are basically this
class Fn {
    new (function) {
        this.Call = function
    }
}

## You can mimmick a function like this
class CustomFunction {
    Call {
        IO.Print("I am a function")
    }
}

## These are the samme
func()
func.Call()

# Looping

var items = ["one", "two", "three", "four"]

for item in items {
    IO.Print(item)
}

IO.Print(items[0])

items[] = "five"

## Different types
[123, "four", true]