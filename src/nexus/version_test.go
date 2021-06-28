package nexus

import (
	"testing"
)

func TestVersion_Parse(t *testing.T) {
	v := "1.0.0"
	ver := &Version{}
	if got := ver.Parse(v); got != true {
		t.Errorf("Version.Parse() = %t; want %t", got, true)
	}

	v = "1.0.0.0"
	ver = &Version{}
	if got := ver.Parse(v); got != false {
		t.Errorf("Version.Parse() = %t; want %t", got, false)
	}

	v = "1.0.0-r1-r1"
	ver = &Version{}
	if got := ver.Parse(v); got != false {
		t.Errorf("Version.Parse() = %t; want %t", got, false)
	}
}

func TestVersion_Cmp(t *testing.T) {
	v := "1.0.0"
	ver := &Version{}
	ver.Parse(v)
	want := &Version{
		major: VersionPair{
			Main: 1,
			Sub:  0,
		},
		minor: VersionPair{
			Main: 0,
			Sub:  0,
		},
		patch: VersionPair{
			Main: 0,
			Sub:  0,
		},
		label: 0,
		orig:  v,
	}
	if got := ver.Cmp(want); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1 := &Version{}
	v2 := &Version{}
	v1.Parse("1.2.3")
	v2.Parse("1.2.3")
	if got := v1.Cmp(v2); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1.Parse("v1.09a")
	v2.Parse("v1.09a")
	if got := v1.Cmp(v2); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1.Parse("v1.09a.09a")
	v2.Parse("v1.09a.09a")
	if got := v1.Cmp(v2); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1.Parse("0.01a-r1")
	v2.Parse("0.01a-r1")
	if got := v1.Cmp(v2); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1.Parse("0.01.01a-r1")
	v2.Parse("0.01.01a-r1")
	if got := v1.Cmp(v2); got != SameVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, SameVersion)
	}

	v1.Parse("1.20.0")
	v2.Parse("1.10.0")
	if got := v1.Cmp(v2); got != NewerVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, NewerVersion)
	}

	v1.Parse("2.0.0")
	v2.Parse("1.10.15")
	if got := v1.Cmp(v2); got != NewerVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, NewerVersion)
	}

	v1.Parse("v1.20a")
	v2.Parse("v1.10b")
	if got := v1.Cmp(v2); got != NewerVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, NewerVersion)
	}

	v1.Parse("v1.09a")
	v2.Parse("v1.10b")
	if got := v1.Cmp(v2); got != OlderVersion {
		t.Errorf("Version.Cmp() = %d; want = %d", got, OlderVersion)
	}
}
