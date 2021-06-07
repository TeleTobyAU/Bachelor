package main

func Bwt32(x string, SA []int32, i int) string {
	xIndex := SA[i]
	if xIndex == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[xIndex-1])
	}
}

func GenerateOTable32(info *InfoInt32) {
	oTable := [][]int32{}
	alphabet := info.Alphabet
	sa := info.SA
	x := info.Input

	for range alphabet {
		oTable = append(oTable, []int32{0})
	}
	for i := range sa {
		for j := range alphabet {
			if Bwt32(x, sa, i) == alphabet[j] {
				oTable[j] = append(oTable[j], oTable[j][i]+1)
			} else {
				oTable[j] = append(oTable[j], oTable[j][i])
			}
		}
	}
	info.OTable = oTable
}

func InitBwtSearch32(info *InfoInt32) {
	n := len(info.SA)
	m := len(info.Key)
	key := info.Key
	alphabet := info.Alphabet

	L := int32(0)
	R := int32(n)

	if m > n {
		R = 0
		L = 1
	}
	i := m - 1
	for i >= 0 && L < R {
		//Find Index of key[i] in O table
		var a int
		for j := range alphabet {
			if string(key[i]) == alphabet[j] {
				a = j
			}
		}

		L = info.CTable[a] + info.OTable[a][L]
		R = info.CTable[a] + info.OTable[a][R]
		i -= 1
	}

	info.L = L
	info.R = R
}

func IndexBwtSearch32(info *InfoInt32) []int32 {
	var match []int32

	for i := int32(0); i < (info.R - info.L); i++ {
		match = append(match, info.SA[info.L+i])
	}

	return match
}
