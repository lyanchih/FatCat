package mega

const (
  API_HOST = "https://eu.api.mega.co.nz"
  LINK_API = API_HOST + "/cs"
  JSON_TYPE = "application/json"
)

type MegaRoot struct {}

type MegaInfo struct {
  Url string `json:"url"`
  Name string `json:"name"`
  Part uint `json:"part"`
}

type Mega struct {
  url, filename string
  id, key, iv []byte
  link string
  err error
  done bool
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
