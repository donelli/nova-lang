
? f3Equals(fDouble(1), 2, 1 * 2)

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
public pub
pub = "teste"
return 2

