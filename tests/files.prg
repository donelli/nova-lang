
nFd = fopen("./files.prg", 1)

if nFd = -1
   ? "Error opening file"
   return
endif

cLine = fread(nFd)
? "|" + cLine + "|"

cLine = fread(nFd)
? "|" + cLine + "|"

* do while !feof(nFd)
   
*    cLine = fread(nFd)
   
*    ? cLine
   
* enddo

fclose(nFd)

return

