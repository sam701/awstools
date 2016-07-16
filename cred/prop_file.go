package cred

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type propsFile struct {
	groups []*propsGroup
}

func (cf *propsFile) addGroup(g *propsGroup) {
	for i, v := range cf.groups {
		if v.name == g.name {
			cf.groups[i] = g
			return
		}
	}
	cf.groups = append(cf.groups, g)
}

func (cf *propsFile) write(w io.Writer) {
	for _, g := range cf.groups {
		fmt.Fprintf(w, "[%s]\n", g.name)
		for _, line := range g.lines {
			fmt.Fprintf(w, "%s = %s\n", line.key, line.value)
		}
	}
}

type propertyLine struct {
	key   string
	value string
}
type propsGroup struct {
	name  string
	lines []*propertyLine
}

func readPropsFile(r io.Reader) *propsFile {
	cf := &propsFile{[]*propsGroup{}}
	var gr *propsGroup
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '[' {
			gr = &propsGroup{
				name:  line[1 : len(line)-1],
				lines: []*propertyLine{},
			}
			cf.groups = append(cf.groups, gr)
		} else {
			pp := strings.Split(line, "=")
			for i, v := range pp {
				pp[i] = strings.TrimSpace(v)
			}
			gr.lines = append(gr.lines, &propertyLine{pp[0], pp[1]})
		}
	}
	return cf
}

func updatePropsGroup(fileName string, group *propsGroup) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	cf := readPropsFile(f)
	cf.addGroup(group)

	err = f.Truncate(0)
	if err != nil {
		log.Println(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Println(err)
	}

	cf.write(f)
}
