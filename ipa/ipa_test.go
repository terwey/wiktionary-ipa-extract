package ipa

import (
	"reflect"
	"testing"
)

type parseIpaTemplateTestData struct {
	name string
	text string
	want IPATemplate
}

func getParseIPATemplateTests() []parseIpaTemplateTestData {
	tests := []parseIpaTemplateTestData{
		{
			text: "{{IPA|en|/ˈkɛstɹəl/}}",
			want: IPATemplate{
				Language: "en",
				IPA:      []string{"/ˈkɛstɹəl/"},
			},
		},

		{
			text: "{{IPA|en|/ˈspæɹoʊ/|/ˈspɛɹoʊ/|a=US}}",
			want: IPATemplate{
				Language: "en",
				IPA: []string{
					"/ˈspæɹoʊ/",
					"/ˈspɛɹoʊ/"},
				Variant: "a=US",
			},
		},
	}
	return tests
}

func TestParseIPATemplate(t *testing.T) {
	tests := getParseIPATemplateTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseIPATemplate(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIPATemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_ParseIPATemplate(b *testing.B) {
	tests := getParseIPATemplateTests()
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseIPATemplate(tt.text)
			}
		})
	}
}

type ipaTestData struct {
	name          string
	text          string
	pronunciation Pronunciation
	want          Pronunciation
}

func getFindIPATests() []ipaTestData {
	tests := []ipaTestData{
		{

			text: `==English==
{{wikipedia|GDP (disambiguation)|GDP|lang=en}}

===Pronunciation===
* {{IPA|en|/==English==
{{wikipedia|GDP (disambiguation)|GDP|lang=en}}

===Pronunciation===
* {{IPA|en|/ˌd͡ʒiːdiːˈpiː/}}
* {{audio|en|en-us-GDP.ogg|a=US}}
* {{rhymes|en|iː|s=3}}`,
			pronunciation: Pronunciation{
				Word: "GDP",
			},
			want: Pronunciation{
				Word: "GDP",
				IPA: []IPATemplate{
					{
						Language: "en",
						IPA:      []string{"/ˌd͡ʒiːdiːˈpiː/"},
					},
				},
			},
		},

		{
			pronunciation: Pronunciation{
				Word: "umsonst",
			},
			text: `==German==

===Etymology===
{{affix|de|um|sonst}}, compare {{cog|da|omsonst}}.

===Pronunciation===
* {{IPA|de|/==German==

===Etymology===
{{affix|de|um|sonst}}, compare {{cog|da|omsonst}}.

===Pronunciation===
* {{IPA|de|/ʊmˈzɔnst/}}
* {{audio|de|De-umsonst.ogg}}

===Adverb===
{{de-adv}}

# [[free of charge]], [[gratis]]
#: {{syn|de|gratis|kostenlos|kostenfrei}}
# [[in vain]], without [[success]]
#: {{syn|de|vergebens|vergeblich|erfolglos}}
# for [[nothing]]; for the [[sake]] of doing it (without expecting reply)`,
			want: Pronunciation{
				Word: "umsonst",
				IPA: []IPATemplate{
					{
						Language: "de",
						IPA:      []string{"/ʊmˈzɔnst/"},
					},
				},
			},
		},

		{
			pronunciation: Pronunciation{
				Word: "Abraham men",
			},
			text: `First attested in ''The Fraternity of Vagabonds'' (1561) by {{w|John Awdely}}.

===Pronunciation===
* {{IPA|en|/, in which the beggar Lazarus ends up in [[Abraham]]'s bosom.

First attested in ''The Fraternity of Vagabonds'' (1561) by {{w|John Awdely}}.

===Pronunciation===
* {{IPA|en|/ˈeɪ.bɹəˌhæm mæn/|/ˈeɪ.bɹə.həm mæn/|a=US}}
* {{audio|en|LL-Q1860 (eng)-Persent101-Abraham man.wav|a=US}}

===Noun===
{{en-noun|Abraham men}}`,

			want: Pronunciation{
				Word: "Abraham men",
				IPA: []IPATemplate{
					{
						Language: "en",
						IPA: []string{
							"/ˈeɪ.bɹəˌhæm mæn/", "/ˈeɪ.bɹə.həm mæn/",
						},
						Variant: "a=US",
					},
				},
			},
		},
	}

	return tests
}

func TestFindIPA(t *testing.T) {
	tests := getFindIPATests()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindIPA(tt.text, tt.pronunciation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindIPA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_FindIPA(b *testing.B) {
	tests := getFindIPATests()

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FindIPA(tt.text, tt.pronunciation)
			}
		})
	}
}
