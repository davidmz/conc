package conc_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/davidmz/go-conc"
)

func Example_single() {
	conc.Run(func(onDispose conc.OnDispose) error {
		fmt.Println("Hello #1")
		onDispose(func() { fmt.Println("Dispose #1") })
		return nil
	})
	// Output:
	// Hello #1
	// Dispose #1
}

func Example_parallel() {
	conc.Run(
		func(onDispose conc.OnDispose) error {
			fmt.Println("Hello #1.1")
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Hello #1.2")
			time.Sleep(10 * time.Millisecond)
			onDispose(func() { fmt.Println("Dispose #1") })
			return nil
		},
		func(onDispose conc.OnDispose) error {
			onDispose(func() { fmt.Println("Dispose #2") })
			time.Sleep(5 * time.Millisecond)
			fmt.Println("Hello #2.1")
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Hello #2.2")
			return nil
		},
	)
	// Output:
	// Hello #1.1
	// Hello #2.1
	// Hello #1.2
	// Hello #2.2
	// Dispose #2
	// Dispose #1
}

func Example_tree() {
	err := conc.Run(
		func(onDispose conc.OnDispose) error {
			fmt.Println("Hello #1.1")
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Hello #1.2")
			time.Sleep(10 * time.Millisecond)
			onDispose(func() { fmt.Println("Dispose #1") })
			return nil
		},
		conc.Tasks(
			func(onDispose conc.OnDispose) error {
				onDispose(func() { fmt.Println("Dispose #2") })
				time.Sleep(5 * time.Millisecond)
				fmt.Println("Hello #2.1")
				time.Sleep(10 * time.Millisecond)
				fmt.Println("Hello #2.2")
				return nil
			},
			func(onDispose conc.OnDispose) error {
				time.Sleep(5 * time.Millisecond)
				onDispose(func() { fmt.Println("Dispose #3") })
				time.Sleep(9 * time.Millisecond)
				fmt.Println("Hello #3.1")
				time.Sleep(10 * time.Millisecond)
				fmt.Println("Hello #3.2")
				return nil
			},
		),
	)
	fmt.Println(err)
	// Output:
	// Hello #1.1
	// Hello #2.1
	// Hello #1.2
	// Hello #3.1
	// Hello #2.2
	// Hello #3.2
	// Dispose #2
	// Dispose #3
	// Dispose #1
	// <nil>
}

func Example_parallel_errors() {
	err := conc.Run(
		func(onDispose conc.OnDispose) error {
			fmt.Println("Hello #1.1")
			return errors.New("Error #1")
		},
		func(onDispose conc.OnDispose) error {
			time.Sleep(5 * time.Millisecond)
			fmt.Println("Hello #2.1")
			return errors.New("Error #2")
		},
	)
	fmt.Println(err)
	// Output:
	// Hello #1.1
	// Hello #2.1
	// Error #1
	// Error #2
}
