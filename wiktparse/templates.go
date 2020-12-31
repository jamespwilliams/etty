package wiktparse

var templateParsers = map[string]func([]string) []reference{
	"af":        parseCompound,
	"affix":     parseCompound,
	"com":       parseCompound,
	"compound":  parseCompound,
	"pre":       parsePrefix,
	"prefix":    parsePrefix,
	"suf":       parseSuffix,
	"suffix":    parseSuffix,
	"con":       parseConfix,
	"confix":    parseConfix,
	"inh":       parseInherited,
	"inherited": parseInherited,
	"bor":       parseBorrowed,
	"borrowed":  parseBorrowed,
	"der":       parseDerived,
	"derived":   parseDerived,
	"m":         parseMention,
	"mention":   parseMention,
}

func parseCompound(components []string) []reference {
	if len(components) < 3 {
		return nil
	}

	var refs []reference
	for _, comp := range components[2:] {
		refs = append(refs, reference{
			refType: refTypeComponent,
			word: langWord{
				language: components[1],
				word:     comp,
			},
		})
	}

	return refs
}

func parsePrefix(components []string) []reference {
	if len(components) < 4 {
		return nil
	}

	pre := components[2]
	if len(pre) == 0 {
		return nil
	}

	if pre[len(pre)-1] != '-' {
		pre = pre + "-"
	}

	return []reference{
		{
			refType: refTypePrefix,
			word: langWord{
				language: components[1],
				word:     pre,
			},
		},
		{
			refType: refTypeComponent,
			word: langWord{
				language: components[1],
				word:     components[3],
			},
		},
	}
}

func parseSuffix(components []string) []reference {
	if len(components) < 4 {
		return nil
	}

	suf := components[3]
	if len(suf) == 0 {
		return nil
	}

	if suf[0] != '-' {
		suf = "-" + suf
	}

	return []reference{
		{
			refType: refTypeComponent,
			word: langWord{
				language: components[1],
				word:     components[2],
			},
		},
		{
			refType: refTypeSuffix,
			word: langWord{
				language: components[1],
				word:     suf,
			},
		},
	}
}

func parseConfix(components []string) []reference {
	if len(components) < 4 {
		return nil
	}

	pre := components[2]
	if len(pre) == 0 {
		return nil
	}

	if pre[len(pre)-1] != '-' {
		pre = pre + "-"
	}

	references := []reference{{
		refType: refTypePrefix,
		word: langWord{
			language: components[1],
			word:     pre,
		},
	}}

	if len(components) > 4 {
		references = append(references, reference{
			refType: refTypeComponent,
			word: langWord{
				language: components[1],
				word:     components[3],
			},
		})
	}

	suf := components[len(components)-1]
	if len(suf) == 0 {
		return nil
	}

	if suf[0] != '-' {
		suf = "-" + suf
	}

	references = append(references, reference{
		refType: refTypeSuffix,
		word: langWord{
			language: components[1],
			word:     suf,
		},
	})

	return references
}

func parseMention(components []string) []reference {
	if len(components) < 3 {
		return nil
	}

	word := components[2]
	if word == "" && len(components) >= 4 {
		word = components[3]
	}

	if word == "" {
		return nil
	}

	return []reference{
		{
			refType: refTypeDerived,
			word: langWord{
				language: components[1],
				word:     components[2],
			},
		},
	}
}

func parseSimpleTemplate(components []string, refType referenceType) []reference {
	if len(components) < 4 {
		return nil
	}

	return []reference{{
		refType: refType,
		word: langWord{
			language: components[2],
			word:     components[3],
		},
	}}
}

func parseInherited(components []string) []reference {
	return parseSimpleTemplate(components, refTypeInherited)
}

func parseBorrowed(components []string) []reference {
	return parseSimpleTemplate(components, refTypeBorrowed)
}

func parseDerived(components []string) []reference {
	return parseSimpleTemplate(components, refTypeDerived)
}
