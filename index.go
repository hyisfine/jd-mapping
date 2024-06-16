package jdmapping

func New(config *MappingConfig) *JDMapping {
	jd := &JDMapping{
		config:          config,
		isAddOPList:     make([]*StoreNode, 0),
		isArrayList:     make([]*StoreNode, 0),
		isArrayItemList: make([]*StoreNode, 0),
		excludeKeysList: make([]*StoreNode, 0),
	}

	jd.initNode()
	jd.dfsStoreConfigMap(config.MappingNode, nil, nil)

	return jd
}
