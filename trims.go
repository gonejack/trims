package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
)

type opts struct {
	Overwrite bool     `short:"w" name:"overwrite" help:"Overwrite source file."`
	About     bool     `help:"Show about."`
	File      []string `arg:"" optional:""`
}

type trims struct {
	opts
}

func (c *trims) run() (err error) {
	kong.Parse(&c.opts,
		kong.Name("trims"),
		kong.Description("Trim spaces from lines from files."),
		kong.UsageOnError(),
	)
	if c.About {
		fmt.Println("Visit https://github.com/gonejack/trims")
		return
	}
	if len(c.File) == 0 {
		c.File = append(c.File, "-")
	}
	for _, f := range c.File {
		err = c.process(f)
		if err != nil {
			return fmt.Errorf("process %s failed: %w", f, err)
		}
	}
	return
}
func (c *trims) process(file string) (err error) {
	src, dst := os.Stdin, os.Stdout
	if file != "-" {
		src, err = os.OpenFile(file, os.O_RDWR, 0755)
		if err != nil {
			return
		}
		defer src.Close()
	}
	if src != os.Stdin && c.Overwrite {
		dst, err = os.CreateTemp(os.TempDir(), "")
		if err != nil {
			return
		}
		defer os.Remove(dst.Name())
		defer dst.Close()
		defer func() {
			if err == nil {
				src.Truncate(0)
				src.Seek(0, io.SeekStart)
				dst.Seek(0, io.SeekStart)
				_, err = io.Copy(src, dst)
			}
		}()
	}
	dt := bufio.NewWriter(dst)
	sc := bufio.NewScanner(src)
	for sc.Scan() {
		dat := bytes.TrimSpace(sc.Bytes())
		if len(dat) > 0 {
			_, err = dt.Write(dat)
			if err != nil {
				return
			}
			_, err = dt.WriteRune('\n')
			if err != nil {
				return
			}
		}
	}
	err = sc.Err()
	if err != nil {
		return
	}
	err = dt.Flush()
	if err != nil {
		return
	}
	return
}
