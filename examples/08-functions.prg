
? f3Equals(fDouble(1), 2, 1 * 2)

fPublic()

? pub

return



function f3Equals
parameters pcParam, pcParam2, pcParam3

if pcParam = pcParam and pcParam2 = pcParam3
   return .t.
endif

return .f.


procedure fDouble
   parameters pnVal

return pnVal * fGetNumber()


function fGetNumber
return 2


function fPublic

public pub
pub = 123

private pub
pub = 420

? pub

return 