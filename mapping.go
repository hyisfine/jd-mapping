package jdmapping

type JDMapping struct {
	config          *MappingConfig
	isAddOPList     []*StoreNode
	isArrayList     []*StoreNode
	isArrayItemList []*StoreNode
	excludeKeysList []*StoreNode
}

// 处理mapping节点，添加私有字段，parent、key等
func (jd *JDMapping) initNode() {
	jd.config.MappingNode.IsRootNode = true
	jd.dfsInitNode(nil, jd.config.MappingNode, "")
}

// 深度优先处理节点数据,添加私有字段，parent、key等
func (jd *JDMapping) dfsInitNode(parent, node *MappingNode, key string) {
	if node == nil {
		return
	}

	node.Key = key
	node.Parent = parent

	if isArrayItem(key) {
		if parent != nil {
			parent.IsArray = true
		}
		node.IsArrayItem = true
	}

	if len(node.Children) == 0 {
		node.IsLeafNode = true
	}

	for key, children := range node.Children {
		jd.dfsInitNode(node, children, key)
	}
}

// 处理收集isAddOPMap、isArrayMap、isArrayItemMap
func (jd *JDMapping) dfsStoreConfigMap(node *MappingNode, keyList []string, arrayIndexList []int) {
	if node == nil {
		return
	}

	keyList = append(keyList, node.Key)

	if node.IsArray {
		arrayIndexList = append(arrayIndexList, len(keyList)-1)
	}

	storeNode := &StoreNode{
		MappingNode:    node,
		KeyList:        keyList,
		ArrayIndexList: arrayIndexList,
	}

	if node.IsAddOPFn != nil {
		jd.isAddOPList = append(jd.isAddOPList, storeNode)
	}

	if node.IsArray {
		jd.isArrayList = append(jd.isArrayList, storeNode)
	}

	if node.IsArrayItem {
		jd.isArrayItemList = append(jd.isArrayItemList, storeNode)
	}

	if len(node.ExcludeKeys) != 0 {
		jd.excludeKeysList = append(jd.excludeKeysList, storeNode)
	}

	for _, item := range node.Children {
		newKeyList := make([]string, len(keyList))
		copy(newKeyList, keyList)
		jd.dfsStoreConfigMap(item, newKeyList, arrayIndexList)
	}
}
