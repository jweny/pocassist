import React from "react";
import { ColumnProps } from "antd/es/table";
import { VulDataProps } from "../../../api/vul";
import { FormColumnProps } from "./SearchForm";
import { RuleDataProps } from "../../../api/rule";

export const getAllColumns: () => ColumnProps<RuleDataProps>[] = () => {
  let allColumns: ColumnProps<RuleDataProps>[] = [
    { title: "ID", dataIndex: "id", ellipsis: true },
    { title: "漏洞编号", dataIndex: "vul_id", ellipsis: true },
    {
      title: "是否启用",
      dataIndex: "enable",
      ellipsis: true,
      render: (value: boolean) => (value ? "是" : "否")
    },
    {
      title: "json_poc名称",
      dataIndex: "json_poc",
      ellipsis: true,
      render: (value: any) => value.name
    },
    { title: "漏洞描述", dataIndex: "desp_name", ellipsis: true },
    {
      title: "影响类型",
      dataIndex: "affects",
      ellipsis: true
    }
  ];
  return allColumns;
};

export const richFormColumns: FormColumnProps[] = [
  {
    name: "description",
    label: "描述"
  },
  {
    name: "vulnerability",
    label: "漏洞危害"
  },
  {
    name: "affected_version",
    label: "影响版本"
  },
  {
    name: "suggestion",
    label: "修复建议"
  },
  {
    name: "exploit",
    label: "Exploit"
  }
  // {
  //   name: "remarks",
  //   label: "备注"
  // }
];

export const criteriaOperator = ["AND", "OR", "XOR", "ONE"];
export const criterionOperator = [
  "equal",
  "not equal",
  "contains",
  "not contains",
  "greater than",
  "less than",
  "pattern match"
];
