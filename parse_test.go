package language

import "testing"

var tests = []string{
	"This is a test sentence that I would like to see.",
	"Some random text that is interesting.",
	"Roep hulle op die radio.",
	"Dis a kwessie van HONDERDE rande.",
	"My arm is in water.",
}

func TestParse(t *testing.T) {
	af, err := parse("data/afrikaans.txt")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	en, err := parse("data/esperanto.txt")
	if err != nil {
		t.Errorf("error: %v", err)
	}

	for i := range tests {
		scoreAf := matchString(af, tests[i])
		scoreEn := matchString(en, tests[i])
		ans := ""
		if scoreAf > scoreEn {
			ans = "AFRIKAANS"
		} else {
			ans = "ENGLISH"
		}
		t.Errorf("%q is %v (%.3f vs %.3f)", tests[i], ans, scoreAf, scoreEn)
	}
}
