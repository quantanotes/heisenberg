package query

type Plan []node

func NewPlan() Plan {
	return Plan{}
}

func (p Plan) Delete(key []byte) Plan {
	return append(p, newNode(NodeDelete, key))
}

func (p Plan) Seek() {

}
