package wiktparse

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/jamespwilliams/etty/wiktlang"
	"go.uber.org/zap"
)

type textElem struct {
	Data string `xml:",chardata"`
}

var titleRemovalRegex = regexp.MustCompile(`Reconstruction:[^:]*/`)

func ParseDump(log *zap.Logger, dump io.Reader, out io.Writer, languages wiktlang.Languages) error {
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

				parseTextTag(log, out, currentTitle, text.Data, languages)
			}
		}
	}

	return nil
}

func parseTextTag(log *zap.Logger, out io.Writer, title, text string, languages wiktlang.Languages) {
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
			languageCode, ok := languages.CodeFromName(currentLanguage)
			if !ok {
				log.Debug("couldn't find code for language", zap.String("language", currentLanguage))
				continue
			}

			sections[languageCode] += line + "\n"
		}
	}

	for lang, section := range sections {
		refs := parseEtymologySection(log, section)
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

func parseEtymologySection(log *zap.Logger, section string) []reference {
	// Walk the section to try and find the meaningful bit:
	var (
		templateCount     int
		startIndex        int
		endIndex          int = len(section) - 1
		containedTemplate bool
	)

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

	if endIndex > len(section) {
		log.Debug("endIndex > len(section) in parseEtymologySection, returning nil")
		return nil
	}

	section = section[startIndex : endIndex+1]
	section = withinParenRegex.ReplaceAllString(section, "")

	var refs []reference
	for _, template := range templatesRegex.FindAllString(section, -1) {
		refs = append(refs, parseTemplate(log, template)...)
	}

	return refs
}

// parseTemplate takes a template and returns a list of references which that template makes
func parseTemplate(log *zap.Logger, template string) (refs []reference) {
	defer func() {
		// TODO: clean this up:
		var res []reference
		for _, ref := range refs {
			if ref.word.word != "-" {
				res = append(res, ref)
			}
		}

		refs = res
	}()

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

	if len(components) == 0 {
		return nil
	}

	parser, ok := templateParsers[components[0]]
	if !ok {
		return nil
	}

	return parser(components)
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
