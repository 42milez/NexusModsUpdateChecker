package api

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	Err "github.com/42milez/NexusModsWatcher/src/error"
	"github.com/42milez/NexusModsWatcher/src/log"
	"github.com/42milez/NexusModsWatcher/src/nexus"
	"github.com/42milez/NexusModsWatcher/src/util"
	"github.com/google/uuid"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const apiKeyName = "NEXUS_MODS_API_KEY"
const downloadConcurrentMax = 3
const (
	notStarted downloadStatus = iota
	inProgress
	completed
)

var NexusMods *nexusModsApi

type nexusModsApi struct {
	apiKey string
	client http.Client
}

func (p *nexusModsApi) GetRelease(mods nexus.ModInfoSet) ([]*nexus.Release, error) {
	getLatestFile := func(domain string, modId int, modName string, filter string) (*nexus.File, error) {
		url := fmt.Sprintf("https://api.nexusmods.com/v1/games/%s/mods/%d/files.json?category=main", domain, modId)

		var (
			req *http.Request
			err error
		)

		if req, err = http.NewRequest("GET", url, nil); err != nil {
			return nil, Err.CreateRequestFailed
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("apikey", p.apiKey)

		var resp *http.Response

		if resp, err = p.client.Do(req); err != nil {
			return nil, Err.RequestFailed
		}
		defer closeConn(resp)

		if err = outputNexusModsApiRateLimit(resp); err != nil {
			log.D(err.Error())
		}

		filesResp := &nexus.FilesApiResponse{}

		if err = json.NewDecoder(resp.Body).Decode(filesResp); err != nil {
			log.D(err.Error())
			return nil, Err.DecodeJsonFailed
		}

		if len(filesResp.Files) == 0 {
			return nil, Err.NoFileReceived
		}

		var latestFile nexus.File

		if len(filesResp.Files) > 1 {
			for _, f := range filesResp.Files {
				if strings.Contains(f.Name, filter) {
					latestFile = f
					break
				}
			}
		}
		if len(filesResp.Files) == 1 {
			latestFile = filesResp.Files[0]
		}
		log.D("latest file",
			fmt.Sprintf("domain:    %s", domain),
			fmt.Sprintf("id:        %d", modId),
			fmt.Sprintf("name:      %s", modName),
			fmt.Sprintf("version:   %s", latestFile.Version),
			fmt.Sprintf("size (mb): %0.2f", float64(latestFile.SizeKB)/1024))

		return &latestFile, nil
	}

	var releases []*nexus.Release

	for domain, modsInDomain := range mods {
		for _, mod := range modsInDomain {
			var (
				f   *nexus.File
				err error
			)

			if f, err = getLatestFile(domain, mod.ID, mod.Name, mod.Filter); err != nil {
				log.E(fmt.Sprintf("can't get latest file: modId=%d", mod.ID))
				return nil, Err.GetLatestFileFailed
			}

			if f.CmpTimestamp(mod.File.Timestamp) == nexus.NewerTimestamp {
				log.I(fmt.Sprintf("new version released: %s (%s)", mod.Name, f.Version))
				releases = append(releases, &nexus.Release{
					Domain: domain,
					Mod:    mod,
					File:   f,
				})
			}
		}
	}

	return releases, nil
}

func (p *nexusModsApi) Download(mods nexus.ModInfoSet) ([]*nexus.LocalFile, error) {
	getDownloadLink := func(ctx context.Context, mod *nexus.ModInfo, domain string) (string, error) {
		url := fmt.Sprintf(
			"https://api.nexusmods.com/v1/games/%s/mods/%d/files/%d/download_link.json",
			domain,
			mod.ID,
			mod.File.FileID)

		var (
			req *http.Request
			err error
		)

		if req, err = http.NewRequestWithContext(ctx, "GET", url, nil); err != nil {
			return "", Err.CreateRequestFailed
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("apikey", p.apiKey)

		var resp *http.Response
		defer closeConn(resp)

		errCh := make(chan error)

		go func() {
			resp, err = p.client.Do(req)
			errCh <- err
		}()

		select {
		case <-ctx.Done():
			return "", Err.CtxCanceled
		case e := <-errCh:
			if e != nil {
				return "", Err.RequestFailed
			}
		}

		if err = outputNexusModsApiRateLimit(resp); err != nil {
			log.D(err.Error())
		}

		dlResp := make(nexus.DownloadLinkApiResponse, 0)

		if err = json.NewDecoder(resp.Body).Decode(&dlResp); err != nil {
			return "", Err.DecodeJsonFailed
		}

		if len(dlResp) == 0 {
			return "", Err.NoDownloadLinkReceived
		}

		var uri string

		for _, dl := range dlResp {
			if dl.ShortName == "Nexus CDN" {
				uri = dl.URI
				break
			}
		}

		if uri == "" {
			return "", Err.GetDownloadLinkFailed
		}

		return uri, nil
	}

	downloadFile := func(ctx context.Context, mod *nexus.ModInfo, domain, uri string) (string, error) {
		var err error

		modDir := nexus.ModDir(mod, domain)
		if _, err = os.Stat(modDir); err != nil {
			if err = os.MkdirAll(modDir, 0755); err != nil {
				return "", Err.CreateDirectoryFailed
			}
		}

		dstPath := fmt.Sprintf("/%s/%s", modDir, strings.Replace(mod.File.FileName, " ", "_", -1))

		var tmpFile *os.File

		if tmpFile, err = os.Create(dstPath); err != nil {
			return "", Err.CreateFileFailed
		}
		defer closeFile(tmpFile)

		var req *http.Request

		if req, err = http.NewRequestWithContext(ctx, "GET", uri, nil); err != nil {
			return "", Err.CreateRequestFailed
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("apikey", p.apiKey)

		var resp *http.Response
		defer closeConn(resp)

		errCh := make(chan error)

		go func() {
			resp, err = p.client.Do(req)
			errCh <- err
		}()

		select {
		case <-ctx.Done():
			return "", Err.CtxCanceled
		case e := <-errCh:
			if e != nil {
				return "", Err.RequestFailed
			}
		}

		var bytes int64

		if bytes, err = io.Copy(tmpFile, resp.Body); err != nil {
			return "", Err.CopyFileFailed
		}

		if bytes != int64(mod.File.SizeInBytes) {
			log.E("file size mismatch",
				fmt.Sprintf("domain:        %s", domain),
				fmt.Sprintf("modId:         %d", mod.ID),
				fmt.Sprintf("name:          %s", mod.File.Name),
				fmt.Sprintf("version:       %s", mod.File.Version),
				fmt.Sprintf("size (expect): %d bytes", mod.File.SizeInBytes),
				fmt.Sprintf("size (actual): %d bytes", bytes))
			return "", Err.FileSizeMismatch
		}

		log.D(fmt.Sprintf("file downloaded: %s", dstPath))

		return dstPath, nil
	}

	completeCh := make(chan *nexus.ModInfo)
	errCh := make(chan error)

	ret := struct {
		LocalFiles []*nexus.LocalFile
		Mtx        sync.Mutex
	}{}

	download := func(ctx context.Context, mod *nexus.ModInfo, domain string) {
		log.D(fmt.Sprintf("download started: id=%d", mod.ID))

		var uri string
		var err error

		uri, err = getDownloadLink(ctx, mod, domain)
		if err != nil {
			errCh <- Err.GetDownloadLinkFailed
			return
		}

		var path string
		path, err = downloadFile(ctx, mod, domain, uri)
		if err != nil {
			errCh <- Err.DownloadFailed
			return
		}

		defer ret.Mtx.Unlock()
		ret.Mtx.Lock()

		ret.LocalFiles = append(ret.LocalFiles, &nexus.LocalFile{
			Mod:  mod,
			Path: path,
		})

		log.I("download completed",
			fmt.Sprintf("domain:  %s", domain),
			fmt.Sprintf("id:      %d", mod.ID),
			fmt.Sprintf("name:    %s", mod.Name),
			fmt.Sprintf("version: %s", mod.File.Version))

		completeCh <- mod
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	queue := &downloadQueue{}

	for domain, modsInDomain := range mods {
		for _, mod := range modsInDomain {
			queue.Push(mod, domain)
		}
	}

	// perform initial downloads
	for i := 0; i < downloadConcurrentMax; i += 1 {
		entry := queue.Pick()
		if entry == nil {
			break
		}
		go download(ctx, entry.Mod, entry.Domain)
	}

	for {
		select {
		case mod := <-completeCh:
			queue.Done(mod.UUID)
			if entry := queue.Pick(); entry != nil {
				go download(ctx, entry.Mod, entry.Domain)
			}
		case err = <-errCh:
			return nil, err
		}
		if queue.IncompleteCount() == 0 {
			break
		}
	}

	return ret.LocalFiles, nil
}

type downloadStatus int

type downloadQueueEntry struct {
	Domain string
	Mod    *nexus.ModInfo
	Status downloadStatus
}

type downloadQueue struct {
	queue list.List
}

func (p *downloadQueue) IncompleteCount() (ret int) {
	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Status == notStarted || entry.Status == inProgress {
			ret += 1
		}
	}
	return
}

func (p *downloadQueue) InProgressCount() (ret int) {
	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Status == inProgress {
			ret += 1
		}
	}
	return
}

func (p *downloadQueue) Done(id uuid.UUID) {
	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Mod.UUID == id {
			entry.Status = completed
			break
		}
	}
}

func (p *downloadQueue) Len() int {
	return p.queue.Len()
}

func (p *downloadQueue) Pick() *downloadQueueEntry {
	if p.queue.Len() == 0 {
		return nil
	}

	var onGoingDownloads int

	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Status == inProgress {
			onGoingDownloads += 1
		}
	}

	if onGoingDownloads == downloadConcurrentMax {
		return nil
	}

	var ret *downloadQueueEntry

	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Status == notStarted {
			entry.Status = inProgress
			ret = entry
			break
		}
	}

	return ret
}

