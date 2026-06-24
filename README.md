# FHIR Patient API

HL7 FHIR R4に準拠したバックエンドAPIです。
日本の医療DX推進で定められた6情報（傷病名・アレルギー・感染症・薬剤禁忌・処方情報・検査情報）をREST APIとして提供します。

> FHIR Bundle形式の医療データを読み込み、電子カルテ情報共有サービスで定められた6情報をREST APIとして提供するGo/Gin製のバックエンドデモアプリです。

---

## 主な機能

| 機能 | 詳細 |
|------|------|
| **FHIR Bundle解析** | HL7 FHIR R4形式のJSONを読み込みパース |
| **患者情報API** | 患者一覧・詳細をREST APIで提供 |
| **傷病名API** | 患者に紐づくConditionリソースを提供 |
| **検査情報API** | 患者に紐づくObservationリソースを提供 |
| **アレルギー情報API** | 患者に紐づくAllergyIntoleranceリソースを提供 |
| **処方情報API** | 患者に紐づくMedicationRequestリソースを提供 |
| **感染症情報API** | LOINCコードで絞り込んだ感染症検査結果を提供 |
| **薬剤禁忌API** | category=medicationのAllergyIntoleranceを提供 |
| **JP Core対応** | 日本標準コードシステム（ICD-10・JLAC10）に対応 |
| **FHIRサーバー接続** | HAPI FHIRサーバーからリアルタイムでデータ取得 |

---

## データ構造

```go
type Resource struct {
    ResourceType              string           // Patient / Condition / Observation / AllergyIntolerance / MedicationRequest
    ID                        string           // 患者ID
    Gender                    string           // 性別
    BirthDate                 string           // 生年月日
    Subject                   *Reference       // 患者への参照
    Code                      *CodeableConcept // 疾患・検査コード
    ClinicalStatus            *CodeableConcept // 臨床ステータス（active など）
    OnsetDateTime             string           // 発症日
    ValueQuantity             *Quantity        // 検査値（数値・単位）
    ValueCodeableConcept      *CodeableConcept // 検査結果（陽性/陰性など）
    Status                    string           // ステータス
    Patient                   *Reference       // 患者への参照（AllergyIntolerance）
    Criticality               string           // 重篤度（AllergyIntolerance）
    AllergyCategory           []string         // アレルギーカテゴリ（medication など）
    MedicationCodeableConcept *CodeableConcept // 薬剤コード（MedicationRequest）
    DosageInstruction         []struct{ Text string } // 用法用量
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
│   │   ├── parser.go
│   │   └── client.go
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

### ファイルモードで起動（デフォルト）

```bash
go run main.go
```

### FHIRサーバーモードで起動

```bash
FHIR_MODE=server FHIR_PATIENT_COUNT=5 go run main.go
```

### テスト実行

```bash
go test -v ./tests/...
```

---

## API エンドポイント一覧

| メソッド | パス | 説明 |
|---|---|---|
| GET | /health | ヘルスチェック |
| GET | /api/v1/patients | 患者一覧取得 |
| GET | /api/v1/patients/:id | 患者詳細取得 |
| GET | /api/v1/patients/:id/conditions | 傷病名取得 |
| GET | /api/v1/patients/:id/observations | 検査情報取得 |
| GET | /api/v1/patients/:id/allergies | アレルギー・薬剤禁忌取得 |
| GET | /api/v1/patients/:id/medications | 処方情報取得 |
| GET | /api/v1/patients/:id/infections | 感染症情報取得 |

---

## 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 患者一覧
curl http://localhost:8080/api/v1/patients

# 患者詳細
curl http://localhost:8080/api/v1/patients/p001

# 傷病名
curl http://localhost:8080/api/v1/patients/p001/conditions

# 検査情報
curl http://localhost:8080/api/v1/patients/p001/observations

# アレルギー・薬剤禁忌
curl http://localhost:8080/api/v1/patients/p001/allergies

# 処方情報
curl http://localhost:8080/api/v1/patients/p001/medications

# 感染症情報
curl http://localhost:8080/api/v1/patients/p001/infections
```

---

## アーキテクチャ

レイヤーを `fhir` / `service` / `handler` に分離し、それぞれ独立してテスト・拡張しやすい構成にしています。

```
main.go
  └─ service.NewPatientService()                ← bundle.json をロード
  └─ service.NewPatientServiceFromFHIRMultiple() ← FHIRサーバーから取得
        └─ fhir.LoadBundle() / FHIRClient        ← JSONパース・HTTP取得
handler.PatientHandler                           ← HTTPリクエストを受け取る
  └─ service.PatientService                      ← ビジネスロジック
        └─ fhir.Extract*/Find*()                 ← Bundleからリソース抽出
```

---

## 扱うFHIRリソース

| リソース | 6情報との対応 | 用途 |
|---|---|---|
| Patient | - | 患者基本情報（氏名・性別・生年月日） |
| Condition | 傷病名 | 診断・疾患情報（病名・臨床ステータス・発症日） |
| Observation | 検査情報・感染症情報 | 検査値・バイタル・感染症検査結果 |
| AllergyIntolerance | アレルギー・薬剤禁忌 | アレルギー情報・投与禁忌薬剤 |
| MedicationRequest | 処方情報 | 処方薬剤名・用法用量 |

---

## ライセンス

[MIT License](./LICENSE)
