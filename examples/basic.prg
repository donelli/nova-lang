
set procedure to
set procedure to  /path/to/library
set procedure to  otherlib additive
set procedure to  /path/to/otherlib2 additive

cName      = space(10)
dBirthDate = {01/01/21}

@ 01, 01 say "Name......:"
@ 02, 01 say "Birth Date:"
@ 03, 01 say "Idade.....:"

@ 01, 14 get cName pict "@!" valid(!empty(cName))
@ 02, 14 get dBirthDate valid(dBirthDate <> {})

read

nIdade = int((date() - dBirthDate) / 365)

@ 03, 14 say reverse(nIdade)
@ 05, 01 say "Nicknames.:"

declare vcNicknames[4]
afill(vcNicknames, space(15))

fReadNicknames("vcNicknames")

return


************************
function fReadNicknames
************************
parameters pcArray

nRow   = 5
nCol   = 14
nIndex = 1

do while .t.

   cNickname = &pcArray[nIndex]
   
   @ nRow, nCol get cNickname pict "@!"
   read
   
   if empty(cNickname)
      exit
   endif
   
   &pcArray[nIndex] = cNickname
   ++nIndex
   
   if nIndex > 4
      exit
   endif
   
enddo

return

