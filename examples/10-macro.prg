
cMacro = ".t."

cMacroFun = "dial box 'macro'"
&cMacroFun

if !empty(cMacro) and &cMacro
   dial box "true"
else
   dial box "false"
endif

return

