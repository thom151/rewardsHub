package dropbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DropboxListEntry struct {
	Tag         string `json:".tag"`
	Name        string `json:"name"`
	PathDisplay string `json:"path_display"`
	ID          string `json:"id"`
	Size        int64  `json:"size"`
}

type ListFolderResponse struct {
	Entries []DropboxListEntry `json:"entries"`
	Cursor  string             `json:"cursor"`
	HasMore bool               `json:"has_more"`
}

func ListDropboxFolder(ctx context.Context, path, accToken string) (entries []DropboxListEntry, err error) {
	type params struct {
		IncludeDeleted              bool   `json:"include_deleted"`
		IncludeMediaInfo            bool   `json:"include_media_info"`
		IncludeMountedFolders       bool   `json:"include_mounted_folders"`
		IncludeNonDownloadableFiles bool   `json:"include_non_downloadable_files"`
		Path                        string `json:"path"`
		Recursive                   bool   `json:"recursive"`
	}

	type continueParams struct {
		Cursor string `json:"cursor"`
	}

	listFolderUrl := "https://api.dropboxapi.com/2/files/list_folder"
	listFolderContinueUrl := "https://api.dropboxapi.com/2/files/list_folder/continue"

	var allEntries []DropboxListEntry

	firstReqBody := params{
		IncludeDeleted:              false,
		IncludeNonDownloadableFiles: true,
		IncludeMountedFolders:       true,
		IncludeMediaInfo:            false,
		Path:                        path,
		Recursive:                   true,
	}

	firstReqBytes, err := json.Marshal(firstReqBody)
	if err != nil {
		return []DropboxListEntry{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", listFolderUrl, bytes.NewBuffer(firstReqBytes))
	if err != nil {
		return []DropboxListEntry{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []DropboxListEntry{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return []DropboxListEntry{}, fmt.Errorf("dropbox error: %s", string(b))
	}

	var listResp ListFolderResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&listResp)
	if err != nil {
		return []DropboxListEntry{}, err
	}

	allEntries = append(allEntries, listResp.Entries...)

	for listResp.HasMore {
		fmt.Printf("folder has more\n")
		continueBody := continueParams{
			Cursor: listResp.Cursor,
		}

		continueBytes, err := json.Marshal(continueBody)
		if err != nil {
			return []DropboxListEntry{}, nil
		}

		req2, err := http.NewRequestWithContext(ctx, "POST", listFolderContinueUrl, bytes.NewBuffer(continueBytes))
		if err != nil {
			return []DropboxListEntry{}, err
		}
		req2.Header.Set("Authorization", "Bearer "+accToken)
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			return []DropboxListEntry{}, err
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != 200 {
			b, _ := io.ReadAll(resp2.Body)
			resp.Body.Close()
			return []DropboxListEntry{}, fmt.Errorf("dropbox error: %s", string(b))
		}

		var nextResp ListFolderResponse
		decoder := json.NewDecoder(resp2.Body)
		err = decoder.Decode(&nextResp)
		if err != nil {
			return []DropboxListEntry{}, err
		}
		allEntries = append(allEntries, nextResp.Entries...)
		listResp = nextResp

	}

	return allEntries, nil

}
