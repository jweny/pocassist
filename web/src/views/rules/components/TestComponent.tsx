import React, { useEffect, useState } from "react";
import {
  Button,
  Input,
  Modal,
  Popconfirm,
  Space,
  Tabs
} from "antd";
import { getId } from "./RunTest";
import RuleComponent from "./RuleComponent";

export interface TestComponentProps {
  data: any;
  type: string;
  setData: (value: TestComponentProps["data"]) => void;
  delete: (id: string) => void;
}
interface RequestProps {
  cookies?: string;
  custom_headers?: string;
  method?: string;
  post_text?: string;
  url?: string;
  version?: string;
  name?: string;
  data?: any;
}

const TestComponent: React.FC<TestComponentProps> = props => {
  const { data } = props;

  const handleRequestChange = (
    id: number,
    type: keyof RequestProps,
    val: any
  ) => {
    const current = {
      ...data,
      data: data.data.map((item: any) => {
        if (item.id === id) {
          return {
            ...item,
            [type]: val
          };
        } else {
          return item;
        }
      })
    };
    props.setData(current);
  };

  const handleDelete = () => {
    props.delete(data.id as string);
  };

  const onEdit = (targetKey: any, action: any) => {
    if (action === "add") {
      add();
    } else {
      Modal.confirm({
        title: "真的删除？",
        onOk() {
          remove(targetKey);
        }
      });
    }
  };

  const add = () => {
    const current = {
      ...data,
      data: [...data.data, { id: getId() }]
    };
    props.setData(current);
  };

  const remove = (targetKey: string) => {
    console.log(targetKey);
    const current = {
      ...data,
      data: data.data.filter((item: any) => item.id !== targetKey)
    };
    props.setData(current);
  };

  const handleNameChange = (val: string) => {
    const current = {
      ...data,
      name: val
    };
    props.setData(current);
  };

  return (
    <div className="test-component-wrap">
      {props.type === "groups" && (
        <Input
          value={data.name}
          style={{ width: "200px", margin: "10px 0" }}
          onChange={event => handleNameChange(event.target.value)}
        />
      )}
      <Tabs
        type="editable-card"
        defaultActiveKey="1"
        onEdit={onEdit}
        tabBarExtraContent={
          <Space>
            {props.type === "groups" && (
              <Popconfirm title="真的删除？" onConfirm={handleDelete}>
                <Button type="link" danger>
                  删除rules
                </Button>
              </Popconfirm>
            )}
          </Space>
        }
      >
        {data.data.map((rule: any, index: number) => {
          return (
            <Tabs.TabPane tab={`rule`} key={rule.id}>
              <RuleComponent data={rule} setData={handleRequestChange} />
            </Tabs.TabPane>
          );
        })}
      </Tabs>
    </div>
  );
};

export default TestComponent;
