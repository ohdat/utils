package utils

//DeIn 去掉交集
func DeIn(ids1, ids2 []int) ([]int, []int) {
	var ids2Map map[int]int
	for i := 0; i < len(ids2); i++ {
		ids2Map[ids2[i]] = ids2[i]
	}
	if len(ids2Map) > 0 {
		for i := 0; i < len(ids1); i++ {
			if _, ok := ids2Map[ids1[i]]; ok {
				delete(ids2Map, ids1[i])
				ids1 = append(ids1[:i], ids1[i+1:]...)
				i-- // form the remove item index to start iterate next item
			}
		}
	}
	var ids3 = make([]int, 0)
	if len(ids2Map) > 0 {
		for _, id_ := range ids2Map {
			ids3 = append(ids3, id_)
		}
	}

	return ids1, ids3
}
