package wiktparse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTemplate(t *testing.T) {
	refs := parseTemplate("{{inh|enm|ang|hello}}")
	require.Equal(t, 1, len(refs))
	require.Equal(t, refTypeInherited, refs[0].refType)
	require.Equal(t, langWord{
		language: "ang",
		word:     "hello",
	}, refs[0].word)
}

func TestParseTemplateDerived(t *testing.T) {
	refs := parseTemplate("{{der|eng|deu|qqqqq}}")

	require.Equal(t, 1, len(refs))
	require.Equal(t, refTypeDerived, refs[0].refType)
	require.Equal(t, langWord{
		language: "deu",
		word:     "qqqqq",
	}, refs[0].word)
}

func TestParseTemplateDerivedLongForm(t *testing.T) {
	refs := parseTemplate("{{derived|eng|deu|qqqqq}}")
	require.Equal(t, 1, len(refs))
	require.Equal(t, refTypeDerived, refs[0].refType)
	require.Equal(t, langWord{
		language: "deu",
		word:     "qqqqq",
	}, refs[0].word)
}

func TestParseTemplateMultiline(t *testing.T) {
	refs := parseTemplate(
		`{{derived
|eng
|deu
|qqqqq
}}`)
	require.Equal(t, 1, len(refs))
	require.Equal(t, refTypeDerived, refs[0].refType)
	require.Equal(t, langWord{
		language: "deu",
		word:     "qqqqq",
	}, refs[0].word)
}

func TestParseTemplatePrefix(t *testing.T) {
	refs := parseTemplate("{{prefix|en|re|lay}}")

	require.Equal(t, 2, len(refs))

	require.Equal(t, refTypePrefix, refs[0].refType)
	require.Equal(t, refTypeComponent, refs[1].refType)

	require.Equal(t, langWord{
		language: "en",
		word:     "re-",
	}, refs[0].word)

	require.Equal(t, langWord{
		language: "en",
		word:     "lay",
	}, refs[1].word)
}

func TestParseTemplateSuffix(t *testing.T) {
	refs := parseTemplate("{{suffix|en|deflate|ion}}")

	require.Equal(t, reference{
		refType: refTypeComponent,
		word: langWord{
			language: "en",
			word:     "deflate",
		},
	}, refs[0])

	require.Equal(t, reference{
		refType: refTypeSuffix,
		word: langWord{
			language: "en",
			word:     "-ion",
		},
	}, refs[1])
}

func TestParseTemplateConfix1(t *testing.T) {
	refs := parseTemplate("{{confix|en|neuro|genic}}")

	require.Equal(t, reference{
		refType: refTypePrefix,
		word: langWord{
			language: "en",
			word:     "neuro-",
		},
	}, refs[0])

	require.Equal(t, reference{
		refType: refTypeSuffix,
		word: langWord{
			language: "en",
			word:     "-genic",
		},
	}, refs[1])
}

func TestParseTemplateConfix2(t *testing.T) {
	refs := parseTemplate("{{confix|en|be|dew|ed}}")

	require.Equal(t, reference{
		refType: refTypePrefix,
		word: langWord{
			language: "en",
			word:     "be-",
		},
	}, refs[0])

	require.Equal(t, reference{
		refType: refTypeComponent,
		word: langWord{
			language: "en",
			word:     "dew",
		},
	}, refs[1])

	require.Equal(t, reference{
		refType: refTypeSuffix,
		word: langWord{
			language: "en",
			word:     "-ed",
		},
	}, refs[2])
}

func TestParseTemplateMention(t *testing.T) {
	refs := parseTemplate("{{m|en|hello}}")

	require.Equal(t, reference{
		refType: refTypeDerived,
		word: langWord{
			language: "en",
			word:     "hello",
		},
	}, refs[0])
}

func TestParseTemplateMentionBlank(t *testing.T) {
	refs := parseTemplate("{{m|en}}")

	require.Nil(t, refs)
}

func TestParseSimpleSection(t *testing.T) {
	fmt.Println(parseEtymologySection("{{prefix|en|a|biological|t1=without|t2=relating to life}}"))
}

func TestParseSection(t *testing.T) {
	fmt.Println(parseEtymologySection("From {{derived|eng|deu|qqqqq}}, abcd..."))
}

func TestParseSection2(t *testing.T) {
	fmt.Println(parseEtymologySection("First attested in 1664. From {{etyl|la|en}} {{m|la|circuitōsus}}"))
}

func TestParseSectionParens(t *testing.T) {
	fmt.Println(parseEtymologySection("From {{etyl|frm|en}} {{m|frm|democratie}} (French {{m|fr|démocratie}})"))
}

func TestParseSection3(t *testing.T) {
	fmt.Println(parseEtymologySection("{{suffix|en|deflate|ion}}"))
}

func TestParseSection4(t *testing.T) {
	// We want to ignore bits after the first valid "referencing clause", so we want to ignore " from" on
	// wards here
	fmt.Println(parseEtymologySection("Via {{etyl|de|cs}} {{m|de|Symptom}}<ref>{{R:Rejzek 2007}}</ref> from {{der|cs|grc|σύμπτωμα||a happening, accident, symptom of disease}}, from stem of {{m|grc|συμπίπτω||Ι befall}}, from {{m|grc|συν-||together}} + {{m|grc|πίπτω||I fall}}. "))
}

func TestParseSection5(t *testing.T) {
	fmt.Println(parseEtymologySection(`{{root|en|ine-pro|*ǵenh₁-}}
From {{inh|en|enm|nature}}, {{m|enm|natur}}, borrowed from {{bor|en|fro|nature}}, from 
`))
}

func TestParseSection6(t *testing.T) {
	fmt.Println(parseEtymologySection(`From {{suffix|en|πάθος|lang1=grc|t1=suffering|y}}`))
}

func TestParseSection7(t *testing.T) {
	fmt.Println(parseEtymologySection(`{{root|en|ine-pro|*h₁eḱ-}}
	From {{bor|en|la|-}} and {{der|en|NL.|-}} {{m|la|hippopotamus}}, from {{der|en|grc|ἱπποπόταμος}}, from {{m|grc|ἵππος||horse}} (English {{m|en|hippo-}}) + {{m|grc|ποταμός||river}}.`))
}

func TestParseSection8(t *testing.T) {
	fmt.Println(parseEtymologySection(`From {{af|grc|ῐ̔́ππος|ποτᾰμός|t1=horse|t2=river}}.`))
}
