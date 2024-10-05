package main

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"

	"github.com/orisano/gosax"
	"github.com/schollz/progressbar/v3"

	bzip2 "github.com/cosnicolaou/pbzip2"

	"github.com/spf13/pflag"

	"github.com/terwey/wiktionary-ipa-extract/ipa"
)

var (
	inPage  bool
	inTitle bool
	found   bool
	inText  bool

	inputFile  string
	outputFile string

	bz bool
)

type preprocess struct {
	buf  string
	word ipa.Pronunciation
}

var tracefile string
var cpuprofile string

func init() {
	pflag.StringVarP(&inputFile, "input", "i", "", "input file")
	pflag.StringVarP(&outputFile, "output", "o", "", "output file")
	pflag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	pflag.StringVar(&tracefile, "trace", "", "trace file")
	pflag.BoolVar(&bz, "bz", false, "use bzip2")
}

func main() {

	log.SetFlags(0)

	if err := pflag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
		pflag.Usage()
		os.Exit(1)
	}

	if cpuprofile != "" {
		log.Printf("cpuprofile: %s", cpuprofile)
		// Create a file to write the CPU profile data to
		cpf, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer cpf.Close() // Ensure the file is closed when we're done

		// Start CPU profiling
		if err := pprof.StartCPUProfile(cpf); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile() // Stop the CPU profile when the program exits
	}

	if tracefile != "" {
		log.Printf("tracefile: %s", tracefile)
		f, err := os.Create(tracefile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := trace.Start(f); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	input := os.Stdin

	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		input = f
	}

	writer := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		writer = f
	}

	process(input, writer)
}

func process(input io.Reader, writer io.Writer) {
	bufferedWriter := bufio.NewWriterSize(writer, 8*1024*1024)

	bar := progressbar.Default(-1, "Finding IPA's")

	workers := 4
	size := 1000 * workers
	data := make(chan preprocess, 10*size)
	output := make(chan ipa.Pronunciation, size)

	var word ipa.Pronunciation

	var wg sync.WaitGroup
	var wgWriter sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(d chan preprocess, o chan ipa.Pronunciation, wg *sync.WaitGroup) {
			defer wg.Done()

			for p := range d {
				buf := ipa.FindIPA(p.buf, p.word)
				if len(buf.IPA) != 0 {
					o <- buf
				}
			}
		}(data, output, &wg)
	}

	wgWriter.Add(1)
	go func(o chan ipa.Pronunciation, w *bufio.Writer, wgW *sync.WaitGroup) {
		defer wgW.Done()
		for p := range o {
			bar.Add(1)
			w.Write(p.JSON())
			w.Write([]byte("\n"))
		}
	}(output, bufferedWriter, &wgWriter)

	var ipaBuf strings.Builder

	var r *gosax.Reader
	if bz {
		r = gosax.NewReader(bzip2.NewReader(context.Background(), input))
	} else {
		r = gosax.NewReader(input)
	}

	for {
		e, err := r.Event()
		if err != nil {
			log.Fatal(err)
		}
		if e.Type() == gosax.EventEOF {
			break
		}

		switch e.Type() {
		case gosax.EventStart:
			name, _ := gosax.Name(e.Bytes)
			switch string(name) {
			case "page":
				inPage = true
				ipaBuf.Reset()

			case "title":
				inTitle = true

			case "text":
				inText = true
			}
		case gosax.EventEnd:
			name, _ := gosax.Name(e.Bytes)
			switch string(name) {
			case "page":
				inPage = false
				found = false

			case "title":
				inTitle = false

			case "text":
				inText = false
				if ipaBuf.Len() != 0 {
					data <- preprocess{buf: ipaBuf.String(), word: word}
				}
			}

		case gosax.EventText:
			if inTitle {
				if bytes.Contains(e.Bytes, []byte("Wiktionary")) {
					continue
				}
				word = ipa.Pronunciation{
					Word: string(e.Bytes),
				}
			}

			if inText {
				// we need to somehow check if this is sufficient,
				// but I suspect it won't matter much if we lose a few words
				if bytes.Contains(e.Bytes, []byte("IPA")) {
					found = true
					ipaBuf.Write(e.Bytes)
				}

				if found {
					if bytes.HasPrefix(e.Bytes, []byte("}}")) {
						found = false
						ipaBuf.WriteString("}}")
					} else {
						ipaBuf.Write(e.Bytes)
					}
				}
			}
		default:
		}
	}

	close(data)

	wg.Wait()
	close(output)

	wgWriter.Wait()
}
