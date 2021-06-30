package nexus

import (
	"context"
	"fmt"
	Err "github.com/42milez/NexusModsUpdateChecker/src/error"
	"github.com/42milez/NexusModsUpdateChecker/src/log"
	"github.com/42milez/NexusModsUpdateChecker/src/util"
	"github.com/google/go-github/v35/github"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	SameTimestamp TsCondition = iota
	NewerTimestamp
	OlderTimestamp
)

type TsCondition int

type ModInfo struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Filter string `yaml:"filter,omitempty"`
	Author string `yaml:"author"`
	File   struct {
		FileID      int    `yaml:"fileId"`
		FileName    string `yaml:"fileName"`
		Name        string `yaml:"name"`
		SizeInBytes int    `yaml:"sizeInBytes"`
		Timestamp   int    `yaml:"timestamp"`
		Version     string `yaml:"version"`
	}
	UUID uuid.UUID `yaml:"uuid"`
}

type ModInfoSet map[string][]*ModInfo

func (p *ModInfoSet) Export(filePath string) error {
	var out *os.File
	var err error

	if out, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		return Err.OpenFileFailed
	}
	defer closeFile(out)

	enc := yaml.NewEncoder(out)
	enc.SetIndent(2)

	if err = enc.Encode(p); err != nil {
		log.D(err.Error())
		return Err.EncodeYamlFailed
	}

	return nil
}

func (p *ModInfoSet) Setup(f string) error {
	var in *os.File
	var err error

	if in, err = os.Open(f); err != nil {
		return Err.OpenFileFailed
	}
	defer closeFile(in)

	if err = yaml.NewDecoder(in).Decode(p); err != nil {
		log.D(err.Error())
		return Err.DecodeYamlFailed
	}

	isZeroUUID := func(id uuid.UUID) bool {
		tmp1 := [16]byte(id)
		var tmp2 [16]byte
		return reflect.DeepEqual(tmp1, tmp2)
	}

	for _, modsInDomain := range *p {
		for _, mod := range modsInDomain {
			if isZeroUUID(mod.UUID) {
				mod.UUID = uuid.New()
			}
		}
	}

	return nil
}

func (p *ModInfoSet) Update(release *Release) bool {
	for k, modsInDomain := range *p {
		if k != release.Domain {
			continue
		}
		for _, mod := range modsInDomain {
			if mod.ID == release.Mod.ID {
				mod.File.FileID = release.File.FileID
				mod.File.FileName = release.File.FileName
				mod.File.Name = release.File.Name
				mod.File.SizeInBytes = release.File.SizeInBytes
				mod.File.Version = release.File.Version
				mod.File.Timestamp = release.File.UploadedTimestamp
				return true
			}
		}
	}
	return false
}

type File struct {
	ID                   []int     `json:"id"`
	UID                  int       `json:"uid"`
	FileID               int       `json:"file_id"`
	Name                 string    `json:"name"`
	Version              string    `json:"version"`
	CategoryID           int       `json:"category_id"`
	CategoryName         string    `json:"category_name"`
	IsPrimary            bool      `json:"is_primary"`
	Size                 int       `json:"size"`
	FileName             string    `json:"file_name"`
	UploadedTimestamp    int       `json:"uploaded_timestamp"`
	UploadedTime         time.Time `json:"uploaded_time"`
	ModVersion           string    `json:"mod_version"`
	ExternalVirusScanUrl string    `json:"external_virus_scan_url"`
	Description          string    `json:"description"`
	SizeKB               int       `json:"size_kb"`
	SizeInBytes          int       `json:"size_in_bytes"`
	ChangelogHTML        string    `json:"changelog_html"`
	ContentPreviewLink   string    `json:"content_preview_link"`
}

func (p *File) CmpTimestamp(ts int) TsCondition {
	if p.UploadedTimestamp > ts {
		return NewerTimestamp
	}
	if p.UploadedTimestamp < ts {
		return OlderTimestamp
	}
	return SameTimestamp
}

type FileUpdate struct {
	OldFileID         int       `json:"old_file_id"`
	NewFileID         int       `json:"new_file_id"`
	OldFileName       string    `json:"old_file_name"`
	NewFileName       string    `json:"new_file_name"`
	UploadedTimestamp int       `json:"uploaded_timestamp"`
	UploadedTime      time.Time `json:"uploaded_time"`
}

type FilesApiResponse struct {
	Files       []File
	FileUpdates []FileUpdate
}

type Release struct {
	Domain string
	Mod    *ModInfo
	File   *File
}

type DownloadLink struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	URI       string `json:"URI"`
}

type DownloadLinkApiResponse []DownloadLink

type LocalFile struct {
	Mod  *ModInfo
	Path string
}

func (p *LocalFile) Move(dstDir string) (string, error) {
	if _, err := os.Stat(dstDir); err == nil {
		return "", Err.DirectoryAlreadyExists
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return "", Err.CreateDirectoryFailed
	}

	filePath := dstDir + "/" + p.Mod.File.FileName

	if err := os.Rename(p.Path, filePath); err != nil {
		return "", Err.MoveFileFailed
	}

	return filePath, nil
}

func CreateReleaseNote(releases []*Release) (string, error) {
	gh := github.NewClient(nil)

	var input string
	for _, r := range releases {
		input += fmt.Sprintf(
			"- [%s: %s](%s)\n",
			r.Mod.Name,
			r.File.Version,
			fmt.Sprintf("https://www.nexusmods.com/%s/mods/%d?tab=files", r.Domain, r.Mod.ID))
	}

	opt := &github.MarkdownOptions{
		Mode:    "gfm",
		Context: "google/go-github",
	}

	output, _, err := gh.Markdown(context.Background(), input, opt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func ModDir(mod *ModInfo, domain string) string {
	dirName := mod.Name
	dirName = strings.ToLower(dirName)
	dirName = strings.Replace(dirName, " ", "_", -1)
	dirName = regexp.MustCompile("[-!$%^&*()_+|~=`{}\\[\\]:\";'<>?,./]").ReplaceAllLiteralString(dirName, "_")
	dirName = regexp.MustCompile(`_{2,}`).ReplaceAllString(dirName, "_")
	dirName = strings.TrimPrefix(dirName, "_")
	dirName = strings.TrimSuffix(dirName, "_")
	return fmt.Sprintf("%s/%s/%s/%s/%d/%s", util.WorkDir, "mods", domain, dirName, mod.File.Timestamp, mod.File.Version)
}

func closeFile(f *os.File) {
	if f == nil {
		return
	}
	if err := f.Close(); err != nil {
		util.Exit(Err.CloseFileFailed)
	}
}
