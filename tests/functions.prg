
? str(fDouble(12), 10, 2)

* Recursion is not working right now
* ? fRecursive(1)

return


*****************
procedure fDouble
*****************
parameters pnNumber

return pnNumber * 2


********************
function fRecursive
********************
parameters pnLevel

if pnLevel > 10
   return 1
endif

nVal = fRecursive(pnLevel + 1)

return nVal

