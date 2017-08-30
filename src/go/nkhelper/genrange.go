// +build ignore

package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

var source = `#ifndef JPRANGE_H____
#define JPRANGE_H____
const nk_rune nk_font_japanese_glyph_ranges[] = {
0x0020, 0x007f,
{{range .ranges}}{{printf "0x%04x, 0x%04x,\n" (index . 0) (index . 1)}}{{end}}0
};
#endif
`

func main() {
	m := map[uint32]struct{}{}
	for _, url := range []string{
		"https://encoding.spec.whatwg.org/index-jis0208.txt",
		"https://encoding.spec.whatwg.org/index-jis0212.txt",
	} {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		tsv := csv.NewReader(resp.Body)
		tsv.Comma = '\t'
		tsv.Comment = '#'
		for {
			record, err := tsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			var i uint32
			if _, err = fmt.Sscanf(record[1], "0x%04x", &i); err != nil {
				log.Fatal(err)
			}
			m[i] = struct{}{}
		}
	}
	a := []uint32{}
	for v := range m {
		a = append(a, v)
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	ranges := [][2]uint32{}
	l, prev := a[0], a[0]
	for _, v := range a[1:] {
		if prev+1 == v {
			prev++
			continue
		}
		ranges = append(ranges, [2]uint32{l, prev})
		l = v
		prev = v
	}
	ranges = append(ranges, [2]uint32{l, prev})

	f, err := os.Create("jprange.h")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t := template.Must(template.New("").Parse(source))
	if err = t.Execute(f, map[string]interface{}{
		"ranges": ranges,
	}); err != nil {
		log.Fatal(err)
	}
}
