package cli

import (
	"net/http"
	"os"
)

func Download(url string, file string) error {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.ReadFrom(resp.Body)
	return err
}
