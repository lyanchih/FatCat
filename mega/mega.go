package mega

import (
  "io"
  "os"
)

func New(url, file string) *Mega {
  return &Mega{url, file, nil, nil, nil, "", nil}
}

func (m *Mega) Download() {
  var f *os.File
  var r io.Reader
  var c io.Closer

  for a := 0; m.err == nil; a++ {
    switch a {
    case 0:
      m.id, m.key, m.iv, m.err = parseUrl(m.Url)
    case 1:
      m.link, m.err = requestStorageLink(string(m.id))
    case 2:
      f, m.err = os.OpenFile(m.FileName, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
      defer f.Close()
    case 3:
      r, c, m.err = requestStorageReader(m.key, m.iv, m.link)
      defer c.Close()
    case 4:
      _, m.err = io.Copy(f, r)
    default: return
    }
  }
}

func (m *Mega) Error() error {
  return m.err
}
