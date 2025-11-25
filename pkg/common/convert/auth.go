package convert

func TokenMapDB2Pb(tokenMapDB map[string]int) map[string]int32 {
    if tokenMapDB == nil {
        return nil
    }
    
    tokenMapPB := make(map[string]int32, len(tokenMapDB))
    for k, v := range tokenMapDB {
        tokenMapPB[k] = int32(v)
    }
    return tokenMapPB
}

func TokenMapPb2DB(tokenMapPB map[string]int32) map[string]int {
    if tokenMapPB == nil {
        return nil
    }
    
    tokenMapDB := make(map[string]int, len(tokenMapPB))
    for k, v := range tokenMapPB {
        tokenMapDB[k] = int(v)
    }
    return tokenMapDB
}