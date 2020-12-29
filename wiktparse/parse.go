package wiktparse

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/jamespwilliams/etymology/wiktlang"
)

type textElem struct {
	Data string `xml:",chardata"`
}

var titleRemovalRegex = regexp.MustCompile(`Reconstruction:[^:]*/`)

func ParseDump(dump io.Reader, out io.Writer, languages wiktlang.Languages) error {
	var currentTitle string

	d := xml.NewDecoder(dump)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding token: %w", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			switch ty.Name.Local {
			case "title":
				var title textElem
				if err = d.DecodeElement(&title, &ty); err != nil {
					return fmt.Errorf("error decoding title: %w", err)
				}
				currentTitle = title.Data
				currentTitle = titleRemovalRegex.ReplaceAllString(currentTitle, "")
			case "text":
				var text textElem
				if err = d.DecodeElement(&text, &ty); err != nil {
					return fmt.Errorf("error decoding text: %w", err)
				}

				parseTextTag(out, currentTitle, text.Data, languages)
			}
		}
	}

	return nil
}

func parseTextTag(out io.Writer, title, text string, languages wiktlang.Languages) {
	var (
		currentLanguage    string
		inEtymologySection bool
	)
	sections := make(map[string]string)

	for _, line := range strings.Split(text, "\n") {
		if len(line) < 4 {
			continue
		}

		if line[0] == '=' && line[1] == '=' && line[2] != '=' {
			currentLanguage = line[2 : len(line)-2]
			continue
		}

		if line == "===Etymology===" {
			inEtymologySection = true
			continue
		}

		if line[0] == '=' {
			inEtymologySection = false
			// TODO: call parseEtymologySection here, to avoid needing to store it in a map
			continue
		}

		if inEtymologySection {
			sections[languages.CodeFromName(currentLanguage)] += line + "\n"
		}
	}

	for lang, section := range sections {
		refs := parseEtymologySection(section)
		for _, ref := range unique(refs) {
			fmt.Fprintf(out, "%v:%v\trel:%s\t%v:%v\n", lang, title, ref.refType, ref.word.language, ref.word.word)
		}
	}
}

var (
	withinParenRegex  = regexp.MustCompile(`\([^)]*\)`)
	rootTemplateRegex = regexp.MustCompile(`\{\{root[^}]*\}\}`)
	imgTemplateRegex  = regexp.MustCompile(`(?s)\{\{multiple[^}]*\}\}`)
	templatesRegex    = regexp.MustCompile(`\{\{[^}]*\}\}`)
)

func parseEtymologySection(section string) []reference {
	// Walk the section to try and find the meaningful bit:
	var (
		templateCount     int
		startIndex        int
		endIndex          int = len(section) - 1
		containedTemplate bool
	)

	// TODO: we also want to ignore empty templates (e.g: {{inh|sco|enm|-}} should _not_ be treated as a template)

	// maybe it'd be useful to extract a list of templates and their punctuation. The word "From" also seems important.
	// the first valid template seems very likely to be the correct one. the only complication is when there's a "+" or
	// an "or" separating two (or more) templates
	//
	// I think if we encounter a comma after a valid template, we can just stop there.

	section = rootTemplateRegex.ReplaceAllString(section, "")
	section = imgTemplateRegex.ReplaceAllString(section, "")

outer:
	for i, c := range section {
		switch c {
		case '{':
			templateCount++
			containedTemplate = true
		case '}':
			templateCount--
		case ',', '.', '\n':
			// The useful bit is often before the first comma. But there's also sometimes a prelude, which also usually
			// ends in a comma, dot or newline. But the prelude doesn't usually contain a template.
			if containedTemplate && templateCount == 0 {
				endIndex = i - 1
				break outer
			}

			if templateCount == 0 {
				// if we aren't in a template right now, then reset the start index
				startIndex = i
			}
		}
	}

	if !containedTemplate {
		return nil
	}

	// TODO: Remove the bits in parentheses:
	if endIndex > len(section) {
		println("endIndex > len(section), hm... section=" + section)
		return nil
	}

	section = section[startIndex : endIndex+1]
	section = withinParenRegex.ReplaceAllString(section, "")

	var refs []reference
	for _, template := range templatesRegex.FindAllString(section, -1) {
		refs = append(refs, parseTemplate(template)...)
	}

	return refs
}

// parseTemplate takes a template and returns a list of references which that template makes
func parseTemplate(template string) []reference {
	template = strings.Trim(template, "{}")

	dirtyComponents := strings.Split(template, "|")

	var components []string
	for _, comp := range dirtyComponents {
		if strings.Contains(comp, "=") {
			// Filter out named parameters, which we aren't interested in (yet?)
			continue
		}

		components = append(components, strings.TrimSpace(comp))
	}

	var refType referenceType
	switch components[0] {
	case "pre", "prefix":
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
	case "suf", "suffix":
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
	case "con", "confix":
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
	case "inh", "inherited":
		refType = refTypeInherited
	case "bor", "borrowed":
		refType = refTypeBorrowed
	case "der", "derived":
		refType = refTypeDerived
	case "m", "mention":
		if len(components) < 3 {
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
	default:
		return nil
	}

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

func unique(refs []reference) []reference {
	for i := 0; i < len(refs); i++ {
		for i2 := i + 1; i2 < len(refs); i2++ {
			if refs[i] == refs[i2] {
				// delete
				refs = append(refs[:i2], refs[i2+1:]...)
				i2--
			}
		}
	}
	return refs
}
