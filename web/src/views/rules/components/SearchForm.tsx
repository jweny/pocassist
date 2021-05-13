import React, { useContext } from "react";
import { Button, Form, Input, Select } from "antd";
import { FormItemProps } from "antd/es/form";
import RuleContext from "../../../store/rule/store";

export type FormColumnProps = Omit<FormItemProps, "children"> & {
  render?: () => React.ReactNode;
};

export interface VulComponentProps {}

const VulSearchFrom: React.FC<VulComponentProps> = props => {
  const [form] = Form.useForm();
  const { state, dispatch } = useContext(RuleContext);

  const formColumns: FormColumnProps[] = [
    {
      name: "enableField",
      label: "是否启用",
      render: () => {
        return (
          <Select placeholder="状态" style={{ width: 120 }} allowClear>
            <Select.Option value="True">启用</Select.Option>
            <Select.Option value="False">禁用</Select.Option>
          </Select>
        );
      }
    },
    {
      name: "affectsField",
      label: "影响类型",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 200 }} allowClear>
            {state.basic?.ModuleAffects.map(item => {
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
