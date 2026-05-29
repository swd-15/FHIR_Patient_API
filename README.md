# FHIR Patient API

FHIR Patient API is a backend demo application built with Go and Gin.
It reads healthcare data in HL7 FHIR Bundle format and provides patient, condition, and observation data through REST APIs.

> FHIR Bundle形式の医療データを読み込み、患者情報・疾患・検査値をREST APIとして提供するGo/Gin製のバックエンドデモアプリです。

---

## 主な機能

| 機能 | 詳細 |
|------|------|
| **FHIR Bundle解析** | HL7 FHIR R4形式のJSONを読み込みパース |
| **患者情報API** | 患者一覧・詳細をREST APIで提供 |
| **疾患情報API** | 患者に紐づくConditionリソースを提供 |
| **検査値API** | 患者に紐づくObservationリソースを提供 |

---

## データ構造

```go
type Resource struct {
    ResourceType   string           // Patient / Condition / Observation
    ID             string           // 患者ID
    Gender         string           // 性別
    BirthDate      string           // 生年月日
    Subject        *Reference       // 患者への参照（Condition/Observation）
    Code           *CodeableConcept // 疾患・検査コード
    ClinicalStatus *CodeableConcept // 臨床ステータス（active など）
    OnsetDateTime  string           // 発症日
    ValueQuantity  *Quantity        // 検査値（数値・単位）
    Status         string           // Observationステータス
}
```

---

## ディレクトリ構成

```
FHIR_Patient_API/
├── main.go
├── go.mod
├── sample/
│   └── bundle.json
├── internal/
│   ├── fhir/
│   │   ├── model.go
│   │   └── parser.go
│   ├── service/
│   │   └── patient_service.go
│   └── handler/
│       └── patient_handler.go
└── tests/
    ├── parser_test.go
    ├── service_test.go
    └── handler_test.go
```

---

## セットアップ & 起動

### 必要環境

- Go 1.22 以上

### インストール

```bash
git clone https://github.com/swd-15/FHIR_Patient_API.git
cd FHIR_Patient_API
go mod tidy
```

### サーバー起動

```bash
go run main.go
```

### テスト実行

```bash
go test ./tests/...
```

---

## API エンドポイント一覧

| メソッド | パス | 説明 |
|---|---|---|
| GET | /health | ヘルスチェック |
| GET | /api/v1/patients | 患者一覧取得 |
| GET | /api/v1/patients/:id | 患者詳細取得 |
| GET | /api/v1/patients/:id/conditions | 疾患情報取得 |
| GET | /api/v1/patients/:id/observations | 検査値情報取得 |

---

## 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 患者一覧
curl http://localhost:8080/api/v1/patients

# 患者詳細
curl http://localhost:8080/api/v1/patients/p001

# 疾患情報
curl http://localhost:8080/api/v1/patients/p001/conditions

# 検査値
curl http://localhost:8080/api/v1/patients/p001/observations
```

---

## アーキテクチャ

レイヤーを `fhir` / `service` / `handler` に分離し、それぞれ独立してテスト・拡張しやすい構成にしています。

```
main.go
  └─ service.NewPatientService()   ← bundle.json をロード
        └─ fhir.LoadBundle()       ← JSONパース
handler.PatientHandler             ← HTTPリクエストを受け取る
  └─ service.PatientService        ← ビジネスロジック
        └─ fhir.Extract*/Find*()   ← Bundleからリソース抽出
```

---

## 扱うFHIRリソース

| リソース | 用途 |
|---|---|
| Patient | 患者基本情報（氏名・性別・生年月日） |
| Condition | 診断・疾患情報（病名・臨床ステータス・発症日） |
| Observation | 検査値・バイタル情報（HbA1cなど） |

---

## ライセンス

[MIT License](./LICENSE)
