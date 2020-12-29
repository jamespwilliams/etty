package main

// TODO: perhaps words/fragments should be represented as being in one of these different categories:
// - "compound" - composed of prefix/base/suffix - not actually derived/inherited/borrowed, rather
//   a combination of other sub-words in the same language
// - "borrowed"
// - "derived"
// - "inherited" (each of these taken from Wiktionary's definition)
//
// How should a word be marshalled in the output of this program? options:
// - lists of edges, e.g: a number of "word rel:xyz other_word" lines
//   easy, conforms to existing wordnet format
// - map from a word to an (ordered) list of each relation for that word
// - map from a word to {
//      wordType (as defined above),
//		[]reference/relation
//   }
//

type referenceType int

const (
	refTypeInherited referenceType = iota
	refTypeBorrowed
	refTypeDerived
	refTypeComponent
	refTypePrefix
	refTypeSuffix
)

var referenceTypes = map[referenceType]string{
	refTypeInherited: "inherited",
	refTypeBorrowed:  "borrowed",
	refTypeDerived:   "derived",
	refTypeComponent: "component",
	refTypePrefix:    "prefix",
	refTypeSuffix:    "suffix",
}

type reference struct {
	// references also have a source language, might be useful
	refType referenceType
	word    langWord
}

type langWord struct {
	language string
	word     string
}

func (t referenceType) String() string {
	return referenceTypes[t]
}
