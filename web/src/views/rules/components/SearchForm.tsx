import React, { useContext, useReducer } from "react";
import { Button, Form, Input, Select } from "antd";
import { FormItemProps } from "antd/es/form";
import RuleContext from "../../../store/rule/store";
import { ModuleDataProps } from "../../../api/module";
import { ProductDataProps } from "../../../api/product";
import { ScriptDataProps } from "../../../api/script";

export type FormColumnProps = Omit<FormItemProps, "children"> & {
  render?: () => React.ReactNode;
};

export interface VulComponentProps {}

const VulSearchFrom: React.FC<VulComponentProps> = props => {
  const [form] = Form.useForm();
  const { state, dispatch } = useContext(RuleContext);
  // const { moduleList, productList } = props;

  const formColumns: FormColumnProps[] = [
    // {
    //   name: "moduleField",
    //   label: "模块",
    //   render: () => {
    //     return (
    //       <Select placeholder="请选择" style={{ width: 200 }} allowClear>
    //         {state.moduleList?.map(item => {
    //           return (
    //             <Select.Option value={item.id as number} key={item.id}>
    //               {item.name}
    //             </Select.Option>
    //           );
    //         })}
    //       </Select>
    //     );
    //   }
    // },
    {
      name: "enableField",
      label: "是否启用",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 120 }} allowClear>
            <Select.Option value="True">是</Select.Option>
            <Select.Option value="False">否</Select.Option>
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
    // {
    //   name: "hasDesField",
    //   label: "关联",
    //   render: () => {
    //     return (
    //       <Select placeholder="请选择" style={{ width: 120 }} allowClear>
    //         <Select.Option value={1}>已关联</Select.Option>
    //         <Select.Option value={0}>未关联</Select.Option>
    //       </Select>
    //     );
    //   },
    //   initialValue: 1
    // },
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
