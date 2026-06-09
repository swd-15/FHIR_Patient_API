package fhir

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//FHIRサーバーへのHTTPクライアント
type FHIRClient struct {
	BaseURL string
}

//FHIRClientを生成
func NewFHIRClient(baseURL string) *FHIRClient {
	return &FHIRClient{BaseURL: baseURL}
}

//指定患者IDの全リソースを取得
func (c *FHIRClient) FetchPatientEverything(patientID string) (*Bundle, error) {
	url := fmt.Sprintf("%s/Patient/%s/$everything", c.BaseURL, patientID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var bundle Bundle
	if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return &bundle, nil
}

//患者一覧を取得
func (c *FHIRClient) FetchPatients(count int) (*Bundle, error) {
	url := fmt.Sprintf("%s/Patient?_count=%d", c.BaseURL, count)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	var bundle Bundle
	if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return &bundle, nil
}

//複数患者のBundleを取得して1つにまとめる
func (c *FHIRClient) FetchMultiplePatientEverything(count int) (*Bundle, error) {
	// まず患者一覧を取得
	patientBundle, err := c.FetchPatients(count)
	if err != nil {
		return nil, err
	}

	// 全患者のデータをまとめるBundleを作成
	mergedBundle := &Bundle{
		ResourceType: "Bundle",
		Type:         "collection",
	}

	// 各患者のEverythingを取得してmergedBundleに追加
	for _, entry := range patientBundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Patient" {
			continue
		}
		everything, err := c.FetchPatientEverything(r.ID)
		if err != nil {
			continue // エラーの患者はスキップ
		}
		mergedBundle.Entry = append(mergedBundle.Entry, everything.Entry...)
	}

	return mergedBundle, nil
}
