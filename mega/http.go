package mega

import (
  "io"
  "errors"
  "net/http"
  "crypto/aes"
  "crypto/cipher"
)

func requestStorageLink(id string) (string, error) {
  r, err := structToReader([]storageLinkRequest{storageLinkRequest{"g", 1, id}})
  if err != nil {
    return "", err
  }
  r, err = structToReader([]storageLinkRequest{storageLinkRequest{"g", 1, id}})
  
  resp, err := http.Post(LINK_API, JSON_TYPE, r)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  
  var link []storageLinkResponse
  err = readerToStruct(resp.Body, &link)
  if err != nil {
    return "", err
  } else if len(link) < 1 {
    return "", errors.New("Response of storage link don't have any object")
  }
  
  return link[0].G, nil
}

func requestStorageReader(key, iv []byte, link string) (io.Reader, io.Closer, error) {
  resp, err := http.Get(link)
  if err != nil {
    return nil, nil, err
  }
  
  block, err := aes.NewCipher(key)
  if err != nil {
    defer resp.Body.Close()
    return nil, nil, err
  }
  
  return &cipher.StreamReader{cipher.NewCTR(block, iv),resp.Body}, resp.Body, nil
}
