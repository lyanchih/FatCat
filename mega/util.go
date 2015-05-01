package mega

import (
  "io"
  "io/ioutil"
  "bytes"
  "math/big"
  "encoding/json"
  "encoding/base64"
  "encoding/binary"
)

func xor(x, y int64) (z int64) {
  _x := big.NewInt(x)
  _y := big.NewInt(y)
  _z := big.NewInt(z)
  
  _z.Xor(_x, _y)
  
  return _z.Int64()
}

func bs2i64(bsArr ...[]byte) (is []int64, err error) {
  is = make([]int64, len(bsArr))
  for i, bs := range bsArr {
    if len(bs) < 8 {
      bs = append(bs, []byte{0,0,0,0,0,0,0,0}[0:8-len(bs)]...)
    }
    
    buf := bytes.NewBuffer(bs)
    err = binary.Read(buf, binary.BigEndian, &is[i])
    if err != nil {
      return nil, err
    }
  }
  
  return is, nil
}

func parseKey(ks []byte) (key []byte, iv []byte, err error) {
  is, err := bs2i64(ks[0:8], ks[16:24], ks[8:16], ks[24:32])
  if err != nil {
    return nil, nil, err
  }
  
  buf := &bytes.Buffer{}
  err = binary.Write(buf, binary.BigEndian, xor(is[0], is[1]))
  err = binary.Write(buf, binary.BigEndian, xor(is[2], is[3]))
  
  return buf.Bytes(), append(ks[16:24], []byte{0,0,0,0,0,0,0,0}...), err
}


func parseUrl(url string) (id, key, iv []byte, err error) {
  arr := bytes.Split([]byte(url), []byte("!"))
  
  id = arr[1]
  
  tmpKey := bytes.Replace(arr[2], []byte("-"), []byte("+"), -1)
  tmpKey = bytes.Replace(tmpKey, []byte("_"), []byte("/"), -1)
  tmpKey = bytes.Replace(tmpKey, []byte(","), []byte(""), -1)
  tmpKey = append(tmpKey, []byte("=")...)
  
  dst := make([]byte, base64.URLEncoding.DecodedLen(len(tmpKey)))
  _, err = base64.StdEncoding.Decode(dst, tmpKey)
  if err != nil {
    return nil, nil, nil, err
  }
  
  key, iv, err = parseKey(dst)
  return
}

func structToReader(i interface{}) (io.Reader, error) {
  bs, err := json.Marshal(i)
  if err != nil {
    return nil, err
  }
  
  return bytes.NewBuffer(bs), nil
}

func readerToStruct(r io.Reader, i interface{}) (error) {
  bs, err := ioutil.ReadAll(r)
  if err != nil {
    return err
  }
  return json.Unmarshal(bs, i)
}
