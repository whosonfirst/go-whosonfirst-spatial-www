package writer

import (
	"context"
	"fmt"
	chain "github.com/g8rswimmer/error-chain"
	"io"
)

// Type MultiWriter holds mutltiple Writer instances.
type MultiWriter struct {
	Writer
	writers []Writer
}

// NewMultiWriter returns a Writer instance that will send all writes to each instance in 'writers'.
// Writes happen synchronolously in the order in which the underlying Writer instances are specified.
func NewMultiWriter(writers ...Writer) Writer {

	wr := &MultiWriter{
		writers: writers,
	}

	return wr
}

func (mw *MultiWriter) Write(ctx context.Context, key string, fh io.ReadSeeker) (int64, error) {

	errors := make([]error, 0)
	count := int64(0)

	for _, wr := range mw.writers {

		i, err := wr.Write(ctx, key, fh)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		count += i
		_, err = fh.Seek(0, 0)

		if err != nil {
			return count, err
		}
	}

	if len(errors) > 0 {

		err := fmt.Errorf("One or more Write operations failed")
		err = errorChain(err, errors...)

		return count, err
	}

	return count, nil
}

func (mw *MultiWriter) WriterURI(ctx context.Context, key string) string {
	return ""
}

func (mw *MultiWriter) Close(ctx context.Context) error {

	errors := make([]error, 0)

	for _, wr := range mw.writers {

		err := wr.Close(ctx)

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {

		err := fmt.Errorf("One or more Close operations failed")
		err = errorChain(err, errors...)

		return err
	}

	return nil
}

func errorChain(first error, others ...error) error {
	ec := chain.New()
	ec.Add(first)

	for _, e := range others {
		ec.Add(e)
	}

	return ec
}
