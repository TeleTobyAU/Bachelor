package main

type SaisStruct struct {
	SA             []int
	LSTypes        string
	buckets        []int
	beginnings     []int
	ends           []int
	LMSAlphabet    []string
	summaryOffsets []int
	summaryString  []int
	newStrLen      int
	newAlphSize    int
}

const UNDEFINED = int(^uint(0) >> 1)

/**
Linear suffix array construction by almost pure induced sorting.
Algorithm derived from Zhang and Chan 2009
*/
func SAIS(info *Info, n string) []int {
	saisStruct := new(SaisStruct)

	sortSA(info, saisStruct, n)

	return saisStruct.SA
}

func classifyLS(n string) string {
	LSTypes := make([]rune, len(n))
	LSTypes[len(n)-1] = 'S'

	for i := len(n) - 2; i >= 0; i-- {
		if n[i] == n[i+1] {
			LSTypes[i] = LSTypes[i+1]
		} else {
			if n[i] > n[i+1] {
				LSTypes[i] = 'L'
			} else {
				LSTypes[i] = 'S'
			}
		}
	}
	return string(LSTypes)
}

func isLMSIndex(LSString string, i int) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] == 'S' && LSString[i-1] != 'S'
	}
}

func placeLMS(saisStruct *SaisStruct, info *Info, n string) []int {
	//Initialize SA to UNDEFINED
	SA := make([]int, len(n))
	for i := range SA {
		SA[i] = UNDEFINED
	}

	//SA-IS step 1, placing LMS substrings in saisStruct
	for i := 0; i < len(n); i++ {
		if isLMSIndex(saisStruct.LSTypes, i) {
			remappedi := IndexOf(string(n[i]), info.alphabet)
			saisStruct.ends[remappedi]--
			SA[saisStruct.ends[remappedi]] = i
		}
	}
	return SA
}

func induceL(saisStruct *SaisStruct, info *Info, n string) []int {
	SA := saisStruct.SA
	for i := 0; i < len(n); i++ {
		if SA[i] == UNDEFINED {
			continue
		}

		if SA[i] == 0 {
			continue
		}

		j := SA[i] - 1
		if saisStruct.LSTypes[j] == 'L' {
			remappedi := IndexOf(string(n[j]), info.alphabet)
			SA[saisStruct.beginnings[remappedi]-1] = j
			saisStruct.beginnings[remappedi]++
		}
	}
	return SA
}

func induceS(saisStruct *SaisStruct, info *Info, n string) []int {
	//SA-IS step 3, placing and sorting S types
	SA := saisStruct.SA
	saisStruct.ends = bucketEnds(info, saisStruct)
	for i := len(n); i > 0; i-- {
		if saisStruct.SA[i-1] == 0 {
			continue
		}

		j := SA[i-1] - 1

		if saisStruct.LSTypes[j] == 'S' {
			remappedi := IndexOf(string(n[j]), info.alphabet)
			saisStruct.ends[remappedi]--
			SA[saisStruct.ends[remappedi]] = j
		}

	}
	return SA
}

func computeBuckets(info *Info) []int {
	buckets := make([]int, len(info.alphabet))
	for i := 0; i < len(info.input); i++ {
		remappedi := IndexOf(string(info.input[i]), info.alphabet)
		if remappedi != -1 {
			buckets[remappedi]++
		}
	}
	return buckets
}

func bucketEnds(info *Info, saisStruct *SaisStruct) []int {
	ends := make([]int, len(info.alphabet))
	ends[0] = saisStruct.buckets[0]
	for i := 1; i < len(info.alphabet); i++ {
		ends[i] = ends[i-1] + saisStruct.buckets[i]
	}
	return ends
}

func bucketBeginnings(info *Info, saisStruct *SaisStruct) []int {
	beginnings := make([]int, len(info.alphabet))
	beginnings[0] = saisStruct.buckets[0]
	for i := 1; i < len(info.alphabet); i++ {
		beginnings[i] = beginnings[i-1] + saisStruct.buckets[i-1]
	}
	return beginnings
}

func equalLMS(saisStruct *SaisStruct, n string, i int, j int) bool {
	if i == len(n) || j == len(n) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndex(saisStruct.LSTypes, i+k)
		jLMS := isLMSIndex(saisStruct.LSTypes, j+k)
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+k] != n[j+k] || (saisStruct.LSTypes[i+k] == 'S') != (saisStruct.LSTypes[j+k] == 'S') {
			return false
		}
		k++
	}
}

func reduceSA(saisStruct *SaisStruct, n string) {
	saisStruct.summaryOffsets = make([]int, len(n))
	name := 0
	names := make([]int, len(n)+1)
	for i := range names {
		names[i] = UNDEFINED
	}
	names[saisStruct.SA[0]] = name
	lastS := saisStruct.SA[0]

	for i := 1; i < len(n); i++ {
		j := saisStruct.SA[i]
		if !isLMSIndex(saisStruct.LSTypes, j) {
			continue
		}
		if !equalLMS(saisStruct, n, lastS, j) {
			name++
		}
		lastS = j
		names[j] = name
	}
	saisStruct.newAlphSize = name + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = names[i]
		if name == UNDEFINED {
			continue
		}
		saisStruct.summaryOffsets[j] = i
		saisStruct.summaryString = append(saisStruct.summaryString, name)
		j++
	}

	saisStruct.newStrLen = j - 1
}

func recursiveSorting(info *Info, saisStruct *SaisStruct, n string) []int {
	saisStruct.LSTypes = classifyLS(n)

	saisStruct.buckets = computeBuckets(info)

	saisStruct.beginnings = bucketBeginnings(info, saisStruct)

	saisStruct.ends = bucketEnds(info, saisStruct)

	saisStruct.SA = placeLMS(saisStruct, info, n)

	saisStruct.SA = induceL(saisStruct, info, n)

	saisStruct.beginnings = bucketBeginnings(info, saisStruct)

	saisStruct.SA = induceS(saisStruct, info, n)

	reduceSA(saisStruct, n)

	SA := saisStruct.SA

	//SA = sortSA(info, saisStruct, n)

	return SA
}

func sortSA(info *Info, saisStruct *SaisStruct, n string) []int {
	SA := saisStruct.SA
	if len(n) == 0 {
		SA[0] = 0
		return SA
	}

	if len(info.alphabet) == len(n)+1 {
		SA[0] = len(n)
		for i := 0; i < len(n); i++ {
			j := IndexOf(n[i], info.alphabet)
			SA[j] = i
		}
	} else {
		recursiveSorting(info, saisStruct, n)
	}

	return SA
}

func remapLMS(info *Info, saisStruct *SaisStruct, n string) []int {
	SA := saisStruct.SA
	saisStruct.ends = bucketEnds(info, saisStruct)
	for i := saisStruct.newStrLen + 1; i > 0; i-- {
		idx := saisStruct.summaryOffsets[saisStruct.summaryString[i-1]]
		bucketIdx := IndexOf(string(n[idx]), info.alphabet)
		SA[saisStruct.ends[bucketIdx]] = idx
	}
	SA[0] = len(n) - 1

	return SA
}
