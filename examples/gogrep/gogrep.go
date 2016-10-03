package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bketelsen/trace"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func search(ctx context.Context, root string, pattern string) ([]string, error) {
	t, ctx := trace.NewContext(ctx, "gogrep", "search")
	defer t.Finish()
	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string, 100)
	// get all the paths

	g.Go(func() error {
		tr, ctx := trace.NewContext(ctx, "gogrep", "walk")
		defer tr.Finish()
		defer close(paths)

		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			tr, ctx := trace.NewContext(ctx, "gogrep", "filepathwalk")
			defer tr.Finish()
			if err != nil {
				tr.SetError()
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if !info.IsDir() && !strings.HasSuffix(info.Name(), ".go") {
				return nil
			}

			select {
			case paths <- path:
			case <-ctx.Done():
				tr.SetError()
				return ctx.Err()
			}
			return nil
		})

	})

	c := make(chan string, 100)
	for path := range paths {
		p := path
		g.Go(func() error {

			tr, ctx := trace.NewContext(ctx, "gogrep", "filesearch")
			defer tr.Finish()
			data, err := ioutil.ReadFile(p)
			if err != nil {
				tr.SetError()
				return err
			}
			if !bytes.Contains(data, []byte(pattern)) {
				return nil
			}
			select {
			case c <- p:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}
	go func() {
		g.Wait()
		close(c)
	}()

	var m []string
	for r := range c {
		m = append(m, r)
	}
	return m, g.Wait()
}
