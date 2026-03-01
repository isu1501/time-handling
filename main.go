package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func main() {

	//1.write a function that returns current time in this format: "2026-02-28 14:35:12"
	currentTime := getCurrentTime()
	fmt.Println(currentTime) //2026-03-01 09:31:20

	//2.Convert it into UTC time.Time.
	localTime := "28-02-2026 10:30 PM"
	UTCTime := convertToUTC(localTime)
	fmt.Println(UTCTime) //2026-02-28 22:30:00 +0000 UTC

	//3.Time Difference Calculator
	time1 := time.Date(2026, 02, 28, 10, 00, 00, 0, time.UTC)
	time2 := time.Date(2026, 02, 28, 10, 05, 32, 0, time.UTC)
	timeDifference := time2.Sub(time1)
	fmt.Println(timeDifference) //5m32s

	//4.Check expired token
	expiresAt := time.Date(2026, 02, 28, 10, 00, 00, 00, time.UTC)
	//use after function to check if the token is expired
	if time.Now().After(expiresAt) {
		fmt.Println("Token expired")
	} else {
		fmt.Println("Token is valid")
	}
	//use until function to check if the token is expired
	if time.Until(expiresAt) <= 0 {
		fmt.Println("Token expired")
	} else {
		fmt.Println("Token is valid")
	}
	//op - Token is expired

	//5.Timeout Wrapper
	err := RunWithTimeout(func() error {
		time.Sleep(2 * time.Second)
		return nil
	}, 1*time.Second)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	//op - Error: timeout exceeded

	//6.Safe Reusable Timer
	stop := make(chan struct{})
	go RepeatEvery(2*time.Second, func() {
		fmt.Println("tick")
	}, stop)
	time.Sleep(10 * time.Second)
	close(stop)

}

func RepeatEvery(d time.Duration, fn func(), stop <-chan struct{}) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			fn()
			timer.Reset(d)
		case <-stop:
			return
		}
	}
}

func RunWithTimeout(fn func() error, d time.Duration) error {
	done := make(chan error)

	//Run function in a seperate goroutine
	go func() {
		done <- fn()
	}()

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case err := <-done:
		return err
	case <-timer.C:
		return errors.New("timeout exceeded")
	}
}

func convertToUTC(localTime string) time.Time {
	layout := "02-01-2006 15:04 PM"
	parsedTime, err := time.Parse(layout, localTime)
	if err != nil {
		log.Fatalf("error parsing time: %v", err)
	}
	return parsedTime.UTC()
}

func getCurrentTime() string {

	return time.Now().Format("2006-01-02 15:04:05")
}
