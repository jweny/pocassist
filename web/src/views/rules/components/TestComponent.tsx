import React, { useEffect, useImperativeHandle, useState } from "react";
import { ColumnProps } from "antd/es/table";
import {
  Button,
  Col,
  Dropdown,
  Form,
  Input,
  Menu,
  message,
  Modal,
  Popconfirm,
  Radio,
  Row,
  Select,
  Space,
  Table,
  Tabs
} from "antd";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons/lib";
import { getId } from "./RunTest";
import { debounce } from "lodash";
import { criteriaOperator, criterionOperator } from "./columns";
import RuleComponent from "./RuleComponent";

export interface TestComponentProps {
  // data: {
  //   request?: RequestProps;
  //   response?: { var: ResponseProps };
  //   criteria?: CriterionProps;
  //   id?: string;
  //   method?: string;
  //   path?: string;
  //   header?: any;
  //   body?: string;
  //   follow_redirects?: boolean;
  //   search?: string;
  //   expression?: string;
  // };
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
interface ResponseProps {
  "#text"?: string;
  "@description"?: string;
  "@name"?: string;
  "@source"?: string;
}
export interface CriterionProps {
  "@operator"?: string;
  "@value"?: string;
  "@variable"?: string;
  "@comment"?: string;
  criterion?: CriterionProps[];
  criteria?: CriterionProps[];
  children?: CriterionProps[];
  key?: string;
  id: string;
}
const TestComponent: React.FC<TestComponentProps> = props => {
  const { data } = props;
  console.log(data);
  const [test, setTest] = useState<CriterionProps>({} as CriterionProps);
  const [expand, setExpand] = useState<string[]>([]);

  useEffect(() => {
    setTest(data?.criteria as CriterionProps);
    setExpand(prev =>
      Array.from(new Set([...prev, data?.criteria?.id ?? "1"]))
    );
  }, [data]);

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

  const handleReset = () => {
    const current = {
      id: data.id,
      request: {},
      response: {
        var: {}
      }
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
        title: "确认删除该规则吗？",
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
            {/*自动保存吧，体验好一点*/}
            {/*<Button type="link" onClick={() => generateJson()}>*/}
            {/*  保存*/}
            {/*</Button>*/}
            {/*<Popconfirm title="确定重置吗？" onConfirm={handleReset}>*/}
            {/*  <Button type="link">重置</Button>*/}
            {/*</Popconfirm>*/}
            {props.type === "groups" && (
              <Popconfirm title="确定删除吗" onConfirm={handleDelete}>
                <Button type="link" danger>
                  删除
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
