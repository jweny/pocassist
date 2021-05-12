import { FormColumnProps } from "../../vul/components/SearchForm";
import { Input } from "antd";
import React from "react";

export const allModuleColumns = [
  { title: "模块名称", dataIndex: "name", width: "20%" },
  { title: "影响类型", dataIndex: "affects", width: "20%" },
  { title: "备注", dataIndex: "remarks", width: "30%" }
];

export const allProductColumns = [
  { title: "组件名称", dataIndex: "name", width: "20%" },
  { title: "厂商", dataIndex: "provider", width: "20%" },
  { title: "备注", dataIndex: "remarks", width: "30%" }
];

export const moduleFormColumns: FormColumnProps[] = [
  {
    name: "name",
    label: "模块名称",
    rules: [{ required: true }]
  },
  {
    name: "affects",
    label: "影响类型"
  },
  {
    name: "remarks",
    label: "备注",
    render: () => {
      return <Input.TextArea />;
    }
  }
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
