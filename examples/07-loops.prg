
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

for i = 0 to 5
   ? i
next

for i = 5 to 0 step -1
   ? i
next

for i = 0 to 5
   for j = 0 to 5
      ? i * j
   next
next

return

