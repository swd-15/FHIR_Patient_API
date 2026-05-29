package fhir

//FHIR Bundleリソースのトップレベル構造体
type Bundle struct {
	ResourceType string  `json:"resourceType"`
	Type         string  `json:"type"`
	Entry        []Entry `json:"entry"`
}

//Bundle内の各リソースエントリ
type Entry struct {
	Resource Resource `json:"resource"`
}

//FHIRリソースの共通フィールドを持つ構造体
type Resource struct {
	ResourceType string `json:"resourceType"`

	// Patient
	ID   string `json:"id"`
	Name []struct {
		Family string   `json:"family"`
		Given  []string `json:"given"`
	} `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`

	// Condition / Observation
	Subject *Reference       `json:"subject"`
	Code    *CodeableConcept `json:"code"`
	Status  string           `json:"status"`

	// Condition
	ClinicalStatus *CodeableConcept `json:"clinicalStatus"`
	OnsetDateTime  string           `json:"onsetDateTime"`

	// Observation
	Category       []CodeableConcept `json:"category"`
	Effective      string            `json:"effectiveDateTime"`
	ValueQuantity  *Quantity         `json:"valueQuantity"`
}

//FHIR Reference型
type Reference struct {
	Reference string `json:"reference"`
}

//コード+テキストによる概念表現
type CodeableConcept struct {
	Coding []Coding `json:"coding"`
	Text   string   `json:"text"`
}

//特定のコードシステムにおけるコード
type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display"`
}

//患者一覧で返す簡略情報
type PatientSummary struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

//疾患情報のレスポンス
type ConditionResponse struct {
	PatientID      string `json:"patient_id"`
	Display        string `json:"display"`
	Code           string `json:"code"`
	ClinicalStatus string `json:"clinical_status"`
	OnsetDate      string `json:"onset_date"`
}

//検査値のレスポンス
type ObservationResponse struct {
	PatientID     string  `json:"patient_id"`
	Display       string  `json:"display"`
	Code          string  `json:"code"`
	Value         float64 `json:"value"`
	Unit          string  `json:"unit"`
	EffectiveDate string  `json:"effective_date"`
	Status        string  `json:"status"`
}

type Quantity struct {
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
	System string  `json:"system"`
	Code   string  `json:"code"`
}
