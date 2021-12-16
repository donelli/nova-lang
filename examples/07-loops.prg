
nNum = 0

do while .t. 
   
   nNum = nNum + 1
   
   if nNum > 10
      exit
   endif
   
   if nNum % 2 = 0
      loop
   endif
   
   ? "Here"
   
enddo

return

