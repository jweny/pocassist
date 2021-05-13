import React, { useContext, useReducer } from "react";
import { Button, Form, Input, Select } from "antd";
import { FormItemProps } from "antd/es/form";
import VulContext from "../../../store/vul/store";

export type FormColumnProps = Omit<FormItemProps, "children"> & {
  render?: () => React.ReactNode;
};

export interface VulComponentProps {}

const VulSearchFrom: React.FC<VulComponentProps> = props => {
  const [form] = Form.useForm();
  const { state, dispatch } = useContext(VulContext);
  // const { moduleList, productList } = props;

  const formColumns: FormColumnProps[] = [
    {
      name: "typeField",
      label: "漏洞类型",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 200 }} allowClear>
            {state.basic?.VulType.map(item => {
              return (
                <Select.Option value={item.name} key={item.name}>
                  {item.label}
                </Select.Option>
              );
            })}
          </Select>
        );
      }
    },
    {
      name: "productField",
      label: "组件",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 120 }} allowClear>
            {state.productList?.map(item => {
              return (
                <Select.Option value={item.id as number} key={item.id}>
                  {item.name}
                </Select.Option>
              );
            })}
          </Select>
        );
      }
    },
    {
      name: "search",
      label: "模糊查询"
    }
  ];

  const handleFinish = (values: any) => {
    dispatch({ type: "SET_SEARCH_QUERY", payload: values });
  };

  return (
    <Form
      form={form}
      layout="inline"
      className="vul-manage-form"
      onFinish={handleFinish}
    >
      {formColumns.map(
        (
          {
            name,
            label,
            render = () => <Input placeholder={`请输入${label}`} allowClear />,
            ...formProps
          },
          index
        ) => (
          <Form.Item label={label} key={index} name={name} {...formProps}>
            {render()}
          </Form.Item>
        )
      )}
      <Button type="primary" htmlType="submit">
        查询
      </Button>
    </Form>
  );
};

export default VulSearchFrom;
