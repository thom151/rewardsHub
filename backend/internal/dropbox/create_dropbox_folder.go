package dropbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func CreateDropboxFolder(ctx context.Context, path, accToken string) (err error) {
	body := []byte(fmt.Sprintf(`{"path": "%s"}`, path))
	createFolderUrl := "https://api.dropboxapi.com/2/files/create_folder_v2"

	req, err := http.NewRequestWithContext(ctx, "POST", createFolderUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dropbox error: %s", string(b))
	}
	return nil
}

func CreateDropboxSharedLink(ctx context.Context, path, accToken string) (link ShareLinkResp, err error) {
	createSharedLink := "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"

	type setting struct {
		Access              string `json:"access"`
		AllowDownload       bool   `json:"allow_download"`
		Audience            string `json:"audience"`
		RequestedVisibility string `json:"requested_visibility"`
	}
	type params struct {
		Path     string  `json:"path"`
		Settings setting `json:"settings"`
	}

	shareLinkBody := params{
		Path: path,
		Settings: setting{
			Access:              "viewer",
			AllowDownload:       true,
			Audience:            "public",
			RequestedVisibility: "public",
		},
	}

	shareLinkBytes, err := json.Marshal(shareLinkBody)
	if err != nil {
		return ShareLinkResp{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", createSharedLink, bytes.NewBuffer(shareLinkBytes))
	if err != nil {
		return ShareLinkResp{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ShareLinkResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return ShareLinkResp{}, fmt.Errorf("dropbox error: %s", string(b))
	}

	var shareLinkResp ShareLinkResp
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&shareLinkResp)
	if err != nil {
		return ShareLinkResp{}, err
	}

	return shareLinkResp, nil

}

type ShareLinkResp struct {
	Tag            string    `json:".tag"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PathLower      string    `json:"path_lower"`
	Rev            string    `json:"rev"`
	ServerModified time.Time `json:"server_modified"`
	Size           int       `json:"size"`
	URL            string    `json:"url"`
}
