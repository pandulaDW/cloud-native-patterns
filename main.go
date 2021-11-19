package main

import (
	"cloud-native/patterns"
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	f := func(_ context.Context) (string, error) {
		return "", errors.New("concurrency error")
	}

	fNew := patterns.Breaker(f, 3)

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		res, err := fNew(ctx)
		if err != nil {
			fmt.Println("Error:", err.Error())
		} else {
			fmt.Println(res)
		}
	}

	time.Sleep(1 * time.Second)
	_, err := fNew(ctx)
	fmt.Println("Error:", err)

	_, err = fNew(ctx)
	fmt.Println("Error:", err)
}
