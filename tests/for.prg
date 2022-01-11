
* Basic for with step = 1
for lnaa = 1 to 5
   ? lnaa
next

? "-------------"

* Inverse for, don't print anything
for lnaa = 5 to 1
   ? lnaa
next

? "-------------"

* Inverse for, this should work
for lnaa = 5 to 1 step -1
   ? lnaa
next

? "-------------"

* The programmer can loop to the start of the loop
for lnaa = 1 to 10
   
   if lnaa % 2 = 0
      loop
   endif
   
   ? lnaa
   
next

? "-------------"

* The programmer can exit out of the loop 
for lnaa = 1 to 100
   
   if lnaa = 5
      exit
   endif
   
   ? lnaa
   
next

? "-------------"

* The counter can be changed during the execution of the loop
for lnaa = 2 to 10
   
   ? lnaa
   
   lnaa = lnaa + 1
   
next

return

