package jdmapping

import (
	"strings"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
)

func isArrayItem(key string) bool {
	return key == ARRAY_NODE_INDEX_REPLACE_KEY
}

func joinSlash(pathList []string) string {
	return strings.Join(pathList, SeqSlash)
}

func joinPoint(pathList []string) string {
	return strings.Join(pathList, SeqPoint)

}

func marshalString(data any) (str string, err error) {
	str, err = cast.ToStringE(data)
	if err != nil {
		str, err = sonic.MarshalString(data)
		return
	}
	if cast.ToString(data) == "" {
		str = "{}"
	}

	return
}
