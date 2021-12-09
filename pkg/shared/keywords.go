package shared

const KeywordCount = 642

var KeywordsMap = make(map[string]string, KeywordCount)

func generateKeyword(keyword string, repeat bool) {
	if repeat && len(keyword) > 4 {
		for i := 4; i <= len(keyword); i++ {
			KeywordsMap[keyword[:i]] = keyword
		}
	} else {
		KeywordsMap[keyword] = keyword
	}
}

func LoadKeywords() {

	if len(KeywordsMap) > 0 {
		return
	}

	generateKeyword("additive", true)
	generateKeyword("alias", true)
	generateKeyword("alternate", true)
	generateKeyword("and", false)
	generateKeyword("append", true)
	generateKeyword("autojoin", false)
	generateKeyword("autorepair", false)
	generateKeyword("autosave", false)
	generateKeyword("bell", false)
	generateKeyword("blank", false)
	generateKeyword("blink", false)
	generateKeyword("blocksize", true)
	generateKeyword("border", false)
	generateKeyword("bottom", false)
	generateKeyword("box", false)
	generateKeyword("browse", true)
	generateKeyword("cacheload", false)
	generateKeyword("capture", false)
	generateKeyword("carry", false)
	generateKeyword("case", false)
	generateKeyword("catalog", false)
	generateKeyword("century", false)
	generateKeyword("clear", false)
	generateKeyword("clipper", true)
	generateKeyword("clock", false)
	generateKeyword("clockrate", false)
	generateKeyword("close", false)
	generateKeyword("color", false)
	generateKeyword("command", false)
	generateKeyword("commandwindow", false)
	generateKeyword("compatible", false)
	generateKeyword("compile", false)
	generateKeyword("compress", false)
	generateKeyword("confirm", true)
	generateKeyword("console", true)
	generateKeyword("copy", false)
	generateKeyword("count", false)
	generateKeyword("currency", true)
	generateKeyword("cursor", true)
	generateKeyword("databases", true)
	generateKeyword("date", false)
	generateKeyword("dbtrap", true)
	generateKeyword("dcache", true)
	generateKeyword("debug", false)
	generateKeyword("decimals", false)
	generateKeyword("declare", false)
	generateKeyword("default", true)
	generateKeyword("deleted", true)
	generateKeyword("delimiters", true)
	generateKeyword("descriptions", true)
	generateKeyword("design", true)
	generateKeyword("development", true)
	generateKeyword("device", true)
	generateKeyword("dialog", true)
	generateKeyword("dictionary", true)
	generateKeyword("directory", true)
	generateKeyword("display", true)
	generateKeyword("doescape", true)
	generateKeyword("dohistory", true)
	generateKeyword("echo", false)
	generateKeyword("edit", false)
	generateKeyword("editfield", true)
	generateKeyword("eject", false)
	generateKeyword("else", false)
	generateKeyword("elseif", false)
	generateKeyword("emacros", true)
	generateKeyword("encryption", true)
	generateKeyword("enddo", false)
	generateKeyword("endif", false)
	generateKeyword("erase", false)
	generateKeyword("error", false)
	generateKeyword("errorwindow", false)
	generateKeyword("escape", true)
	generateKeyword("exact", false)
	generateKeyword("except", false)
	generateKeyword("exclusive", true)
	generateKeyword("exit", false)
	generateKeyword("fastindex", true)
	generateKeyword("fcache", true)
	generateKeyword("field", false)
	generateKeyword("fields", false)
	generateKeyword("fieldval", false)
	generateKeyword("filecase", false)
	generateKeyword("filetype", false)
	generateKeyword("filter", false)
	generateKeyword("fixed", false)
	generateKeyword("fklabel", false)
	generateKeyword("flush", false)
	generateKeyword("for", false)
	generateKeyword("format", false)
	generateKeyword("formstate", false)
	generateKeyword("formupdate", false)
	generateKeyword("from", false)
	generateKeyword("fullpath", true)
	generateKeyword("function", true)
	generateKeyword("gateway", true)
	generateKeyword("gcache", true)
	generateKeyword("get", false)
	generateKeyword("go", false)
	generateKeyword("heading", true)
	generateKeyword("help", false)
	generateKeyword("helpfile", false)
	generateKeyword("helpwindow", false)
	generateKeyword("hiddenfield", true)
	generateKeyword("highlight", true)
	generateKeyword("history", true)
	generateKeyword("hours", false)
	generateKeyword("iblock", false)
	generateKeyword("icache", false)
	generateKeyword("if", false)
	generateKeyword("in", false)
	generateKeyword("index", false)
	generateKeyword("indexext", false)
	generateKeyword("inkeydelay", false)
	generateKeyword("instruct", false)
	generateKeyword("intensity", false)
	generateKeyword("journal", false)
	generateKeyword("kbedit", false)
	generateKeyword("key", false)
	generateKeyword("keys", false)
	generateKeyword("keyboard", true)
	generateKeyword("language", true)
	generateKeyword("ldcheck", true)
	generateKeyword("local", true)
	generateKeyword("locate", false)
	generateKeyword("lock", false)
	generateKeyword("lockwait", false)
	generateKeyword("loop", false)
	generateKeyword("mackey", false)
	generateKeyword("macros", true)
	generateKeyword("mail", false)
	generateKeyword("margin", true)
	generateKeyword("mark", false)
	generateKeyword("maxdbo", false)
	generateKeyword("mblock", false)
	generateKeyword("mconfirm", false)
	generateKeyword("memoclear", false)
	generateKeyword("memoext", false)
	generateKeyword("memoformat", false)
	generateKeyword("memosoftcr", false)
	generateKeyword("memowidth", false)
	generateKeyword("memowindow", false)
	generateKeyword("menu", false)
	generateKeyword("menubar", true)
	generateKeyword("message", true)
	generateKeyword("mouse", true)
	generateKeyword("multiuser", true)
	generateKeyword("navigate", true)
	generateKeyword("near", false)
	generateKeyword("next", false)
	generateKeyword("odometer", false)
	generateKeyword("off", false)
	generateKeyword("on", false)
	generateKeyword("optimize", true)
	generateKeyword("or", false)
	generateKeyword("order", true)
	generateKeyword("otherwise", false)
	generateKeyword("pagelength", false)
	generateKeyword("pagewidth", false)
	generateKeyword("parameters", true)
	generateKeyword("path", false)
	generateKeyword("pause", false)
	generateKeyword("pcache", false)
	generateKeyword("pcedit", false)
	generateKeyword("pcexact", false)
	generateKeyword("pcfilter", false)
	generateKeyword("pcgraphics", false)
	generateKeyword("pckeys", false)
	generateKeyword("pclocking", false)
	generateKeyword("pcpicture", false)
	generateKeyword("pcsays", false)
	generateKeyword("pcunique", false)
	generateKeyword("perfdial", false)
	generateKeyword("perfmeter", false)
	generateKeyword("picture", true)
	generateKeyword("point", true)
	generateKeyword("postfield", true)
	generateKeyword("postform", true)
	generateKeyword("postmenu", false)
	generateKeyword("postrecord", false)
	generateKeyword("precision", true)
	generateKeyword("prefield", true)
	generateKeyword("preform", false)
	generateKeyword("premenu", true)
	generateKeyword("prerecord", true)
	generateKeyword("print", false)
	generateKeyword("printer", false)
	generateKeyword("private", true)
	generateKeyword("procedure", true)
	generateKeyword("prompt", true)
	generateKeyword("pshare", false)
	generateKeyword("public", true)
	generateKeyword("querymode", true)
	generateKeyword("read", false)
	generateKeyword("readexit", false)
	generateKeyword("readinsert", false)
	generateKeyword("recordview", true)
	generateKeyword("refresh", true)
	generateKeyword("relation", true)
	generateKeyword("release", true)
	generateKeyword("replace", true)
	generateKeyword("reprocess", true)
	generateKeyword("restore", true)
	generateKeyword("retainmenu", true)
	generateKeyword("return", true)
	generateKeyword("rollback", true)
	generateKeyword("run", false)
	generateKeyword("runclear", true)
	generateKeyword("runwait", true)
	generateKeyword("safety", true)
	generateKeyword("save", false)
	generateKeyword("say", false)
	generateKeyword("schedule", true)
	generateKeyword("scoreboard", true)
	generateKeyword("screen", false)
	generateKeyword("screenio", false)
	generateKeyword("screenmap", false)
	generateKeyword("scroll", true)
	generateKeyword("seek", false)
	generateKeyword("select", true)
	generateKeyword("separator", true)
	generateKeyword("seqno", true)
	generateKeyword("set", false)
	generateKeyword("shadow", true)
	generateKeyword("skip", false)
	generateKeyword("sleep", false)
	generateKeyword("softseek", true)
	generateKeyword("space", true)
	generateKeyword("sql", false)
	generateKeyword("sqlprompt", false)
	generateKeyword("status", true)
	generateKeyword("step", true)
	generateKeyword("store", false)
	generateKeyword("stringbuf", true)
	generateKeyword("sum", false)
	generateKeyword("sysmenu", true)
	generateKeyword("talk", false)
	generateKeyword("tbufsize", true)
	generateKeyword("tedit", true)
	generateKeyword("terminal", true)
	generateKeyword("textmerge", true)
	generateKeyword("time", false)
	generateKeyword("timeout", true)
	generateKeyword("title", true)
	generateKeyword("to", false)
	generateKeyword("top", false)
	generateKeyword("total", false)
	generateKeyword("tracewindow", true)
	generateKeyword("trap", true)
	generateKeyword("typeahead", true)
	generateKeyword("underline", true)
	generateKeyword("unique", true)
	generateKeyword("unlock", true)
	generateKeyword("update", true)
	generateKeyword("use", false)
	generateKeyword("valid", false)
	generateKeyword("validate", false)
	generateKeyword("vaxtime", true)
	generateKeyword("verify", true)
	generateKeyword("view", false)
	generateKeyword("when", false)
	generateKeyword("while", false)
	generateKeyword("window", true)
	generateKeyword("with", false)
	generateKeyword("wp", true)
	generateKeyword("zap", false)

}
