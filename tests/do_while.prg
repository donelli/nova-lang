
nNumber = 1

do while nNumber <= 10
   ? nNumber
   nNumber = nNumber + 1
enddo

? "-----------------"

do while .t.
   
   if nNumber = 5
      exit
   endif
   
   ? nNumber

   nNumber = nNumber - 1
   
enddo

? "-----------------"

do while nNumber != 0
   
   if nNumber = 2
      nNumber = nNumber - 1
      loop
   endif
   
   
   if nNumber % 2 = 0
      ? nNumber
   endif
   
   nNumber = nNumber - 1
   
enddo

return

