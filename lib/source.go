package ecslogs

import (
	"os"
	"sort"
	"sync"
)

type Source interface {
	Open() (Reader, error)
}

type SourceFunc func() (Reader, error)

func (f SourceFunc) Open() (Reader, error) {
	return f()
}

func RegisterSource(name string, source Source) {
	srcmtx.Lock()
	srcmap[name] = source
	srcmtx.Unlock()
}

func DeregisterSource(name string) {
	srcmtx.Lock()
	delete(srcmap, name)
	srcmtx.Unlock()
}

func GetSource(name string) (source Source) {
	srcmtx.RLock()
	source = srcmap[name]
	srcmtx.RUnlock()
	return
}

func GetSources(names ...string) (sources []Source) {
	if len(names) != 0 {
		sources = make([]Source, 0, len(names))

		for _, name := range names {
			if source := GetSource(name); source != nil {
				sources = append(sources, source)
			}
		}
	}
	return
}

func SourcesAvailable() (sources []string) {
	srcmtx.RLock()
	sources = make([]string, 0, len(srcmap))

	for name := range srcmap {
		sources = append(sources, name)
	}

	srcmtx.RUnlock()
	sort.Strings(sources)
	return
}

var (
	srcmtx sync.RWMutex
	srcmap = map[string]Source{
		"stdin": SourceFunc(func() (Reader, error) {
			return NewMessageDecoder(os.Stdin), nil
		}),
	}
)
