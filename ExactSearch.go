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

func InitBwtSearch32(exact *BwtExact32) {
	n := len(exact.bwtTable.SA)
	m := len(exact.Key)
	key := exact.Key
	alphabet := exact.bwtTable.Alphabet
	CTable := exact.bwtTable.CTable
	OTable := exact.bwtTable.OTable

	L := int32(0)
	R := int32(n)

	if m > n {
		R = 0
		L = 1
	}
	i := m - 1
	for i >= 0 && L < R {
		var a int
		for j := range alphabet {
			if string(key[i]) == alphabet[j] {
				a = j
			}
		}

		L = CTable[a] + OTable[a][L]
		R = CTable[a] + OTable[a][R]
		i -= 1
	}
	exact.L = L
	exact.R = R
}

func IndexBwtSearch32(exact *BwtExact32) []int32 {
	var match []int32

	for i := int32(0); i < (exact.R - exact.L); i++ {
		match = append(match, exact.bwtTable.SA[exact.L+i])
	}
	return match
}
