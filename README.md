# SmartChannel Package (`go_smartchan`)

The `go_smartchan` package provides a thread-safe, "smart" channel implementation for Go.

## Description

This package wraps the basic Go channel with additional functionalities such as:

- Keeping track of the number of writes made to the channel (`Count`).
- Allowing safe concurrent access with internal locking mechanism.
- Checking whether more data can be written to the channel (`CanWrite`).
- Gracefully handling attempts to write to or read from a closed channel.

## Installation

Use the standard `go get` command to install this package:

```bash
go get github.com/andreimerlescu/go-smartchan
```

## Usage

Here's an example of how to use the `SmartChan` type:

```go
package main

import (
	"context"
	"fmt"
	"github.com/andreimerlescu/go-smartchan"
	"time"
)

type Data struct {
	Name string
	Age  int
}

func main() {
	// Initialize a new SmartChan with a capacity of 5
	smartChan := go_smartchan.NewSmartChan(5)

	// Create a context that will automatically cancel after 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Launch a goroutine to write data to the channel
	go func() {
		for i := 0; i < 10; i++ {
			err := smartChan.Write(Data{Name: fmt.Sprintf("Name%d", i), Age: i})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	// Read data from the channel until the context is done
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done, exiting read loop.")
			smartChan.Close()
			return
		case iData, ok := <-smartChan.Chan():
			if ok {
				data, ok := iData.(Data)
				if !ok {
					fmt.Println("Could not convert data to type Data")
					continue
				}
				fmt.Println("Name read from channel:", data.Name)
			} else {
				fmt.Println("Channel closed, exiting read loop.")
				return
			}
		}
	}
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT
