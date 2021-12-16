

f3Equals(fDouble(2), 2, 3)

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

