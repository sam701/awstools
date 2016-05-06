package cred

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type credentialsFile struct {
	groups []*credentialsGroup
}

func (cf *credentialsFile) addGroup(g *credentialsGroup) {
	for i, v := range cf.groups {
		if v.profile == g.profile {
			cf.groups[i] = g
			return
		}
	}
	cf.groups = append(cf.groups, g)
}

func (cf *credentialsFile) write(w io.Writer) {
	for _, g := range cf.groups {
		fmt.Fprintf(w, "[%s]\n", g.profile)
		for _, line := range g.lines {
			fmt.Fprintf(w, "%s = %s\n", line.key, line.value)
		}
	}
}

type propertyLine struct {
	key   string
	value string
}
type credentialsGroup struct {
	profile string
	lines   []*propertyLine
}

func credentialsFilePath() string {
	return path.Join(os.Getenv("HOME"), ".aws/credentials")
}
func SaveCredentials(profile, keyId, secret, token string) {
	f, err := os.OpenFile(credentialsFilePath(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	cf := readCredentials(f)
	cf.addGroup(newCredentialsGroup(profile, keyId, secret, token))

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

func newCredentialsGroup(profile, keyId, secret, token string) *credentialsGroup {
	g := &credentialsGroup{
		profile: profile,
		lines: []*propertyLine{
			&propertyLine{"aws_access_key_id", keyId},
			&propertyLine{"aws_secret_access_key", secret},
		},
	}
	if token != "" {
		g.lines = append(g.lines, &propertyLine{"aws_session_token", token})
	}
	return g
}

func readCredentials(r io.Reader) *credentialsFile {
	cf := &credentialsFile{[]*credentialsGroup{}}
	var gr *credentialsGroup
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line[0] == '[' {
			gr = &credentialsGroup{
				profile: line[1 : len(line)-1],
				lines:   []*propertyLine{},
			}
			cf.groups = append(cf.groups, gr)
		} else {
			pp := strings.Split(line, " = ")
			gr.lines = append(gr.lines, &propertyLine{pp[0], pp[1]})
		}
	}
	return cf
}

func GetMainAccountKeyId(profileName string) string {
	f, err := os.Open(credentialsFilePath())
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	cf := readCredentials(f)
	for _, g := range cf.groups {
		if g.profile == profileName {
			for _, l := range g.lines {
				if l.key == "aws_access_key_id" {
					return l.value
				}
			}
		}
	}
	log.Fatalf("Credentials file does not contain profile '%s'\n", profileName)
	return ""
}
