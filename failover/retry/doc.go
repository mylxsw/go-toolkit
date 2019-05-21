/*
Package retry 实现了一个错误重试函数。当函数执行出错时，能够自动根据设置的重试次数进行重试

	retryTimes, err := Retry(func(rt int) error {
		fmt.Printf("%d retry execute time: %s\n", rt, time.Now().String())
		if rt == 2 {
			return nil
		}

		return errors.New("test error")
	}, 3).Run()
*/
package retry
