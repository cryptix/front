// Package front is a frontmatter extraction library.
package front

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sigs.k8s.io/yaml"
	"strings"
)

var (
	//ErrIsEmpty is an error indicating no front matter was found
	ErrIsEmpty = errors.New("front: no front matter found")
)

//Matter is all what matters here.
type Matter struct {
	Delim string
}

//NewMatter creates a new Matter instance
func NewMatter(delim string) *Matter {
	return &Matter{
		Delim: delim,
	}
}

// JSONToMap parses the input and extract JSON frontmatter as map[string]interface{}
func (m *Matter) JSONToMap(input io.Reader) (front map[string]interface{}, body string, err error) {
	f, body, err := m.splitFront(input)
	if err != nil {
		return map[string]interface{}{}, body, err
	}
	front, err = JSONHandler(f)
	if err != nil {
		return nil, "", err
	}
	return front, body, nil
}

// YAMLToMap parses the input and extract YAML frontmatter as map[string]interface{}
func (m *Matter) YAMLToMap(input io.Reader) (front map[string]interface{}, body string, err error) {
	f, body, err := m.splitFront(input)
	if err != nil {
		return map[string]interface{}{}, body, err
	}
	front, err = YAMLHandler(f)
	if err != nil {
		return nil, "", err
	}
	return front, body, nil
}

// YAMLToJSON parses the input and extract YAML frontmatter as []byte containing JSON
func (m *Matter) YAMLToJSON(input io.Reader) (front []byte, body string, err error) {
	f, body, err := m.splitFront(input)
	if err != nil {
		return []byte{}, body, err
	}
	front, err = yaml.YAMLToJSON([]byte(f))
	if err != nil {
		return nil, "", err
	}
	return front, body, nil
}

func sniffDelim(input []byte) (string, error) {
	if len(input) < 4 {
		return "", ErrIsEmpty
	}
	return string(input[:3]), nil
}

func (m *Matter) splitFront(input io.Reader) (front, body string, err error) {
	bufsize := 1024 * 1024
	buf := make([]byte, bufsize)

	s := bufio.NewScanner(input)
	// Necessary so we can handle larger than default 4096b buffer
	s.Buffer(buf, bufsize)

	var bodyBuilder strings.Builder
	s.Split(m.split)
	n := 0
	for s.Scan() {
		text := s.Text()
		hasDelim := strings.HasPrefix(text, m.Delim)
		if n == 0 && hasDelim {
			front = strings.TrimSpace(text[3:])
		} else if n == 1 && hasDelim {
			bodyBuilder.WriteString(text[3:])
		} else {
			bodyBuilder.WriteString(text)
		}
		n++
	}
	body = strings.TrimSpace(bodyBuilder.String())
	if len(front) < 3 {
		return front, body, ErrIsEmpty
	}
	return front, body, nil
}

//split implements bufio.SplitFunc for spliting front matter from the body text.
func (m *Matter) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	delim, err := sniffDelim(data)
	if err != nil || delim != m.Delim {
		return len(data), data, nil
	}
	if x := bytes.Index(data, []byte(delim)); x >= 0 {
		// check the next delim index
		if next := bytes.Index(data[x+len(delim):], []byte(delim)); next > 0 {
			return next + len(delim), data[:next+len(delim)], nil
		}
		return len(data), data, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func dropSpace(d []byte) []byte {
	return bytes.TrimSpace(d)
}

//JSONHandler decodes JSON string into a go map[string]interface{}
func JSONHandler(front string) (map[string]interface{}, error) {
	var rst interface{}
	err := json.Unmarshal([]byte(front), &rst)
	if err != nil {
		return nil, err
	}
	return rst.(map[string]interface{}), nil
}

//YAMLHandler decodes yaml string into a go map[string]interface{}
func YAMLHandler(front string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(front), &out)
	if err != nil {
		return nil, err
	}
	// clean maps
	for k, v := range out {
		out[k] = convert(v)
	}
	return out, nil
}

// convert converts all map[interface{}]interface{} children to map[string]interface{}
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
