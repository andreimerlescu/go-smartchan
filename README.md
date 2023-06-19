# SmartChannel Package (`go_smartchan`)

The `go_smartchan` package provides a thread-safe, "smart" channel implementation for Go.

## Use Cases

### Handling large volumes of data 
In standard Go channels, you might run into an issue where you're trying to write to a channel that's full. This situation could block the goroutine that's trying to write to the channel and could potentially lead to a deadlock situation. By using the SmartChan, you can check if you can write to the channel before attempting to write, preventing the blocking situation.

### Safely managing channel closure 
In Go, closing a channel that's already closed will cause a panic. This can become an issue in complex concurrent programs, where it might not always be straightforward to manage channel lifecycle. The SmartChan struct includes a built-in atomic check to see if the channel is already closed before attempting to close it, avoiding a potential panic situation.

### Non-blocking reads
The Read method in SmartChan provides a non-blocking alternative to reading from a standard Go channel. In typical Go channels, trying to read from an empty channel will cause the goroutine to block until there's data available. In certain use cases, you might prefer to get an error immediately instead of blocking. The SmartChan.Read() method provides this capability.

### Counting the number of written items
In certain scenarios, you might need to keep a track of how many items have been written to a channel. With standard Go channels, you'd need to manage this separately, which can be error-prone in a concurrent scenario. SmartChan has a built-in atomic counter that keeps track of the number of written items, providing a threadsafe way to get this count.

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
	ch "github.com/andreimerlescu/go-smartchan"
	"time"
)

type Data struct {
	Name string
	Age  int
}

func main() {
	// Initialize a new SmartChan with a capacity of 5
	smartChan := ch.NewSmartChan(5)

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
