
* Trim string
? alltrim(" hello world!     ")

* Number to string
? str(154)
? str(154.77)
? str(24.7, 8, 2)
? str(7541, 12, 2)
? str(152.1159, 10, 3)

* Evaluate expression and return its type, or `U` if undefined or error
? type("150")
? type('"hello world!"')
? type('.t.')
? type('5 + 7')
? type('str(14)')
? type('someFunc(14)')
? type('dial box "wasd"')

* String to number
? val("154.1")
? val("9745142714571")
? val("wasd")
? val("9.14,1")
? val(".157946")

* Verify is value is empty
? empty('wasd')
? empty('')
? empty(150)
? empty(0)
? empty(.t.)
? empty(.f.)

* Generaters a empty string of specified length
? "." + space(10) + "."
? "." + space(100) + "."
? "." + space(0) + "."

return

