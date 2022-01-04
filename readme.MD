
# Nova programming language

Nova is a interpreted programming language written in [Golang](https://go.dev/).

The main ideia of this language is to offer a open source alternative to paid programming languages based on [dBASE](https://pt.wikipedia.org/wiki/DBASE), like [FoxPro](https://en.wikipedia.org/wiki/Visual_FoxPro), [Recital](https://www.recitalsoftware.com/products/recital), [Clipper](https://pt.wikipedia.org/wiki/Clipper_(linguagem_de_programa%C3%A7%C3%A3o))...

After the basic implementation, I will include new commands and features that I think are important, like native SQL support...

## Compiling

```sh
go build -o nova.exe cmd/recital/nova.go
```

## Legend
```
❌ Not implemented and probably will not be
✔️ Parsing process OK
⭐ Interpreted
🔨 Planned / In progress
```

## Basics

```
⭐ Strings
⭐ Numbers
⭐ Boolean/logic
⭐ Variables
⭐ Logical operations (and, or, equals, ...)
⭐ Mathematic operations (plus, minus, ...)
Arrays
⭐ Functions
⭐ Base for built-in functions
✔️ Macros
```

## Ideas
```
Assert (Tests)
Sql
```

## Commands / Blocks

### If statement
```
⭐ IF <condition>
⭐      <statements>
⭐ [ELSEIF <condition>
⭐      <statements>]
⭐ [ELSE
⭐      <statements>]
⭐ ENDIF
```

### Do case statement (like switch/case)
```
✔️ DO CASE
✔️    CASE <condition>
✔️       <statements>
✔️    [OTHERWISE
✔️       <statements>]
✔️ ENDCASE
```

### Do while statement
```
✔️ DO WHILE <condition>
✔️    <statements>
✔️ ENDDO
```

### For statement
```
✔️ FOR <var> = <start_value> TO <end_value> [STEP <number>]
✔️    <statements>
✔️ NEXT
```

### Save commands
```
❌ SAVE COLOR TO <memvar>
❌ SAVE COLOR TO FILE <filename>
🔨 SAVE GETS [TO <memvar>]
❌ SAVE KEYS TO <memvar>
❌ SAVE MENU [TO <memvar>]
❌ SAVE RECORDVIEW TO <memvar>
🔨 SAVE SCREEN TO <memvar>
❌ SAVE SCREEN TO FILE <.img file>
❌ SAVE SCREEN TO FILE (<exp>)
🔨 SAVE SCREEN [AT <row>, <col> TO <endrow>, <endcol>]
❌ SAVE TO <.mem file> / (<expC>) [ERROR / ALL LIKE <skeleton> / ALL EXCEPT <skeleton>]    |  ex: save to monday all like mon_*
❌ SAVE WINDOW <window-name list>/ALL TO <.win filename>                                   | ex: save windows invoice,payables to acct
```

### Restore command
```
❌ RESTORE COLOR FROM <memvar>
❌ RESTORE COLOR FROM FILE <filename>
❌ RESTORE FROM <.mem file> [ADDITIVE]
🔨 RESTORE GETS [TO <memvar>]      (??? FROM ???)
❌ RESTORE KEYS FROM <memvar>
❌ RESTORE MENU [FROM <memvar>]
❌ RESTORE RECORDVIEW FROM <memvar>
🔨 RESTORE SCREEN [AT <row>,<col> TO <row>,<col>]
🔨 RESTORE SCREEN FROM [<memvar>] / [FILE <.img file>] [AT <row>,<col> TO <row>,<col>]
❌ RESTORE WINDOW <window-name list>/ALL FROM <.win filename>
```

### Append command
```
❌ APPEND [NOCLEAR]
🔨 APPEND BLANK [<expN>]
🔨 APPEND FROM <filename>/(<expC>)
🔨	[FOR <condition>] [WHILE <condition>]
🔨	[[TYPE] SDF/FIXED/2020/DELIMITED WITH BLANK / DELIMITED [WITH <delimiter>]
🔨 APPEND FROM ARRAY <array> [FOR <condition>] [WHILE <condition>] [REINDEX]
❌ APPEND MEMO <memo field> FROM <filename> [OVERWRITE]
```

### @ ... say command
```
🔨 @<row>,<col>
🔨 @<row>,<col>,<endrow>,<endcol> [BOX <expC>]
🔨 @<row>,<col> CLEAR
🔨 @<row>,<col> CLEAR TO <endrow>,<endcol> [REVERSE] [BOLD]
❌ @<expN1>,<expN2> FILL TO <expN3>,<expN4> [COLOR <color>/COLOR SCHEME <expN5>]
❌ @<row>,<col> MENU <expC1> [COMMAND <expC2>[;<expC>...]] [HELP <expC3>] [HELPFILE <.hlp filename>]
❌    [PULLDOWN <expC4>/@<procedure> WITH <parameters> [NOREFRESH]
🔨 @<row>,<col> PROMPT <expC1> [MESSAGE <expC2>]
🔨 @<row>,<col> SAY <exp1> [PICTURE <expC1>]
❌    [GET <variable>] [FUNCTION <expC2>] [COLOR <foreground/background>]
❌ @<row>,<col> TO <endrow>,<endcol> SCROLL
❌    [UP/DOWN/LEFT/RIGHT] [BY <expN>] [WRAP]
❌ @<row>,<col> TO <endrow>,<endcol>
❌    [SCROLL [UP/DOWN/LEFT/RIGHT] [BY <expN1>] [WRAP]] / [LABEL <expC>] / [BOX <expC>]
❌    [COLOR <expC>] [DOUBLE/SINGLE/PANEL] [CLEAR] [BOLD] [REVERSE] [SHADOW]
```

### @ ... get command
```
🔨 @<row>,<col> GET <memvar/field>
❌   [BUTTON <expC> [LABEL <expC>] [GROUP <expC>] [HELP <expC>] [TRIGGER <expC>]]
❌   [CALCULATE] [CALCULATED BY <expression>]
❌   [CHOICELIST <expC>] [ERROR <expC>]
❌   [FUNCTION <expC>]
🔨   [PICTURE <template> [NOECHO]
❌   [HELP <expC>] [LOOKUP IN <alias>]
❌   [MESSAGE <expC>] [MUST_ENTER]
🔨   [PREFIELD <procedure-name>] [POSTFIELD <procedure-name>]
🔨   [RANGE <minimum>, <maximum>]
🔨   [READONLY [IF <condition>]]
❌   [RELATION INTO <alias>]
🔨   [VALID <condition>]
❌   [VALID IN <expC>] [VALIDATE WITH <program/procedure name>]
🔨   [WHEN <condition>]
```

### Sequence command
```
🔨 BEGIN SEQUENCE
🔨 	<commands>...
🔨 	[BREAK]
🔨 END [SEQUENCE]
```

### Classes
```
🔨 CLASS <class name> [[OF <class name list>] DYNAMIC]
🔨    [PUBLIC [NOTIFY / PROPAGATE]: / LOCAL: / STATIC:]
🔨    PROPERTY <property name> [AS CHARACTER / NUMERIC / DATE / LOGICAL / <class name>]
🔨    [METHOD <method name> / EXTERNAL]
🔨 	...
🔨    RETURN]
🔨 ENDCLASS
```

### Copy command
```
❌ COPY DICTIONARY TO <.xad file> [ADDITIVE]
❌ COPY FILE <from file> TO <to file>
❌ COPY INDEXES <.ndx filelist> [TO <.dbx filename>]
❌ COPY MEMO <memo field> TO <filename> [ADDITIVE]
❌ COPY STRUCTURE TO <filename> [FIELDS <field>[,<field>]...] [[WITH] PRODUCTION]
❌ COPY STRUCTURE EXTENDED TO <file> [FIELDS <field list>]
❌ COPY TAG <tagname> [OF <.dbx filename>] TO <.ndx filename>
🔨 COPY TO <file>/(<expC>) [<scope>] [FIELDS <field list> [FOR <condition>] [WHILE <condition>] [[TYPE] SDF/2020/FIXED/DELIMITED WITH BLANK/DELIMITED [WITH <delimiter>]
❌ COPY TO ARRAY <array> [FIELDS <field list>] [WHILE <condition>] [FOR <condition>]
```

### Create command
```
🔨 CREATE <.dbf file>
❌ CREATE BRIDGE <.brg file>
❌ CREATE DICTIONARY FROM <.xad file>
🔨 CREATE <new file> FROM <structure extended file>
❌ CREATE GATEWAY <.gtw filename>/(<expC>)
❌ CREATE GATEWAY FROM <server>/<gateway catalog>
❌ CREATE LABEL <.lbl file>/(<expC>)
❌ CREATE REPORT <.frm file>/(<expC>)
❌ CREATE SCREEN <.scr file>/(<expC>)
❌ CREATE VIEW <.vue file>/(<expC>) [FROM ENVIRONMENT]
```

### Delete command
```
🔨 DELETE [<scope>] [FOR <condition>] [WHILE <condition>]
❌ DELETE FILE <filename>/(<expC>)
❌ DELETE TAG <tagname1> [OF <.dbx filename>] [, <tagname2> [OF <.dbx filename>]...]
```

### Dialog command
```
✔️ DIALOG BOX <expC1> [LABEL <expC2>]
✔️ DIALOG FIELDS [LABEL <expC>]
❌ DIALOG FILES LIKE <skeleton> [TRIM] [LABEL <expC>]
❌ DIALOG GET <memvar> [PICTURE <expC>] [HELP <expC>] [LABEL <expC>] [TITLE <expC>]
✔️ DIALOG MESSAGE <expC>
✔️ DIALOG QUERY [LOCK]
✔️ DIALOG SCOPE
```

### DBF commands
```
✔️ BROWSE
✔️   [FIELDS <field> [,<field>]...]
✔️   [FOR <condition>]
❌   [COMPRESS]
❌   [FREEZE <expC>]
❌   [FORMAT]
❌   [KEY <exp>]
❌   [LAST]
❌   [LOCK <expN>]
❌   [NOAPPEND]
❌   [NOCLEAR]
❌   [NODELETE]
❌   [NOEDIT]
❌   [NOFOLLOW]
❌   [NOINIT]
❌   [NOMENU]
❌   [NOMODIFY]
❌   [NORMAL]
❌   [NOORGANIZE]
❌   [NOWAIT]
❌   [OVERLAY]
❌   [WIDTH <expN>]
❌   [WINDOW <window name>]
❌   [COLOR SCHEME <expN> / COLOR <colors>]

✔️ COUNT [<scope>] [WHILE <condition>] [FOR <condition>] [TITLE <expC>] [TO <memvar>]   |   <scope>  ->  ALL/REST  (default: ALL)

🔨 FIND <key>

🔨 INDEX ON
🔨   <key expression> TO <.ndx file> /
❌   TAG <tagname> [OF <.dbx file>]
🔨   [FOR <condition>] [UNIQUE]
❌   [TITLE <expC>] [DESCENDING]

🔨 LOCATE [<scope>] [FOR <condition>] [WHILE <condition>]   |   <scope>  ->  ALL/REST  (default: ALL)

🔨 REPLACE [<scope>] <field1> WITH <exp1> [ADDITIVE]
🔨   [,<field2> WITH <exp2> [ADDITIVE]...]
🔨   [FOR <condition>] [WHILE <condition>]
🔨   [FROM ARRAY <array-name>]
❌   [REINDEX] [<memo-field> WITH <exp> ADDITIVE] [BLANK]
<!-- The default <scope> is the current record, unless
the FOR clause is present, in which case the scope is ALL
scope can be REST when used with WHILE -->

🔨 SEEK <key expression>

🔨 SKIP <expN> [IN <alias name | workarea number>]

🔨 SUM [<scope>] <expN> [,<expN>] [TO <memvar list>] [FOR <condition>] [WHILE <condition>]
<!-- If <scope> is not specified, then ALL is used. scope can be REST when used with WHILE -->

🔨 TOTAL ON <key field> TO <.dbf file> [<scope>] [FIELDS <field list>] [FOR <condition>] [WHILE <condition>]
❌   [SUMMARY]
<!-- If <scope> is not specified, then ALL is used. scope can be REST when used with WHILE -->

🔨 UNLOCK [ALL] [IN <alias>]

🔨 USE [<.dbf file> /
❌     <.brg file>/<.gtw file>]
🔨   [IN <workarea>/<alias>] [ALIAS <name>]
❌   [INDEX <index list>]
❌   [ORDER [TAG] <tag name> [OF <.dbx file>]] [ORDER BY <exp1>] [WHERE <exp2>]
🔨   [EXCLUSIVE]
❌   [NOUPDATE] [NODBX]

🔨 ZAP
```

### Others command
```
✔️ ALIAS <synonym> [<expC>]

✔️ COMPILE <skeleton>

❌ CONTINUE

✔️ DO
✔️   <.prg file> / 
✔️   <procedure name> /
❌   (<expC>)
✔️   [WITH <parameter>[,<parameter>]...]

✔️ EJECT

✔️ ERASE <filename>/<exprC>

🔨 FLUSH

🔨 KEYBOARD <expC> [CLEAR]

🔨 READ
❌  [SAVE] [NOUPDATE] [NOREFRESH] [NOCLEAR] [CYCLE] [OBJECT <expN>] [MODAL]

✔️ RELEASE <memvar list> / ALL [LIKE / EXCEPT <skeleton>]

✔️ RETURN [TO MASTER] [<exp>]

🔨 RUN <command line>
🔨 !   <command line>

✔️ SLEEP <expN>

✔️ STORE <exp> TO <memvar>
```
