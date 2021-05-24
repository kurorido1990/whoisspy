package whoisspy

var spyNum = map[int]int{
	8:  1,
	9:  2,
	16: 3,
	20: 4,
}

func getSpyNum(num int) int {
	currentSpyNum := 1

	for playerNum, spyNum := range spyNum {
		if playerNum >= num {
			currentSpyNum = spyNum
			continue
		} else {
			return currentSpyNum
		}
	}

	return 5
}

var winNum = 3
