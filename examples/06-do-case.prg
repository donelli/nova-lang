
nNumber = 1

do case
   case nNumber = 1
      ? "one"
   case nNumber = 2
      ? "two"
   otherwise
      
      do case
         case nNumber > 2
            ? "grater than 2"
         case nNumber <= 1
            ? "less than 1"
      endcase
      
endcase

return

