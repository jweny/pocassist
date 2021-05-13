import { FormColumnProps } from "../../vul/components/SearchForm";
import { Input } from "antd";
import React from "react";


export const allProductColumns = [
  { title: "组件名称", dataIndex: "name", width: "20%" },
  { title: "厂商", dataIndex: "provider", width: "20%" },
  { title: "备注", dataIndex: "remarks", width: "30%" }
];


export const productFormColumns: FormColumnProps[] = [
  {
    name: "name",
    label: "组件名称",
    rules: [{ required: true }]
  },
  {
    name: "provider",
    label: "厂商"
  },
  {
    name: "remarks",
    label: "备注",
    render: () => {
      return <Input.TextArea />;
    }
  }
];
