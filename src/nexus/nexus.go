package nexus

import (
	"context"
	"fmt"
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"github.com/42milez/NexusModsWatcher/src/log"
	"github.com/42milez/NexusModsWatcher/src/util"
	"github.com/google/go-github/v35/github"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ModInfo struct {
	ID      int    `yaml:"id"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Path    string `yaml:"path"`
	Filter  string `yaml:"filter,omitempty"`
	Author  string `yaml:"author"`
}

type ModInfoSet map[string][]ModInfo

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

	for domain, mods := range *p {
		for _, mod := range mods {
			dirPath := fmt.Sprintf("%s/%s", domain, mod.Path)
			if _, err = os.Stat(dirPath); err != nil {
				if err = os.MkdirAll(dirPath, 0755); err != nil {
					return Err.CreateDirectoryFailed
				}
			}
		}
	}

	return nil
}

func (p *ModInfoSet) Update(domain string, id int, version string) bool {
	for k, v := range *p {
		if k != domain {
			continue
		}
		for i, mod := range v {
			if mod.ID == id {
				(*p)[k][i].Version = version
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
	ID     uuid.UUID
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
	Path    string
	Release *Release
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

func closeFile(f *os.File) {
	if f == nil {
		return
	}
	if err := f.Close(); err != nil {
		util.Exit(Err.CloseFileFailed)
	}
}
