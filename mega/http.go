package mega

import (
  "io"
  "fmt"
  "bytes"
  "net/http"
  "crypto/aes"
  "crypto/cipher"
  "encoding/json"
)

var nullIV []byte = bytes.Repeat([]byte{0}, 16)

func parseNodeName(at, key, iv []byte) (string, error) {
  block, err := aes.NewCipher(key)
  if err != nil {
    return "", err
  }

  dst, err := base64Dec(at)
  if err != nil {
    return "", err
  } else if len(dst) % block.BlockSize() != 0 {
    s := len(dst)
    bs := block.BlockSize()
    dst = dst[0:(bs * int(s / bs))]
  }
  
  mode := cipher.NewCBCDecrypter(block, nullIV)
  data := make([]byte, len(dst))
  mode.CryptBlocks(data, dst)

  var node *megaNode
  err = json.Unmarshal(bytes.Trim(data[4:], string([]byte{0})), &node)
  if err != nil {
    return "", err
  }
  return node.N, nil
}

func requestStorageLink(key, id []byte) (storageLinkResponse, error) {
  r, err := structToReader([]storageLinkRequest{storageLinkRequest{"g", 1, string(id)}})
  if err != nil {
    return storageLinkResponse{}, err
  }
  
  resp, err := http.Post(LINK_API, JSON_TYPE, r)
  if err != nil {
    return storageLinkResponse{}, err
  }
  defer resp.Body.Close()
  
  var link []storageLinkResponse
  err = readerToStruct(resp.Body, &link)
  if err != nil {
    return storageLinkResponse{}, err
  } else if len(link) < 1 {
    return storageLinkResponse{}, fmt.Errorf("Response of storage link don't have any object")
  }
  
  return link[0], nil
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
