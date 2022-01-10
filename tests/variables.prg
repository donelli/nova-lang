
* Simple assigniment
nNumVar = 1
? nNumVar

* Assigniment with a more complex expression
nNumVar2 = (nNumVar + 1) * 2
? nNumVar2

* Reassignement of a existing variable
nNumVar = "Hello"
? nNumVar

* Store command
store "Hello world!" to cText, cText2
? cText + " " + cText2

* variables are declared as .f. if not assigned

private cPrivVar
public cPubVar

? type("cPrivVar")
? type("cPubVar")

cPrivVar = "Hello"
? type("cPrivVar")

? type("1 + 1")

return

