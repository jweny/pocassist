import React from "react";
import { ColumnProps } from "antd/es/table";
import { FormColumnProps } from "./SearchForm";
import { RuleDataProps } from "../../../api/rule";

export const getAllColumns: () => ColumnProps<RuleDataProps>[] = () => {
  let allColumns: ColumnProps<RuleDataProps>[] = [
    { title: "漏洞编号", dataIndex: "vul_id", ellipsis: false },
    {
      title: "状态",
      dataIndex: "enable",
      ellipsis: true,
      render: (value: boolean) => (value ? "启用" : "禁用")
    },
    {
      title: "poc名称",
      dataIndex: "json_poc",
      ellipsis: false,
      render: (value: any) => value.name
    },
    { title: "漏洞名", dataIndex: "desp_name", ellipsis: true },
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
];

