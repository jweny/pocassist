import React from "react";
import { ColumnProps } from "antd/es/table";
import { VulDataProps } from "../../../api/vul";
import { FormColumnProps } from "./SearchForm";

export const getAllColumns: () => ColumnProps<VulDataProps>[] = () => {
  let allColumns: ColumnProps<VulDataProps>[] = [
    { title: "漏洞名称", dataIndex: "name_zh", ellipsis: false, width: "30%" },
    // {
    //   title: "XML文件名称",
    //   dataIndex: "xml_name",
    //   ellipsis: false,
    //   width: "15%"
    // },
    { title: "影响组件", dataIndex: "webapp_name", ellipsis: true },
    // { title: "脚本名称", dataIndex: "script_name", ellipsis: true },
    { title: "漏洞类型", dataIndex: "category", ellipsis: true },
    // { title: "所属模块", dataIndex: "module_name", ellipsis: true },
    { title: "漏洞危害", dataIndex: "vulnerability", ellipsis: true },
    { title: "影响版本", dataIndex: "affected_version", ellipsis: true },
    { title: "漏洞描述", dataIndex: "description", ellipsis: true },
    { title: "cve编号", dataIndex: "cve", ellipsis: true },
    { title: "cnnvd编号", dataIndex: "cnnvd", ellipsis: true },
    { title: "exploit", dataIndex: "exploit", ellipsis: true },
    { title: "修复建议", dataIndex: "suggestion", ellipsis: true },
    { title: "创建时间", dataIndex: "created_at", ellipsis: true },
    { title: "更新时间", dataIndex: "updated_at", ellipsis: true },
    {
      title: "扫描类型",
      dataIndex: "verifiability",
      ellipsis: true,
      render: (value: string) => {
        switch (value) {
          case "utilizable":
            return "精准扫描-可利用";
          case "general":
            return "通用扫描";
          case "non-utilizable":
            return "精准扫描-不可利用";
          default:
            return "";
        }
      }
    },
    {
      title: "漏洞级别",
      dataIndex: "severity",
      ellipsis: true,
      render: (value: string) => {
        switch (value) {
          case "high":
            return "高危";
          case "middle":
            return "中危";
          case "low":
            return "低危";
          case "info":
            return "信息";
          default:
            return "";
        }
      }
    },
    { title: "漏洞语言", dataIndex: "language", ellipsis: true },
    { title: "上传者", dataIndex: "writer_name", ellipsis: true },
    { title: "ID", dataIndex: "id", ellipsis: true }
    // { title: "备注", dataIndex: "remarks", ellipsis: true }
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
