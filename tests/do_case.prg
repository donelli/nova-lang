
do case
   case .t.
      ? "true"
   case .t.
      ? "true"
endcase
   
do case
   case .f.
      ? "false"
endcase

fPrintNumber(7)
fPrintNumber(751)
fPrintNumber(500)
fPrintNumber(-15)

return


**********************
function fPrintNumber
**********************
parameters pnNumber

do case
   case pnNumber < 0
      ? "negative"
   case pnNumber < 10
      ? "less then 10"
   case pnNumber < 100
      ? "less then 100"
   otherwise
   
      do case
         case pnNumber % 2 = 0
            ? "more or equals then 100 and even"
         otherwise
            ? "more or equals then 100 and odd"
      endcase
   
endcase

return

