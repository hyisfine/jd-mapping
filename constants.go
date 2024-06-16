package jdmapping

// 比对类型
type DiffOperator string

const (
	DiffOperatorAdd     DiffOperator = "add"
	DiffOperatorRemove  DiffOperator = "remove"
	DiffOperatorReplace DiffOperator = "replace"
)

const (
	ARRAY_NODE_INDEX_REPLACE_KEY = "ARRAY_NODE_INDEX_REPLACE_KEY" // 数组下标代替字段
	// RootNodeReplaceKey       = "ROOT_NODE_INDEX_REPLACE_KEY"  // 根节点代替字段
)

const (
	SeqSlash = "/"
	SeqPoint = "."
)
