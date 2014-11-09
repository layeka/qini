// Copyright 2014 layeka Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package qini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	regDoubleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$")
	regSingleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*'([^']*)'$")
	regNoQuote     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([^#;]+)")
	regNoValue     = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*([#;].*)?")
)

type QIni struct {
	//section key value
	currSection string
	data        map[string]map[string]string
}

func Load(filename string) *QIni {
	ini := &QIni{currSection: "default", data: make(map[string]map[string]string)}
	ini.data[ini.currSection] = make(map[string]string)
	if f, ok := os.Open(filename); ok == nil {
		defer f.Close()
		reader := bufio.NewReader(f)
		ini.parseReader(reader)
	}
	return ini
}
func (this *QIni) parseReader(reader *bufio.Reader) {
	if b, _, err := reader.ReadLine(); err == nil {
		line := strings.TrimFunc(string(b), unicode.IsSpace)
		if len(line) > 0 {
			for line[len(line)-1] == '\\' {
				line = line[:len(line)-1]
				if b, _, err := reader.ReadLine(); err == nil {
					line += strings.TrimFunc(string(b), unicode.IsSpace)
				}
			}
			this.parseLine(line)
		}
		this.parseReader(reader)
	}
}
func (this *QIni) addSection(section string) {
	section = strings.ToLower(section)
	if _, ok := this.data[section]; !ok {
		this.data[section] = make(map[string]string)
	}
	this.currSection = section
}
func (this *QIni) addValue(key string, value string) {
	this.data[this.currSection][strings.ToLower(key)] = value
}
func (this *QIni) parseLine(line string) {
	// commets
	if line[0] == ';' || line[0] == '#' {
		return
	}
	// section name
	if line[0] == '[' && line[len(line)-1] == ']' {
		section := strings.TrimFunc(line[1:len(line)-1], unicode.IsSpace)
		this.addSection(section)
		return
	}
	// key = value
	if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
		this.addValue(m[0][1], m[0][2])
	} else if m = regSingleQuote.FindAllStringSubmatch(line, 1); m != nil {
		this.addValue(m[0][1], m[0][2])
	} else if m = regNoQuote.FindAllStringSubmatch(line, 1); m != nil {
		this.addValue(m[0][1], strings.TrimFunc(m[0][2], unicode.IsSpace))
	} else if m = regNoValue.FindAllStringSubmatch(line, 1); m != nil {
		this.addValue(m[0][1], "")
	}
}
func (this *QIni) GetValue(section string, key string) (string, error) {
	if s, ok := this.data[strings.ToLower(section)]; ok {
		if v, ok := s[strings.ToLower(key)]; ok {
			return v, nil
		}
	}
	return "", errors.New(fmt.Sprintf("the section %s not exists the key %s", section, key))
}
func (this *QIni) DefaultString(section string, key string, value string) string {
	if v, err := this.GetValue(section, key); err == nil {
		return v
	}
	return value
}
func (this *QIni) DefaultBool(section string, key string, value bool) bool {
	if v, err := this.GetValue(section, key); err == nil {
		lowerStr := strings.ToUpper(v)
		if lowerStr == "T" || lowerStr == "TRUE" || lowerStr == "Y" || lowerStr == "YES" || lowerStr == "1" {
			return true
		}
		return false
	}
	return value
}
func (this *QIni) DefaultInt(section string, key string, value int) int {
	if v, err := this.GetValue(section, key); err == nil {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return value
}
func (this *QIni) DefaultInt64(section string, key string, value int64) int64 {
	if v, err := this.GetValue(section, key); err == nil {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return value
}
func (this *QIni) DefaultFloat(section string, key string, value float32) float32 {
	if v, err := this.GetValue(section, key); err == nil {
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		}
	}
	return value
}
func (this *QIni) DefaultFloat64(section string, key string, value float64) float64 {
	if v, err := this.GetValue(section, key); err == nil {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return value
}
