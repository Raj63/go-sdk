package yeqown

import (
	"fmt"

	"github.com/Raj63/go-sdk/qrcode"
)

// CreateN implements qrcode.QRCode
func (y *yeqownQrcode) CreateN(input []qrcode.ProductInfo) ([]qrcode.QRResult, error) {
	numInput := len(input)
	results := make([]qrcode.QRResult, numInput)
	resultCh := make(chan struct {
		result qrcode.QRResult
		err    error
	})
	doneCh := make(chan struct{})
	defer close(resultCh)

	// Define the worker function that runs concurrently for each input product
	worker := func(index int) {
		qrResult, err := y.Create(input[index])
		select {
		case resultCh <- struct {
			result qrcode.QRResult
			err    error
		}{qrResult, err}:
		case <-doneCh:
		}
	}

	// Launch a goroutine for each input product
	for i := 0; i < numInput; i++ {
		go worker(i)
	}

	// Collect results from the result channel
	for i := 0; i < numInput; i++ {
		select {
		case res := <-resultCh:
			if res.err != nil {
				close(doneCh) // stop all goroutines by closing doneCh
				return nil, res.err
			}
			results[i] = res.result
		case <-doneCh:
			return nil, fmt.Errorf("stopped all the go routines")
		}
	}

	return results, nil
}
