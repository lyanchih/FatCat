package mega

const (
  API_HOST = "https://eu.api.mega.co.nz"
  LINK_API = API_HOST + "/cs"
  JSON_TYPE = "application/json"
)

type Mega struct {
  Url string
  FileName string
  id, key, iv []byte
  link string
  err error
}

type storageLinkRequest struct {
  A string `json:"a"`
  G int `json:"g"`
  P string `json:"p"`
}

type storageLinkResponse struct {
  S int `json:"s"`
  At string `json:"at"`
  G string `json:"g"`
}
