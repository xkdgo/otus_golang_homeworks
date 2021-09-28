package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	s       = `[\,\!\.\:\?\'\"\ \r\n\@\#\$\%\^\&\*\+\~\/\>\<]`
	tire    = `(?:-)+`
	notTire = `[^-]+`
)

var (
	regex        = regexp.MustCompile(s)
	tireregex    = regexp.MustCompile(tire)
	nottireregex = regexp.MustCompile(notTire)
)

type freqinfo struct {
	word    string
	counter int
}

func isOnlyTire(str string) bool {
	matchedTire := tireregex.MatchString(str)
	matchedSomeOther := nottireregex.MatchString(str)
	return !matchedSomeOther && matchedTire
}

func sortInfoSlice(resultInfo []freqinfo) {
	sort.Slice(resultInfo, func(i, j int) bool {
		if resultInfo[i].counter == resultInfo[j].counter {
			return resultInfo[i].word < resultInfo[j].word
		}
		return resultInfo[i].counter > resultInfo[j].counter
	})
}

func buildFrequencyMap(words []string) map[string]int {
	frequency := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		if isOnlyTire(word) {
			continue
		}
		if word == "" {
			continue
		}
		frequency[word]++
	}
	return frequency
}

func Top10(text string) []string {
	splitwords := regex.Split(text, -1)
	var words []string
	for _, str := range splitwords {
		words = append(words, strings.Fields(str)...)
	}
	frequency := buildFrequencyMap(words)

	resultInfo := make([]freqinfo, 0, len(frequency))
	for word, counter := range frequency {
		if len(word) == 0 {
			continue
		}
		info := freqinfo{
			word:    word,
			counter: counter,
		}
		resultInfo = append(resultInfo, info)
	}
	sortInfoSlice(resultInfo)

	if len(resultInfo) >= 10 {
		resultInfo = resultInfo[:10]
	}
	result := make([]string, 0, 10)
	for _, info := range resultInfo {
		if info.counter == 0 {
			continue
		}
		result = append(result, info.word)
	}

	return result
}
