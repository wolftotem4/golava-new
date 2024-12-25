package pkg

import (
	"fmt"
	"strings"
)

var nativePkg = map[string]struct{}{
	"archive":   {},
	"bufio":     {},
	"builtin":   {},
	"bytes":     {},
	"cmp":       {},
	"compress":  {},
	"container": {},
	"context":   {},
	"crypto":    {},
	"database":  {},
	"debug":     {},
	"embed":     {},
	"encoding":  {},
	"errors":    {},
	"expvar":    {},
	"flag":      {},
	"fmt":       {},
	"go":        {},
	"hash":      {},
	"html":      {},
	"image":     {},
	"index":     {},
	"io":        {},
	"iter":      {},
	"log":       {},
	"maps":      {},
	"math":      {},
	"mime":      {},
	"net":       {},
	"os":        {},
	"path":      {},
	"plugin":    {},
	"reflect":   {},
	"regexp":    {},
	"runtime":   {},
	"slices":    {},
	"sort":      {},
	"strconv":   {},
	"strings":   {},
	"structs":   {},
	"sync":      {},
	"syscall":   {},
	"testing":   {},
	"text":      {},
	"time":      {},
	"unicode":   {},
	"unsafe":    {},
}

func IsNativePackage(pkg string) bool {
	_, ok := nativePkg[pkg]
	if ok {
		return true
	}

	for p := range nativePkg {
		if strings.HasPrefix(pkg, fmt.Sprintf("%s/", p)) {
			return true
		}
	}

	return false
}
