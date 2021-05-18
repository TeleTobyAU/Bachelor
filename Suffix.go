package main

import (
	"fmt"
	"sort"
	"time"
)

const UNDEFINED = int(^uint(0) >> 1)

func SAIS(x string) []int {
	start := time.Now()
	n, alphSize := str2int(x)
	SA := make([]int, len(n))
	names := make([]int, len(n))
	sumString := make([]int, len(n))
	sumOffset := make([]int, len(n))
	LSTypes := make([]bool, len(n))
	maxAlph := len(n) + 1
	if alphSize > len(n) {
		maxAlph = alphSize
	}
	buckets := make([]int, maxAlph)
	bucketEnd := make([]int, maxAlph)

	sortSA(n, &SA, &names, &sumString, &sumOffset, &buckets, &bucketEnd, &LSTypes, alphSize)
	fmt.Println("total", time.Since(start))

	return SA
}

func str2int(x string) ([]int, int) {
	alpha := map[byte]int{}
	for i := range x {
		alpha[x[i]] = 1
	}
	tempAlph := make([]byte, len(alpha))
	i := 0
	for c := range alpha {
		tempAlph[i] = c
		i++
	}

	sort.Slice(tempAlph, func(i, j int) bool {
		return tempAlph[i] < tempAlph[j]
	})

	for i, c := range tempAlph {
		alpha[c] = i
	}

	out := make([]int, len(x))
	for i := range x {
		out[i] = alpha[x[i]]
	}

	return out, len(alpha)
}

func classifyLS(n []int, LSTypes *[]bool) {
	(*LSTypes)[len(n)-1] = true

	for i := len(n) - 2; i >= 0; i-- {
		if n[i] == n[i+1] {
			(*LSTypes)[i] = (*LSTypes)[i+1]
		} else {
			if n[i] > n[i+1] {
				(*LSTypes)[i] = false
			} else {
				(*LSTypes)[i] = true
			}
		}
	}
}

func isLMSIndex(LSString []bool, i int) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] && !LSString[i-1]
	}
}

func placeLMS(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	for i := range *SA {
		(*SA)[i] = UNDEFINED
	}

	findBucketEnds(alphSize, buckets, bucketEnds)

	//SA-IS step 1, placing LMS substrings in saisStruct
	for i := 0; i < len(n); i++ {
		if isLMSIndex(*LSTypes, i) {
			(*bucketEnds)[n[i]]--
			(*SA)[(*bucketEnds)[n[i]]] = i
		}
	}

}

func induceL(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketStarts *[]int) {
	bucketBeginnings(alphSize, buckets, bucketStarts)
	for i := 0; i < len(n); i++ {
		if (*SA)[i] == UNDEFINED {
			continue
		}

		if (*SA)[i] == 0 {
			continue
		}

		j := (*SA)[i] - 1
		if !(*LSTypes)[j] {
			(*SA)[(*bucketStarts)[n[j]]] = j
			(*bucketStarts)[n[j]]++
		}
	}
}

func induceS(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	findBucketEnds(alphSize, buckets, bucketEnds)
	for i := len(*SA); i > 0; i-- {
		if (*SA)[i-1] == 0 {
			continue
		}

		j := (*SA)[i-1] - 1
		if (*LSTypes)[j] {
			(*bucketEnds)[n[j]]--
			(*SA)[(*bucketEnds)[n[j]]] = j
		}

	}
}

func computeBuckets(n []int, buckets *[]int) {
	for i := 0; i < len(n); i++ {
		if n[i] != -1 {
			(*buckets)[n[i]]++
		}
	}
}

func findBucketEnds(alphSize int, buckets *[]int, bucketEnds *[]int) {
	(*bucketEnds)[0] = (*buckets)[0]
	for i := 1; i < alphSize; i++ {
		(*bucketEnds)[i] = (*bucketEnds)[i-1] + (*buckets)[i]
	}
}

