package jdmapping

import (
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
)

func TestMappingInitNode(t *testing.T) {

	config := &MappingConfig{
		MappingNode: &MappingNode{
			Text: "报警信息",
			Children: MappingChildren{
				"alarm_rule_list": {
					Text:       "报警规则列表",
					ConvertKey: "id",
					Children: MappingChildren{
						ARRAY_NODE_INDEX_REPLACE_KEY: {
							ExcludeKeys: []string{"value_type", "status"},
							Children: MappingChildren{
								"alarm_type": {
									Text: "报警类型",
									ValueMap: ValueMap{
										"1": "持续低迷",
										"2": "大幅下滑",
										"3": "连续下滑",
										"4": "危险预警",
										"5": "风险预警",
										"6": "大幅波动",
									},
								},
								"value": {
									Text: "值",
								},
								"time": {
									Text: "次数",
								},
								"send_ids": {
									Text: "发送群id",
								},
								"display_days": {
									Text: "展示天数",
								},
								"send_user": {
									Text:       "发送人",
									ConvertKey: "id",
									Children: MappingChildren{
										ARRAY_NODE_INDEX_REPLACE_KEY: {
											// IsAddOPFn: func(oldValue, newValue any, mappingNode *MappingNode) bool {
											// 	m, ok := newValue.(map[string]any)
											// 	if !ok {
											// 		return ok
											// 	}

											// 	return cast.ToInt(m["id"]) == 0
											// },
											Children: MappingChildren{
												"name": {
													Text: "姓名",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"warning_rule": {
					Text: "预警规则",
					IsAddOPFn: func(oldValue, newValue any, mappingNode *MappingNode) bool {
						m, ok := newValue.(map[string]any)
						if !ok {
							return ok
						}

						return cast.ToInt(m["id"]) == 0
					},
					Children: MappingChildren{
						"period": {
							Text: "报警周期",
							ValueMap: ValueMap{
								"1": "日",
								"2": "周",
								"3": "半月",
								"4": "月",
								"5": "季",
							},
						},
						"type": {
							Text: "预警规则",
							ValueMap: ValueMap{
								"1": "达成率-时间进度",
								"2": "指标达成率",
								"3": "目标值-实际值",
								"4": "达成值波动",
							},
						},
						"index_time_period": {
							Text: "指标计算周期",
							ValueMap: ValueMap{
								"1": "日",
								"2": "周",
								"3": "半月",
								"4": "月",
								"5": "季",
							},
						},
						"warning_thresholds": {
							Text:       "预警阈值",
							ConvertKey: "id",
							Children: MappingChildren{
								ARRAY_NODE_INDEX_REPLACE_KEY: {
									ExcludeKeys: []string{"value_type", "status"},
									IsAddOPFn: func(oldValue, newValue any, mappingNode *MappingNode) bool {
										m, ok := newValue.(map[string]any)
										if !ok {
											return ok
										}

										return cast.ToInt(m["id"]) == 0
									},
									Children: MappingChildren{
										"threshold_level": {
											Text: "颜色等级",
											ValueMap: ValueMap{
												"1": "绿色",
												"2": "黄色",
												"3": "红色",
											},
										},
										"threshold_operator": {
											Text: "运算符",
											ValueMap: ValueMap{
												"1":  "小于",
												"2":  "小于等于",
												"3":  "等于",
												"4":  "大于",
												"5":  "大于等于",
												"6":  "小于且大于",
												"7":  "小于等于且大于",
												"8":  "小于且大于等于",
												"9":  "小于等于且大于等于",
												"10": "达成",
												"11": "未达成",
											},
										},
										"top_value": {
											Text: "阈值上限",
										},
										"bottom_value": {
											Text: "阈值下限",
										},
									},
								},
							},
						},
					},
				},
				"effective_date": {
					Text: "生效日期",
				},
			},
		},
	}

	data1 := `{
    "alarm_rule_list": [
        {
            "alarm_type": 4,
            "value": 0,
            "value_type": 1,
            "time": 0,
            "send_method": 1,
            "send_ids": "oc_65b7404b6f89c030439cb4824b6c7ed2,oc_3d19d74656d76c591f1edfaca992e3db",
            "display_days": 3,
            "id": 1728,
            "status": 1,
			"send_user": [{"id":1,"name":"ypc"},{"id":0,"name":"ypc"}]
        }
    ],
    "warning_rule": {
        "period": 2,
        "type": 1,
        "index_time_period": 3,
        "id": 0,
        "warning_thresholds": [
            {
                "threshold_level": 3,
                "threshold_operator": 2,
                "top_value": -0.05,
                "bottom_value": 0,
                "value_type": 1,
                "id":0
            },
            {
                "threshold_level": 2,
                "threshold_operator": 6,
                "top_value": 0,
                "bottom_value": -0.05,
                "value_type": 1,
                "id": 4947
            },
            {
                "threshold_level": 1,
                "threshold_operator": 5,
                "top_value": 0,
                "bottom_value": 0,
                "value_type": 1,
                "id": 4948
            },
            {
                "threshold_level": 1,
                "threshold_operator": 5,
                "top_value": 0,
                "bottom_value": 0,
                "value_type": 1,
                "id": 4948
            }
        ]
    },
    "effective_date": "2024/05/27"
}`

	data2 := `{
    "alarm_rule_list": [
        {
            "alarm_type": 4,
            "value": 0,
            "value_type": 1,
            "time": 0,
            "send_method": 1,
            "send_ids": "oc_65b7404b6f89c030439cb4824b6c7ed2,oc_3d19d74656d76c591f1edfaca992e3db",
            "display_days": 3,
            "id": 1728,
            "status": 1,
			"send_user": [{"id":1,"name":"ypc"}]
        }
        {
            "alarm_type": 4,
            "value": 0,
            "value_type": 1,
            "time": 0,
            "send_method": 1,
            "send_ids": "oc_65b7404b6f89c030439cb4824b6c7ed2,oc_3d19d74656d76c591f1edfaca992e3db",
            "display_days": 3,
            "id": 0,
            "status": 1,
			"send_user": [{"id":1,"name":"ypc"}]
        }
    ],
    "warning_rule": {
        "period": 2,
        "type": 1,
        "index_time_period": 3,
        "id": 111,
        "warning_thresholds": [
            {
                "threshold_level": 3,
                "threshold_operator": 2,
                "top_value": -0.05,
                "bottom_value": 0,
                "value_type": 1,
                "id": 0
            },
            {
                "threshold_level": 1,
                "threshold_operator": 5,
                "top_value": 0,
                "bottom_value": 0,
                "value_type": 1,
                "id": 4948
            }
        ]
    },
    "effective_date": "2024/05/27"
}`

	jd := New(config)
	j, err := NewJson(data1, data2, jd)
	if err != nil {
		panic(err)
	}

	var a, b any
	sonic.UnmarshalString(j.oldJson, &a)
	sonic.UnmarshalString(j.newJson, &b)

	fmt.Println(sonic.MarshalString(a))
	fmt.Println(sonic.MarshalString(b))

}
