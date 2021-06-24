package mega

import (
  "io"
  "os"
  "fmt"
)

func NewMega(url, file string) *Mega {
  return &Mega{url, file, nil, nil, nil, "", nil, false}
}

func (m *Mega) Download() {
  if m.done {
    return
  }
  m.err = nil
  
  var f *os.File
  var r io.Reader
  var c io.Closer
  var link storageLinkResponse

  for a := 0; m.err == nil; a++ {
    switch a {
    case 0:
      if m.id != nil && m.key != nil && m.iv != nil {
        continue
      }
      m.id, m.key, m.iv, m.err = parseUrl(m.url)
    case 1:
      if m.link != "" {
        continue
      }
      if link, m.err = requestStorageLink(m.key, m.id); m.err == nil {
        m.link = link.G
        if m.name, m.err = parseNodeName([]byte(link.At), m.key, m.iv); m.err != nil || m.name == "" {
          m.name, m.err = string(m.id), nil
        }
      }
    case 2:
      if _, err := os.Stat(m.name); err == nil {
        m.err = fmt.Errorf("File %s already exist", m.name)
        continue
      }
      f, m.err = os.OpenFile(m.name, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
      defer f.Close()
    case 3:
      r, c, m.err = requestStorageReader(m.key, m.iv, m.link)
    case 4:
      defer c.Close()
    case 5:
      _, m.err = io.Copy(f, r)
    default:
      m.done = true
      return
    }
  }
}

func (m *Mega) Error() error {
  return m.err
}

func (m *Mega) Name() string {
  return m.name
}
