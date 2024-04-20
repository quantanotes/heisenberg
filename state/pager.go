package state

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Pager struct {
	disk   Disk
	pages  []File
	prefix string
	suffix string
}

func NewPager(disk Disk, prefix string, suffix string) Pager {
	p := Pager{
		prefix: prefix,
		disk:   disk,
		suffix: suffix,
	}
	p.load()
	return p
}

func (p *Pager) New() (File, error) {
	f, err := p.disk.Open(p.nextPageName())
	if err != nil {
		return nil, err
	}
	p.pages = append([]File{f}, p.pages...)
	return f, nil
}

func (p *Pager) Get(index int) File {
	return p.pages[index]
}

func (p *Pager) Ensure() error {
	if len(p.pages) == 0 {
		_, err := p.New()
		return err
	}
	return nil
}

func (p *Pager) Head() File {
	return p.pages[p.HeadIndex()]
}

func (p *Pager) HeadIndex() int {
	return len(p.pages) - 1
}

// TODO: more robust page loading algorithm

func (p *Pager) load() error {
	files, err := p.disk.Scan()
	if err != nil {
		return err
	}
	maxPage := p.getMaxPage(files)
	return p.openPages(maxPage)
}

func (p *Pager) getMaxPage(files []string) uint32 {
	var maxPage uint32
	for _, file := range files {
		_, name := filepath.Split(file)
		pageStr := strings.TrimPrefix(strings.TrimSuffix(name, p.suffix), p.prefix)
		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			continue
		}
		if uint32(page) > maxPage {
			maxPage = uint32(page)
		}
	}
	return maxPage
}

func (p *Pager) openPages(maxPage uint32) error {
	p.pages = make([]File, maxPage)
	for i := uint32(0); i <= maxPage; i++ {
		name := fmt.Sprintf("%s%d%s", p.prefix, i, p.suffix)
		file, err := p.disk.Open(name)
		if err != nil {
			return err
		}
		p.pages[i] = file
	}
	return nil
}

func (p *Pager) nextPageName() string {
	return fmt.Sprintf("%s%d%s", p.prefix, len(p.pages), p.suffix)
}
