package main

import (
	"sort"
)

const UNDEFINED = int32(^uint32(0) >> 1)

func SAIS(x string) []int32 {
	n, alphSize := str2int32(x)
	SA := make([]int32, len(n))
	names := make([]int32, len(n))
	sumString := make([]int32, len(n))
	sumOffset := make([]int32, len(n))
	LSTypes := make([]bool, len(n))
	maxAlph := int32(len(n) + 1)
	if alphSize > int32(len(n)) {
		maxAlph = alphSize
	}
	buckets := make([]int32, maxAlph)
	bucketEnd := make([]int32, maxAlph)

	sortSA(n, &SA, &names, &sumString, &sumOffset, &buckets, &bucketEnd, &LSTypes, alphSize)

	return SA
}

func str2int32(x string) ([]int32, int32) {
	alpha := map[byte]int32{}
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
		return int32(tempAlph[i]) < int32(tempAlph[j])
	})

	for i, c := range tempAlph {
		alpha[c] = int32(i)
	}

	out := make([]int32, len(x))
	for i := range x {
		out[i] = alpha[x[i]]
	}

	return out, int32(len(alpha))
}

func classifyLS(n []int32, LSTypes *[]bool) {
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

func isLMSIndex(LSString []bool, i int32) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] && !LSString[i-1]
	}
}

func placeLMS(n []int32, alphSize int32, SA *[]int32, LSTypes *[]bool, buckets *[]int32, bucketEnds *[]int32) {
	for i := range *SA {
		(*SA)[i] = UNDEFINED
	}

	findBucketEnds(alphSize, buckets, bucketEnds)

	for i := 0; i < len(n); i++ {
		if isLMSIndex(*LSTypes, int32(i)) {
			(*bucketEnds)[n[i]]--
			(*SA)[(*bucketEnds)[n[i]]] = int32(i)
		}
	}
}

func induceL(n []int32, alphSize int32, SA *[]int32, LSTypes *[]bool, buckets *[]int32, bucketStarts *[]int32) {
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

func induceS(n []int32, alphSize int32, SA *[]int32, LSTypes *[]bool, buckets *[]int32, bucketEnds *[]int32) {
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

func computeBuckets(n []int32, buckets *[]int32) {
	for i := 0; i < len(n); i++ {
		if n[i] != -1 {
			(*buckets)[n[i]]++
		}
	}
}

func findBucketEnds(alphSize int32, buckets *[]int32, bucketEnds *[]int32) {
	(*bucketEnds)[0] = (*buckets)[0]
	for i := 1; int32(i) < alphSize; i++ {
		(*bucketEnds)[i] = (*bucketEnds)[i-1] + (*buckets)[i]
	}
}

func bucketBeginnings(alphSize int32, buckets *[]int32, bucketStarts *[]int32) {
	(*bucketStarts)[0] = 0
	for i := 1; int32(i) < alphSize; i++ {
		(*bucketStarts)[i] = (*bucketStarts)[i-1] + (*buckets)[i-1]
	}
}

func equalLMS(n []int32, LSTypes *[]bool, i int32, j int32) bool {
	if i == int32(len(n)) || j == int32(len(n)) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndex(*LSTypes, i+int32(k))
		jLMS := isLMSIndex(*LSTypes, j+int32(k))
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+int32(k)] != n[j+int32(k)] || ((*LSTypes)[i+int32(k)]) != ((*LSTypes)[j+int32(k)]) {
			return false
		}
		k++
	}
}

func reduceSA(n []int32, SA *[]int32, names *[]int32, LSTypes *[]bool, newAlphSize *int32, sumString *[]int32, sumOffset *[]int32, newStrLen *int32) {
	name := 0

	for i := range *names {
		(*names)[i] = UNDEFINED
	}
	(*names)[(*SA)[0]] = int32(name)

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
		(*names)[j] = int32(name)
	}
	*newAlphSize = int32(name) + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = int((*names)[i])
		if int32(name) == UNDEFINED {
			continue
		}

		(*sumOffset)[j] = int32(i)
		(*sumString)[j] = int32(name)
		j++
	}

	var temp []int32
	for i := 0; i < len(*sumString); i++ {
		if (*sumString)[i] != 0 {
			temp = append(temp, (*sumString)[i])
		}
	}
	*sumString = append(temp, 0)

	*newStrLen = int32(j - 1)
}

func recursiveSorting(n []int32, SA *[]int32, names *[]int32, LSTypes *[]bool, buckets *[]int32, bucketEnds *[]int32, reducedString *[]int32, reducedOffset *[]int32, alphSize int32) {
	classifyLS(n, LSTypes)

	computeBuckets(n, buckets)

	placeLMS(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	var newAlphSize int32
	var newStrLen int32

	reduceSA(n, SA, names, LSTypes, &newAlphSize, reducedString, reducedOffset, &newStrLen)

	newSa := make([]int32, len(*reducedString))
	newNames := make([]int32, len(*reducedString))
	newLSTypes := make([]bool, len(*reducedString))
	newSumStr := make([]int32, len(*reducedString))
	newSumOffset := make([]int32, len(*reducedString))
	newBuckets := make([]int32, newAlphSize)
	newBucketEnds := make([]int32, newAlphSize)
	sortSA(*reducedString, &newSa, &newNames, &newSumStr, &newSumOffset, &newBuckets, &newBucketEnds, &newLSTypes, newAlphSize)

	for i := 0; i < len(*SA); i++ {
		(*SA)[i] = UNDEFINED
	}

	remapLMS(n, buckets, bucketEnds, alphSize, newStrLen, &newSa, reducedOffset, SA)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)
}

func sortSA(n []int32, SA *[]int32, names *[]int32, sumString *[]int32, sumOffset *[]int32, buckets *[]int32, bucketEnds *[]int32, LSTypes *[]bool, alphSize int32) {
	if int32(len(n)) == 1 {
		(*SA)[0] = 0
		return
	}

	if alphSize == int32(len(n)) {
		(*SA)[0] = int32(len(n) - 1)
		for i := 0; i < len(n)-1; i++ {
			j := n[i]
			(*SA)[j] = int32(i)
		}
	} else {
		recursiveSorting(n, SA, names, LSTypes, buckets, bucketEnds, sumString, sumOffset, alphSize)
	}
}

func remapLMS(n []int32, buckets *[]int32, bucketEnds *[]int32, alphSize int32, reducedLength int32, reducedSA *[]int32, reducedOffset *[]int32, SA *[]int32) {
	findBucketEnds(alphSize, buckets, bucketEnds)
	for i := reducedLength + 1; i > 0; i-- {
		idx := (*reducedOffset)[(*reducedSA)[i-1]]
		(*bucketEnds)[n[idx]]--
		(*SA)[(*bucketEnds)[n[idx]]] = idx
	}
	(*SA)[0] = int32(len(n) - 1)
}
