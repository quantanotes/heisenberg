package core

// import (
// 	"github.com/quantanotes/heisenberg/hnsw"
// 	"github.com/quantanotes/heisenberg/utils"
// )

// type indexConfig struct {
// 	Name  string
// 	Count uint
// 	Dim   uint
// 	Space uint
// }

// type Indexa struct {
// 	config indexConfig
// 	hnsw   *hnsw.HNSW
// }

// func NewIndex(conf indexConfig) *Index {
// 	hnsw := hnsw.NewHNSW(utils.SpaceType(conf.Space), int(conf.Dim), 1000000, &hnsw.HNSWOptions{M: 50, Ef: 100}, 69420)
// 	return &Index{conf, hnsw}
// }

// func LoadIndex(path string, conf indexConfig) (*Index, error) {
// 	hnsw, err := hnsw.LoadHNSW(path, int(conf.Dim), utils.SpaceType(conf.Space))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Index{
// 		conf,
// 		hnsw,
// 	}, nil
// }

// func (i *Index) Close() {
// 	i.hnsw.Free()
// }

// func (i *Index) Save(path string) {
// 	i.hnsw.Save(path)
// }

// func (i *Index) Insert(id uint, vec []float32) error {
// 	if len(vec) != int(i.config.Dim) {
// 		return utils.DimensionMismatch(len(vec), int(i.config.Dim))
// 	}
// 	return i.hnsw.Add(int(id), vec)
// }

// func (i *Index) Delete(id uint) error {
// 	return i.hnsw.Delete(int(id))
// }

// func (i *Index) Search(query []float32, k int) ([]int, error) {
// 	return i.hnsw.Search(query, k)
// }

// // Generates unique index keys
// func (i *Index) Next() uint {
// 	i.config.Count++
// 	return i.config.Count
// }