func (p *downloadQueue) Push(mod *nexus.ModInfo, domain string) {
	entry := &downloadQueueEntry{
		Domain: domain,
		Mod:    mod,
		Status: notStarted,
	}
	p.queue.PushBack(entry)
}

func (p *downloadQueue) RemoveCompletedEntry() {
	for elem := p.queue.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*downloadQueueEntry)
		if entry.Status == completed {
			p.queue.Remove(elem)
		}
	}
}

func outputNexusModsApiRateLimit(resp *http.Response) error {
	now := time.Now()

	hRstAt := resp.Header.Get("X-RL-Hourly-Reset")
	hRstAt = strings.Replace(hRstAt, " ", "T", 1)
	hRstAt = strings.Replace(hRstAt, " +0000", "Z", 1)
	var hRstAfter int

	var t time.Time
	var err error

	if t, err = time.Parse(time.RFC3339, hRstAt); err != nil {
		return Err.ParseTimeFailed
	}

	hRstAfter = int(math.Ceil(t.Sub(now).Minutes()))

	dRstAt := resp.Header.Get("X-RL-Daily-Reset")
	dRstAt = strings.Replace(dRstAt, " ", "T", 1)
	dRstAt = strings.Replace(dRstAt, " +0000", "Z", 1)
	var dRstAfter int

	if t, err = time.Parse(time.RFC3339, dRstAt); err != nil {
		return Err.ParseTimeFailed
	}

	dRstAfter = int(math.Ceil(t.Sub(now).Minutes()))

	log.D(fmt.Sprintf("rate limit (nexus mods public api): %s, %s",
		fmt.Sprintf("hourly > %s/%s (%d min)",
			resp.Header.Get("X-RL-Hourly-Remaining"),
			resp.Header.Get("X-RL-Hourly-Limit"),
			hRstAfter),
		fmt.Sprintf("daily > %s/%s (%d min)",
			resp.Header.Get("X-RL-Daily-Remaining"),
			resp.Header.Get("X-RL-Daily-Limit"),
			dRstAfter)))

	return nil
}

func init() {
	NexusMods = &nexusModsApi{}

	var err error

	if NexusMods.apiKey, err = getSecret("auth", apiKeyName); err != nil {
		util.Exit(fmt.Errorf("%s (%s)", Err.GetSecretFailed, apiKeyName))
	}

	if NexusMods.apiKey == "" {
		util.Exit(Err.InvalidApiKey)
	}

	NexusMods.client = http.Client{
		Timeout: clientTimeout,
	}
}
