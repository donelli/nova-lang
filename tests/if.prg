
if .t.
   ? "true"
else
   ? "false"
endif

fPrintNumber(1)
fPrintNumber(88)
fPrintNumber(5712)
fPrintNumber(-150)

return


**********************
function fPrintNumber
**********************
parameters pnNumber

if pnNumber < 0
   ? "negative"
elseif pnNumber < 10
   ? "less than 10"
elseif pnNumber < 100
   ? "less than 100"
else
   
   if pnNumber = 100
      ? "equal to 100"
   else
      ? "greater than 100"
   endif
   
endif

return

