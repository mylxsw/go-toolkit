package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetryPanic(t *testing.T) {
	retryTimes, err := Retry(func(retryTimes int) error {
		if retryTimes < 1 {
			panic("sorry")
		}
		return nil
	}, 2).Run()
	if err != nil {
		t.Errorf("still error: %s", err.Error())
	}

	if retryTimes != 2 {
		t.Errorf("test failed, expect %d, got %d", 2, retryTimes)
	}
}

func TestRetryLatter(t *testing.T) {
	fmt.Println("current time: " + time.Now().String())

	retryTimes, err := Retry(func(rt int) error {
		fmt.Printf("%d retry execute time: %s\n", rt, time.Now().String())
		if rt < 1 {
			return errors.New("test error")
		}
		return nil
	}, 3).Run()

	if err != nil {
		t.Errorf("still error: %s", err.Error())
	}

	if retryTimes != 2 {
		t.Errorf("test failed, expect %d, got %d", 2, retryTimes)
	}

	fmt.Printf("retry %d times\n", retryTimes)

	succeed := false
	<-Retry(func(rt int) error {
		fmt.Printf("%d retry execute time: %s\n", rt, time.Now().String())
		if rt < 2 {
			return errors.New("test error")
		}
		return nil
	}, 3).Success(func(rt int) {
		fmt.Printf("retry %d times\n", rt)
		succeed = true
	}).Failed(func(err error) {
		fmt.Println("still error: " + err.Error())
		succeed = false
	}).RunAsync()


	if !succeed {
		t.Errorf("test failed")
	}

}
