package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var s = `[\,\!\.\:\?\'\"\ \r\n\@\#\$\%\^\&\*\+\~\/\>\<]`

var regex = regexp.MustCompile(s)

type freqinfo struct {
	word    string
	counter int
}

func sortInfoSlice(resultInfo []freqinfo) {
	sort.Slice(resultInfo, func(i, j int) bool {
		switch {
		case resultInfo[i].counter == resultInfo[j].counter:
			return resultInfo[i].word < resultInfo[j].word
		default:
			return resultInfo[i].counter > resultInfo[j].counter
		}
	})
}

func buildFrequencyMap(words []string) map[string]freqinfo {
	frequency := make(map[string]freqinfo)
	for _, word := range words {
		word = strings.ToLower(word)
		if word == "" {
			continue
		}
		if _, ok := frequency[word]; !ok {
			newInfo := freqinfo{
				word,
				1,
			}
			frequency[word] = newInfo
		} else {
			info := frequency[word]
			info.counter++
			frequency[word] = info
		}
	}
	return frequency
}

func Top10(text string) []string {
	splitwords := regex.Split(text, -1)
	var words []string
	for _, str := range splitwords {
		if str == "" {
			continue
		}
		if str == "-" {
			continue
		}

		words = append(words, strings.Fields(str)...)
	}
	frequency := buildFrequencyMap(words)

	resultInfo := make([]freqinfo, 0, len(frequency))
	for word, info := range frequency {
		if len(word) == 0 {
			continue
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
