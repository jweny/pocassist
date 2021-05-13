import React from "react";
import { ColumnProps } from "antd/es/table";
import { VulDataProps } from "../../../api/vul";
import { FormColumnProps } from "./SearchForm";

export const getAllColumns: () => ColumnProps<VulDataProps>[] = () => {
  let allColumns: ColumnProps<VulDataProps>[] = [
    { title: "漏洞名称", dataIndex: "name_zh", ellipsis: false, width: "30%" },
    { title: "影响组件", dataIndex: "webapp_name", ellipsis: true },
    { title: "漏洞类型", dataIndex: "category", ellipsis: true },
    { title: "cve编号", dataIndex: "cve", ellipsis: true },
    { title: "cnnvd编号", dataIndex: "cnnvd", ellipsis: true },
    { title: "修复建议", dataIndex: "suggestion", ellipsis: true },
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
  ];
  return allColumns;
};

export const richFormColumns: FormColumnProps[] = [
  {
    name: "description",
    label: "描述"
  },
  {
    name: "suggestion",
    label: "修复建议"
  },
];

