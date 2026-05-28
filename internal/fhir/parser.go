package fhir

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//JSONファイルからFHIR Bundleを読み込んでパースする
func LoadBundle(path string) (*Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("bundle file read error: %w", err)
	}
	return ParseBundle(data)
}

//JSONバイト列をBundleにパースする
func ParseBundle(data []byte) (*Bundle, error) {
	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("bundle parse error: %w", err)
	}
	return &bundle, nil
}

//BundleからPatientリソースのみを取り出す
func ExtractPatients(bundle *Bundle) []PatientSummary {
	var patients []PatientSummary
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Patient" {
			continue
		}
		patients = append(patients, PatientSummary{
			ID:        r.ID,
			FullName:  buildFullName(r),
			Gender:    r.Gender,
			BirthDate: r.BirthDate,
		})
	}
	return patients
}

//指定IDのPatientを返す
func FindPatient(bundle *Bundle, id string) (*PatientSummary, bool) {
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Patient" || r.ID != id {
			continue
		}
		p := &PatientSummary{
			ID:        r.ID,
			FullName:  buildFullName(r),
			Gender:    r.Gender,
			BirthDate: r.BirthDate,
		}
		return p, true
	}
	return nil, false
}

//姓名を結合して返す
func buildFullName(r Resource) string {
	if len(r.Name) == 0 {
		return ""
	}
	n := r.Name[0]
	given := strings.Join(n.Given, " ")
	if given != "" {
		return n.Family + " " + given
	}
	return n.Family
}
