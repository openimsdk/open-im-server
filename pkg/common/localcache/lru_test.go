package localcache

//func TestName(t *testing.T) {
//	target := &cacheTarget{}
//	l := NewCache[string](100, 1000, time.Second*20, time.Second*5, target, nil)
//	//l := NewLRU[string, string](1000, time.Second*20, time.Second*5, target)
//
//	fn := func(key string, n int, fetch func() (string, error)) {
//		for i := 0; i < n; i++ {
//			//v, err := l.Get(key, fetch)
//			//if err == nil {
//			//	t.Log("key", key, "value", v)
//			//} else {
//			//	t.Error("key", key, err)
//			//}
//			l.Get(key, fetch)
//			//time.Sleep(time.Second / 100)
//		}
//	}
//
//	tmp := make(map[string]struct{})
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10000; i++ {
//		wg.Add(1)
//		key := fmt.Sprintf("key_%d", i%200)
//		tmp[key] = struct{}{}
//		go func() {
//			defer wg.Done()
//			//t.Log(key)
//			fn(key, 10000, func() (string, error) {
//				//time.Sleep(time.Second * 3)
//				//t.Log(time.Now(), "key", key, "fetch")
//				//if rand.Uint32()%5 == 0 {
//				//	return "value_" + key, nil
//				//}
//				//return "", errors.New("rand error")
//				return "value_" + key, nil
//			})
//		}()
//
//		//wg.Add(1)
//		//go func() {
//		//	defer wg.Done()
//		//	for i := 0; i < 10; i++ {
//		//		l.Del(key)
//		//		time.Sleep(time.Second / 3)
//		//	}
//		//}()
//	}
//	wg.Wait()
//	t.Log(len(tmp))
//	t.Log(target.String())
//
//}
