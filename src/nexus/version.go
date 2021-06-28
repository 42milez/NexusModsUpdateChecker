package nexus

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	SameVersion Condition = iota
	NewerVersion
	OlderVersion
	InvalidVersion
)

type Condition int

type VersionPair struct {
	Main int
	Sub  int
}

type Version struct {
	major VersionPair
	minor VersionPair
	patch VersionPair
	label int
	orig  string
}

func (p *Version) Cmp(ver *Version) Condition {
	if p.orig == ver.orig {
		return SameVersion
	}
	if p.major.Main > ver.major.Main {
		return NewerVersion
	}
	if p.major.Sub < ver.major.Sub {
		return OlderVersion
	}
	if p.minor.Main > ver.minor.Main {
		return NewerVersion
	}
	if p.minor.Sub < ver.minor.Sub {
		return OlderVersion
	}
	if p.patch.Main > ver.patch.Main {
		return NewerVersion
	}
	if p.patch.Sub < ver.patch.Sub {
		return OlderVersion
	}
	if p.label > ver.label {
		return NewerVersion
	}
	return OlderVersion
}

func (p *Version) Parse(ver string) bool {
	if ver == "" {
		return false
	}

	p.orig = ver

	ver = strings.TrimPrefix(ver, "v")
	versions := strings.Split(ver, ".")
	count := len(versions)

	if count > 3 {
		return false
	}

	var label string
	if tmp := strings.Split(versions[count-1], "-"); len(tmp) > 1 {
		if len(tmp) > 2 {
			return false
		}
		versions[count-1] = tmp[0]
		label = splitLabel(tmp[1])
	}

	parse := func(v string) *VersionPair {
		pair := &VersionPair{}
		tmp := splitVersion(v)
		var err error
		if pair.Main, err = strconv.Atoi(tmp[0]); err != nil {
			return nil
		}
		if len(tmp) > 1 {
			runes := []rune(tmp[1])
			if len(runes) > 1 {
				return nil
			}
			n := int(runes[0])
			if n < 97 || 122 < n {
				return nil
			}
			pair.Sub = n
		}
		return pair
	}

	var pair *VersionPair
	if pair = parse(versions[0]); pair == nil {
		return false
	}
	p.major.Main = pair.Main
	p.major.Sub = pair.Sub

	if count > 1 {
		if pair = parse(versions[1]); pair == nil {
			return false
		}
		p.minor.Main = pair.Main
		p.minor.Sub = pair.Sub
	}

	if count > 2 {
		if pair = parse(versions[2]); pair == nil {
			return false
		}
		p.patch.Main = pair.Main
		p.patch.Sub = pair.Sub
	}

	if label != "" {
		var err error
		if p.label, err = strconv.Atoi(label); err != nil {
			return false
		}
	}

	return true
}

func (p *Version) String() string {
	return p.orig
}

func splitVersion(arg string) (ret []string) {
	r := regexp.MustCompile(`[a-z]`)
	idx := r.FindIndex([]byte(arg))
	if idx == nil {
		ret = append(ret, arg)
		return
	}
	ret = append(ret, arg[:idx[0]])
	ret = append(ret, arg[idx[0]:])
	return
}

func splitLabel(arg string) string {
	r := regexp.MustCompile(`[0-9]`)
	idx := r.FindIndex([]byte(arg))
	if idx == nil {
		return arg
	}
	return arg[idx[0]:]
}
