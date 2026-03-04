package dropbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
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
