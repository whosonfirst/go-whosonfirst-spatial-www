package writer

import (
	"context"
	"fmt"
	"io"
	"os"
)

type CwdWriter struct {
	Writer
	writer Writer
}

func init() {

	ctx := context.Background()

	schemes := []string{
		"cwd",
	}

	for _, scheme := range schemes {

		err := RegisterWriter(ctx, scheme, NewCwdWriter)

		if err != nil {
			panic(err)
		}
	}
}

func NewCwdWriter(ctx context.Context, uri string) (Writer, error) {

	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive current working directory, %w", err)
	}

	uri = fmt.Sprintf("fs://%s", cwd)
	fs_wr, err := NewFileWriter(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new FS writer, %w", err)
	}

	wr := &CwdWriter{
		writer: fs_wr,
	}

	return wr, nil
}

func (wr *CwdWriter) Write(ctx context.Context, path string, fh io.ReadSeeker) (int64, error) {

	return wr.writer.Write(ctx, path, fh)
}

func (wr *CwdWriter) WriterURI(ctx context.Context, path string) string {
	return wr.writer.WriterURI(ctx, path)
}

func (wr *CwdWriter) Close(ctx context.Context) error {
	return wr.writer.Close(ctx)
}
