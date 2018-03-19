package retry

import (
	"fmt"
	"time"
	"testing"
	"errors"
)

func TestRetryLatter(t *testing.T) {
	fmt.Println("current time: " + time.Now().String())

	retryTimes, err := Retry(func(rt int) error {
		fmt.Printf("%d retry execute time: %s\n", rt, time.Now().String())
		if rt == 2 {
			return nil
		}

		return errors.New("test error")
	}, 3).Run()

	if err != nil {
		fmt.Println("still error: " + err.Error())
	}

	fmt.Printf("retry %d times\n", retryTimes)

	<-Retry(func(rt int) error {
		fmt.Printf("%d retry execute time: %s\n", rt, time.Now().String())
		if rt == 2 {
			return nil
		}

		return errors.New("test error")
	}, 3).Success(func(rt int) {
		fmt.Printf("retry %d times\n", rt)
	}).Failed(func(err error) {
		fmt.Println("still error: " + err.Error())
	}).RunAsync()

	fmt.Println("logic")

	time.Sleep(10 * time.Second)
}
