package jdmapping

type MappingConfig struct {
	MappingNode *MappingNode
}

type ValueMap map[string]string

type MappingChildren map[string]*MappingNode

// 映射信息节点
type MappingNode struct {
	Text         string          `json:"text"`        //  字段映射文字
	ConvertKey   string          `json:"convert_key"` // array转map的字段
	ValueMap     ValueMap        `json:"value_map"`   // 值映射map
	Children     MappingChildren `json:"children"`    // 子字段数组
	ExcludeKeys  []string
	TextFn       func()                                                      `json:"-"` //  字段映射文字函数
	ConvertKeyFn func(value any) string                                      `json:"-"` // array转map的字段函数
	IsAddOPFn    func(oldValue, newValue any, mappingNode *MappingNode) bool `json:"-"` // 判断是否是添加操作的字段

	Key    string       `json:"key"` // 字段
	Parent *MappingNode `json:"-"`   // 父节点

	IsArray     bool `json:"is_array"`      // 字段是否是数组
	IsArrayItem bool `json:"is_array_item"` // 是否是数组的子元素
	IsRootNode  bool `json:"is_root_node"`  // 是否是根节点
	IsLeafNode  bool `json:"is_leaf_node"`  //是否是叶子节点
}

type StoreNode struct {
	KeyList        []string
	ArrayIndexList []int
	MappingNode    *MappingNode `json:"-"`
}

type JsonValue struct {
	Path     string
	newValue any
	oldValue any
}

type JsonValueV2 struct {
	Path  string
	Value any
}
