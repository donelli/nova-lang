
* All clear commands

clear
clear all
clear fcache
clear gets
clear iostats
clear keys
clear locks
clear memory
clear menus
clear popups
clear program
clear prompt
clear screen
clear typeahead
clear window


* All close commands

close 
close 1
close tmpfile
close all 
close databases 
close format 
close index 
close procedure
close alternate
close alternate to print


* Dialog Boxes

DIALOG SCOPE

DIALOG QUERY
DIALOG QUERY LOCK

DIALOG MESSAGE "Hello" + " world!"

DIALOG BOX "Hello"
DIALOG BOX "Hello" LABEL "label"

* Not implemented

* DIALOG FIELDS 
* DIALOG FIELDS LABEL "label"

* DIALOG FILES LIKE *.prg
* DIALOG FILES LIKE *.prg TRIM
* DIALOG FILES LIKE *.prg LABEL "aa"

* DIALOG GET cVar pict "999" ;
*          title "Number" ;
*          label "Number"


* Misc commands

compile *.prg
compile /path/to/file/main.prg
compile ./prog???.prg

ALIAS show "dialog box"

eject

sleep 5

nNum = 7
sleep nNum * 2

* STORE

store 1 to nVar, nVar2, nVar3
store "hello world" to cHello

* RELEASE

release nVar
release all
release all like nVar*
release all except nVar*

return to master
return 420
return to master 420
return

