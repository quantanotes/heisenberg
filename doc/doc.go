package doc

import (
	"heisenberg/common"
	"heisenberg/utils"
)

type keyIndex struct {
	nIndex numIndex
	sIndex strIndex
}

type Doc struct {
	index map[string]keyIndex
}

func (d *Doc) Insert(id uint64, doc common.Meta, old common.Meta) {
	flat := flattenDoc(doc)
	flatOld := flattenDoc(old)

	keys := utils.KeysAsArray(flat)
	oldKeys := utils.KeysAsArray(flat)

	deleteKeys := utils.Except(oldKeys, keys)
	for _, k := range deleteKeys {
		switch value := flatOld[k].(type) {
		case float64:
			d.index[k].nIndex.insert(id, value)
		case string:
			d.index[k].sIndex.remove(id, value)
		}
	}

	for k, v := range flat {
		if _, ok := d.index[k]; !ok {
			d.index[k] = keyIndex{
				nIndex: newMemNumIndex(),
				sIndex: newMemStrIndex(),
			}
		}
		switch value := v.(type) {
		case float64:
			d.index[k].nIndex.insert(id, value)
		case string:
			d.index[k].sIndex.insert(id, value, flatOld[k].(string))
		}
	}
}

func (d *Doc) Remove(id uint64, doc common.Meta) {
	flat := flattenDoc(doc)
	for k, v := range flat {
		switch value := v.(type) {
		case float64:
			d.index[k].nIndex.remove(id)
		case string:
			d.index[k].sIndex.remove(id, value)
		}
	}
}

func (d *Doc) Search() {

}

func flattenDoc(doc common.Meta) common.Meta {
	flattened := make(common.Meta)
	var flatten func(d common.Meta, prefix key)
	flatten = func(d common.Meta, prefix key) {
		for k, v := range d {
			currKey := append(prefix, k)
			switch value := v.(type) {
			case common.Meta:
				flatten(value, currKey)
			default:
				flattened[currKey.Compress()] = value
			}
		}
	}
	flatten(doc, key{})
	return flattened
}

func unflattenDoc(flat common.Meta) common.Meta {
	unflattened := make(common.Meta)
	for k, v := range flat {
		keys := UncompressKey(k)
		current := unflattened
		for i := 0; i < len(keys)-1; i++ {
			key := keys[i]
			if _, ok := current[key]; !ok {
				current[key] = make(common.Meta)
			}
			current = current[key].(common.Meta)
		}
		current[keys[len(keys)-1]] = v
	}
	return unflattened
}