func bucketBeginnings(alphSize int, buckets *[]int, bucketStarts *[]int) {
	(*bucketStarts)[0] = 0
	for i := 1; i < alphSize; i++ {
		(*bucketStarts)[i] = (*bucketStarts)[i-1] + (*buckets)[i-1]
	}
}

func equalLMS(n []int, LSTypes *[]bool, i int, j int) bool {
	if i == len(n) || j == len(n) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndex(*LSTypes, i+k)
		jLMS := isLMSIndex(*LSTypes, j+k)
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+k] != n[j+k] || ((*LSTypes)[i+k]) != ((*LSTypes)[j+k]) {
			return false
		}
		k++
	}
}

func reduceSA(n []int, SA *[]int, names *[]int, LSTypes *[]bool, newAlphSize *int, sumString *[]int, sumOffset *[]int, newStrLen *int) {
	name := 0

	for i := range *names {
		(*names)[i] = UNDEFINED
	}
	(*names)[(*SA)[0]] = name

	lastS := (*SA)[0]

	for i := 1; i < len(n); i++ {
		j := (*SA)[i]
		if !isLMSIndex(*LSTypes, j) {
			continue
		}
		if !equalLMS(n, LSTypes, lastS, j) {
			name++
		}
		lastS = j
		(*names)[j] = name
	}
	*newAlphSize = name + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = (*names)[i]
		if name == UNDEFINED {
			continue
		}

		(*sumOffset)[j] = i
		(*sumString)[j] = name
		j++
	}

	var temp []int
	for i := 0; i < len(*sumString); i++ {
		if (*sumString)[i] != 0 {
			temp = append(temp, (*sumString)[i])
		}
	}
	*sumString = append(temp, 0)

	*newStrLen = j - 1
}

func recursiveSorting(n []int, SA *[]int, names *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int, reducedString *[]int, reducedOffset *[]int, alphSize int) {
	classifyLS(n, LSTypes)

	computeBuckets(n, buckets)

	placeLMS(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	var newAlphSize int
	var newStrLen int

	reduceSA(n, SA, names, LSTypes, &newAlphSize, reducedString, reducedOffset, &newStrLen)

	newSa := make([]int, len(*reducedString))
	newNames := make([]int, len(*reducedString))
	newLSTypes := make([]bool, len(*reducedString))
	newSumStr := make([]int, len(*reducedString))
	newSumOffset := make([]int, len(*reducedString))
	newBuckets := make([]int, newAlphSize)
	newBucketEnds := make([]int, newAlphSize)
	sortSA(*reducedString, &newSa, &newNames, &newSumStr, &newSumOffset, &newBuckets, &newBucketEnds, &newLSTypes, newAlphSize)

	for i := 0; i < len(*SA); i++ {
		(*SA)[i] = UNDEFINED
	}

	remapLMS(n, buckets, bucketEnds, alphSize, newStrLen, &newSa, reducedOffset, SA)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)
}

func sortSA(n []int, SA *[]int, names *[]int, sumString *[]int, sumOffset *[]int, buckets *[]int, bucketEnds *[]int, LSTypes *[]bool, alphSize int) {
	if len(n) == 1 {
		(*SA)[0] = 0
		return
	}

	if alphSize == len(n) {
		(*SA)[0] = len(n) - 1
		for i := 0; i < len(n)-1; i++ {
			j := n[i]
			(*SA)[j] = i
		}
	} else {
		recursiveSorting(n, SA, names, LSTypes, buckets, bucketEnds, sumString, sumOffset, alphSize)
	}
}

func remapLMS(n []int, buckets *[]int, bucketEnds *[]int, alphSize int, reducedLength int, reducedSA *[]int, reducedOffset *[]int, SA *[]int) {
	findBucketEnds(alphSize, buckets, bucketEnds)
	for i := reducedLength + 1; i > 0; i-- {
		idx := (*reducedOffset)[(*reducedSA)[i-1]]
		(*bucketEnds)[n[idx]]--
		(*SA)[(*bucketEnds)[n[idx]]] = idx
	}
	(*SA)[0] = len(n) - 1
}
