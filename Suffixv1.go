package main

import (
	"fmt"
	"sort"
	"time"
)

const UNDEFINEDv1 = int(^uint(0) >> 1)

func SAISv1(x string) []int {
	start := time.Now()
	n, alphSize := str2intv1(x)
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

	sortSAv1(n, &SA, &names, &sumString, &sumOffset, &buckets, &bucketEnd, &LSTypes, alphSize)
	fmt.Println("total", time.Since(start))

	return SA
}

func str2intv1(x string) ([]int, int) {
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

func classifyLSv1(n []int, LSTypes *[]bool) {
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

func isLMSIndexv1(LSString []bool, i int) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] && !LSString[i-1]
	}
}

func placeLMSv1(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	for i := range *SA {
		(*SA)[i] = UNDEFINEDv1
	}

	findBucketEndsv1(alphSize, buckets, bucketEnds)

	//SA-IS step 1, placing LMS substrings in saisv1 Struct
	for i := 0; i < len(n); i++ {
		if isLMSIndexv1(*LSTypes, i) {
			(*bucketEnds)[n[i]]--
			(*SA)[(*bucketEnds)[n[i]]] = i
		}
	}

}

func induceLv1(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketStarts *[]int) {
	bucketBeginningsv1(alphSize, buckets, bucketStarts)
	for i := 0; i < len(n); i++ {
		if (*SA)[i] == UNDEFINEDv1 {
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

func induceSv1(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	findBucketEndsv1(alphSize, buckets, bucketEnds)
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

func computeBucketsv1(n []int, buckets *[]int) {
	for i := 0; i < len(n); i++ {
		if n[i] != -1 {
			(*buckets)[n[i]]++
		}
	}
}

func findBucketEndsv1(alphSize int, buckets *[]int, bucketEnds *[]int) {
	(*bucketEnds)[0] = (*buckets)[0]
	for i := 1; i < alphSize; i++ {
		(*bucketEnds)[i] = (*bucketEnds)[i-1] + (*buckets)[i]
	}
}

func bucketBeginningsv1(alphSize int, buckets *[]int, bucketStarts *[]int) {
	(*bucketStarts)[0] = 0
	for i := 1; i < alphSize; i++ {
		(*bucketStarts)[i] = (*bucketStarts)[i-1] + (*buckets)[i-1]
	}
}

func equalLMSv1(n []int, LSTypes *[]bool, i int, j int) bool {
	if i == len(n) || j == len(n) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndexv1(*LSTypes, i+k)
		jLMS := isLMSIndexv1(*LSTypes, j+k)
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+k] != n[j+k] || ((*LSTypes)[i+k]) != ((*LSTypes)[j+k]) {
			return false
		}
		k++
	}
}

func reduceSAv1(n []int, SA *[]int, names *[]int, LSTypes *[]bool, newAlphSize *int, sumString *[]int, sumOffset *[]int, newStrLen *int) {
	name := 0

	for i := range *names {
		(*names)[i] = UNDEFINEDv1
	}
	(*names)[(*SA)[0]] = name

	lastS := (*SA)[0]

	for i := 1; i < len(n); i++ {
		j := (*SA)[i]
		if !isLMSIndexv1(*LSTypes, j) {
			continue
		}
		if !equalLMSv1(n, LSTypes, lastS, j) {
			name++
		}
		lastS = j
		(*names)[j] = name
	}
	*newAlphSize = name + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = (*names)[i]
		if name == UNDEFINEDv1 {
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

func recursiveSortingv1(n []int, SA *[]int, names *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int, reducedString *[]int, reducedOffset *[]int, alphSize int) {
	classifyLSv1(n, LSTypes)

	computeBucketsv1(n, buckets)

	placeLMSv1(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceLv1(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceSv1(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	var newAlphSize int
	var newStrLen int

	reduceSAv1(n, SA, names, LSTypes, &newAlphSize, reducedString, reducedOffset, &newStrLen)

	newSa := make([]int, len(*reducedString))
	newNames := make([]int, len(*reducedString))
	newLSTypes := make([]bool, len(*reducedString))
	newSumStr := make([]int, len(*reducedString))
	newSumOffset := make([]int, len(*reducedString))
	newBuckets := make([]int, newAlphSize)
	newBucketEnds := make([]int, newAlphSize)
	sortSAv1(*reducedString, &newSa, &newNames, &newSumStr, &newSumOffset, &newBuckets, &newBucketEnds, &newLSTypes, newAlphSize)

	for i := 0; i < len(*SA); i++ {
		(*SA)[i] = UNDEFINEDv1
	}

	remapLMSv1(n, buckets, bucketEnds, alphSize, newStrLen, &newSa, reducedOffset, SA)

	induceLv1(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceSv1(n, alphSize, SA, LSTypes, buckets, bucketEnds)
}

func sortSAv1(n []int, SA *[]int, names *[]int, sumString *[]int, sumOffset *[]int, buckets *[]int, bucketEnds *[]int, LSTypes *[]bool, alphSize int) {
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
		recursiveSortingv1(n, SA, names, LSTypes, buckets, bucketEnds, sumString, sumOffset, alphSize)
	}
}

func remapLMSv1(n []int, buckets *[]int, bucketEnds *[]int, alphSize int, reducedLength int, reducedSA *[]int, reducedOffset *[]int, SA *[]int) {
	findBucketEndsv1(alphSize, buckets, bucketEnds)
	for i := reducedLength + 1; i > 0; i-- {
		idx := (*reducedOffset)[(*reducedSA)[i-1]]
		(*bucketEnds)[n[idx]]--
		(*SA)[(*bucketEnds)[n[idx]]] = idx
	}
	(*SA)[0] = len(n) - 1
}
