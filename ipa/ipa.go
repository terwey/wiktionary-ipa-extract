package ipa

import (
	"strings"

	"github.com/goccy/go-json"
)

type Pronunciation struct {
	Word string        `json:"word"`
	IPA  []IPATemplate `json:"ipa"`
}

type IPATemplate struct {
	Language string `json:"lang,omitempty"`
	IPA      []string
	Variant  string `json:"variant,omitempty"`
}

func (p Pronunciation) JSON() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func ParseIPATemplate(text string) IPATemplate {
	if !strings.HasPrefix(text, "{{IPA|") {
		return IPATemplate{}
	}

	if strings.Contains(text, "disambiguation") {
		return IPATemplate{}
	}

	// wikipedia is an insane format
	if strings.Contains(text, "==") {
		return IPATemplate{}
	}

	text = strings.TrimPrefix(text, "{{")
	text = strings.TrimSuffix(text, "}}")

	tpl := IPATemplate{}
	split := strings.Split(text, "|")

	for i, s := range split {
		// IPA literal
		if i == 0 {
			if s != "IPA" {
				return tpl
			}
			continue
		}

		if i == 1 {
			// filter language we care for

			// if s != lang {
			// 	return tpl
			// }
			tpl.Language = s
			continue
		}

		if strings.Contains(s, "=") {
			tpl.Variant = s
			continue
		}

		if s == "/" {
			continue
		}

		if !(strings.HasPrefix(s, "/") && strings.HasSuffix(s, "/")) {
			continue
		}

		tpl.IPA = append(tpl.IPA, s)
	}

	return tpl
}

func FindIPA(text string, pronunciation Pronunciation) Pronunciation {
	if !strings.Contains(text, "{{IPA") {
		return pronunciation
	}

	// there can be multiple instaces of the {{IPA}} template
	// for now we want only US instances
	// example
	// |a=UK means it's UK pronunciation
	// * {{IPA|en|/ˈwɪʃ.iˌwɒʃ.i/|a=UK}}
	// |a=US means it's US pronunciation
	// * {{IPA|en|/ˈwɪʃ.iˌwɑ.ʃi/|/ˈwɪʃ.iˌwɔ.ʃi/|a=US}}
	// without the |a= it means it's unspecified
	// * {{IPA|is|/heiː/}}

	// log.Printf("text: \n---\n%s\n---\n", text)

	count := strings.Count(text, "{{IPA")

	// index of the first {{IPA
	index := strings.Index(text, "{{IPA")
	end := strings.Index(text[index:], "}}")
	// chopped up IPA template
	if end == -1 {
		if buf := ParseIPATemplate(text[index:]); len(buf.IPA) != 0 {
			pronunciation.IPA = append(pronunciation.IPA, ParseIPATemplate(text[index:]))
		}
		// return append(existing, ParseIPATemplate(text[index:]))
	}

	ipa := text[index : index+end+2]
	if buf := ParseIPATemplate(ipa); len(buf.IPA) != 0 {
		pronunciation.IPA = append(pronunciation.IPA, buf)
	}
	if count == 1 {
		return pronunciation
	}

	// log.Printf("found: %v, recursive: %s", pronunciation, text[index+end+2:])
	return FindIPA(text[index+end+2:], pronunciation)
}
