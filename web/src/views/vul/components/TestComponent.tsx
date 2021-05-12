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

export interface TestComponentProps {
  data: {
    request?: RequestProps;
    response?: { var: ResponseProps };
    criteria?: CriterionProps;
    id?: string;
  };
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
  const [request, setRequest] = useState<RequestProps>({} as RequestProps);
  const [response, setResponse] = useState<ResponseProps>({});
  const [test, setTest] = useState<CriterionProps>({} as CriterionProps);
  const [expand, setExpand] = useState<string[]>([]);

  useEffect(() => {
    // 取出props的值
    // setRequest({ ...data?.request });
    // setResponse({ ...data?.response?.var });
    // const tempCriteria = data?.criteria ? { ...data?.criteria } : undefined;
    // console.log(tempCriteria);

    // function getCriteriaChildren(data: CriterionProps | undefined) {
    //   console.log(data);
    //
    //   if (data?.criteria) {
    //     const tempCriteria = Array.isArray(data.criteria)
    //       ? data.criteria
    //       : [data.criteria];
    //
    //     return tempCriteria.map(value => {
    //       console.log(value);
    //       const tempChildren: CriterionProps = {
    //         id: getId(),
    //         ...value,
    //         key: "criteria",
    //         children: getCriteriaChildren(value)
    //       };
    //       return tempChildren;
    //     });
    //   }
    //   return Array.isArray(data?.criterion)
    //     ? data?.criterion?.map(item => ({
    //         key: "criterion",
    //         id: getId(),
    //         ...item
    //       }))
    //     : data?.criterion
    //     ? [
    //         {
    //           key: "criterion",
    //           id: getId(),
    //           ...(data.criterion as CriterionProps)
    //         }
    //       ]
    //     : [];
    // }
    //
    // const curCriteria = {
    //   id: getId(),
    //   ...data?.criteria,
    //   key: "criteria",
    //   // criterion只有1条时，是一个单个对象，分情况判断赋给curCriteria变量
    //   children: getCriteriaChildren(tempCriteria)
    // };

    setTest(data?.criteria as CriterionProps);
    setExpand(prev =>
      Array.from(new Set([...prev, data?.criteria?.id ?? "1"]))
    );
  }, [data]);

  const columns: ColumnProps<CriterionProps>[] = [
    {
      title: "Test",
      dataIndex: "key",
      key: "key",
      width: "15%"
    },
    {
      title: "variable",
      dataIndex: "@variable",
      key: "variable",
      width: "20%",
      render: (value: CriterionProps["@variable"], record: CriterionProps) => {
        return (
          <Select
            value={value}
            style={{ width: "100%" }}
            disabled={record.key === "criteria"}
            onChange={val => handleTestChange("@variable", record.id, val)}
          >
            <Select.Option value="$(response_code)">
              $(response_code)
            </Select.Option>
            <Select.Option value="$(header)">$(header)</Select.Option>
            <Select.Option value="$(body)">$(body)</Select.Option>
            <Select.Option value="$(response_length)">
              $(response_length)
            </Select.Option>
          </Select>
        );
      }
    },
    {
      title: "operator",
      dataIndex: "@operator",
      width: "15%",
      key: "operator",
      render: (value, record) => {
        return (
          <Select
            value={value}
            style={{ width: "100%" }}
            onChange={val => handleTestChange("@operator", record.id, val)}
          >
            {/*<Select.Option value="AND">AND</Select.Option>*/}
            {/*<Select.Option value="OR">OR</Select.Option>*/}
            {/*<Select.Option value="contains">contains</Select.Option>*/}
            {/*<Select.Option value="not contains">not contains</Select.Option>*/}
            {/*<Select.Option value="pattern match">pattern match</Select.Option>*/}
            {record.key === "criteria"
              ? criteriaOperator.map(item => (
                  <Select.Option value={item} key={item}>
                    {item}
                  </Select.Option>
                ))
              : criterionOperator.map(item => (
                  <Select.Option value={item} key={item}>
                    {item}
                  </Select.Option>
                ))}
          </Select>
        );
      }
    },
    {
      title: "value",
      dataIndex: "@value",
      width: "20%",
      key: "value",
      render: (value, record) => {
        return (
          <Input
            value={value}
            disabled={record.key === "criteria"}
            onChange={e =>
              handleTestChange("@value", record.id, e.target.value)
            }
          />
        );
      }
    },
    {
      title: "comment",
      dataIndex: "@comment",
      width: "20%",
      key: "comment",
      render: (value, record) => {
        return (
          <Input
            value={value}
            disabled={record.key === "criteria"}
            onChange={e =>
              handleTestChange("@comment", record.id, e.target.value)
            }
          />
        );
      }
    },
    {
      title: "操作",
      dataIndex: "operation",
      width: "10%",
      render: (value, record) => {
        const menu = (
          <Menu>
            <Menu.Item key="0">
              <Button
                type="link"
                onClick={() => handleAddCriterion("criteria", record.id)}
                disabled={record?.children?.some(
                  item => item.key === "criterion"
                )}
              >
                添加criteria
              </Button>
            </Menu.Item>
            <Menu.Item key="1">
              <Button
                type="link"
                onClick={() => handleAddCriterion("criterion", record.id)}
                disabled={record?.children?.some(
                  item => item.key === "criteria"
                )}
              >
                添加criterion
              </Button>
            </Menu.Item>
          </Menu>
        );
        return (
          <div>
            {record.key === "criteria" && (
              <Dropdown overlay={menu}>
                <Button icon={<PlusOutlined />} type="link" />
              </Dropdown>
            )}
            {record.id !== test.id && (
              <Button
                type="link"
                icon={<DeleteOutlined />}
                onClick={() => handleDeleteCriterion(record.id)}
              />
            )}
          </div>
        );
      }
    }
  ];
  // 获取递归增加新节点后的所有节点
  function addChildren(
    cur: CriterionProps[],
    id: string,
    newChildren: CriterionProps[]
  ): CriterionProps[] {
    return cur.map(value => {
      if (value.id === id) {
        setExpand(prevState => [...prevState, id]);
        return {
          ...value,
          children: (value.children ?? []).concat(newChildren)
        };
      } else if (value.children) {
        return {
          ...value,
          children: addChildren(value.children, id, newChildren)
        };
      } else {
        return value;
      }
    });
  }

  const handleAddCriterion = debounce((type, id) => {
    const newChildren =
      type === "criteria"
        ? [{ id: getId(), "@operator": "AND", key: "criteria" }]
        : [
            {
              "@operator": "equal",
              "@value": "200",
              "@variable": "$(response_code)",
              "@comment": "Response code is 200 OK",
              key: "criterion",
              id: getId()
            }
          ];
    let tempTest: CriterionProps = { ...test };
    if (test.id === id) {
      tempTest = {
        ...test,
        children: (test?.children ?? []).concat(newChildren)
      };
    } else {
      tempTest = {
        ...test,
        children: addChildren(
          test.children as CriterionProps[],
          id,
          newChildren
        )
      };
    }
    props.setData({
      ...data,
      criteria: tempTest
    });
  }, 300);

  const handleDeleteCriterion = (id: string) => {
    // 递归删除节点
    function deleteChildren(
      cur: CriterionProps[],
      id: string
    ): CriterionProps[] {
      return cur.filter(value => {
        if (value.id === id) {
          return false;
        } else if (value.children) {
          value.children = deleteChildren(value.children, id);
          return true;
        } else {
          return true;
        }
      });
    }
    const current = {
      ...data,
      criteria: {
        ...test,
        children: deleteChildren(test.children as CriterionProps[], id)
      }
    };
    props.setData(current);
  };

  const handleExpand = (show: boolean, record: CriterionProps) => {
    if (show) {
      setExpand(prev => [...prev, record.id]);
    } else {
      setExpand(prevState => prevState.filter(item => item !== record.id));
    }
  };
  // 获取递归修改后的所有子节点
  const getNewTestChildren = (
    type: keyof CriterionProps,
    id: string,
    val: any,
    cur: CriterionProps[]
  ): CriterionProps[] => {
    return cur.map(item => {
      if (item.id === id) {
        return {
          ...item,
          [type]: val
        };
      } else if (item.children) {
        return {
          ...item,
          children: getNewTestChildren(type, id, val, item.children)
        };
      } else {
        return item;
      }
    });
  };
  /**
   * 匹配条件表格的onChange事件
   * @param type 当前改变的字段类型
   * @param id   改变的行的id
   * @param val  改变后的值
   */
  const handleTestChange = (
    type: keyof CriterionProps,
    id: string,
    val: any
  ) => {
    const tempTest: CriterionProps = { ...test };
    // 根据id判断修改父级还是子级
    if (id === tempTest.id) {
      tempTest[type] = val;
    } else if (tempTest.children) {
      tempTest.children = getNewTestChildren(type, id, val, tempTest.children);
    }

    const current = {
      ...data,
      criteria: tempTest
    };
    props.setData(current);
  };

  const handleRequestChange = (type: keyof RequestProps, val: any) => {
    const current = {
      ...data,
      request: {
        ...data.request,
        [type]: val
      }
    };
    props.setData(current);
  };

  const handleResponseChange = (type: keyof ResponseProps, val: any) => {
    const current = {
      ...data,
      response: {
        var: {
          ...data.response?.var,
          [type]: val
        }
      }
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

  return (
    <div className="test-component-wrap">
      <Tabs
        defaultActiveKey="1"
        tabBarExtraContent={
          <Space>
            {/*自动保存吧，体验好一点*/}
            {/*<Button type="link" onClick={() => generateJson()}>*/}
            {/*  保存*/}
            {/*</Button>*/}
            <Popconfirm title="确定重置吗？" onConfirm={handleReset}>
              <Button type="link">重置</Button>
            </Popconfirm>
            <Popconfirm title="确定删除吗" onConfirm={handleDelete}>
              <Button type="link" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        }
      >
        {/*请求*/}
        <Tabs.TabPane tab="请求" key="1">
          <Row className="run-test-items">
            <Col span={4}>请求方法：</Col>
            <Col span={14} offset={1}>
              <Select
                value={data?.request?.method}
                onChange={val => handleRequestChange("method", val)}
                style={{ width: "100%" }}
              >
                <Select.Option value="GET">GET</Select.Option>
                <Select.Option value="OPTIONS">OPTIONS</Select.Option>
                <Select.Option value="HEAD">HEAD</Select.Option>
                <Select.Option value="POST">POST</Select.Option>
                <Select.Option value="PUT">PUT</Select.Option>
                <Select.Option value="DELETE">DELETE</Select.Option>
                <Select.Option value="PROPRFIND">PROPRFIND</Select.Option>
                <Select.Option value="TRACE">TRACE</Select.Option>
                <Select.Option value="CONNECT">CONNECT</Select.Option>
                <Select.Option value="MOVE">MOVE</Select.Option>
                <Select.Option value="TRACK">TRACK</Select.Option>
              </Select>
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>URL：</Col>
            <Col span={14} offset={1}>
              <Input
                value={data?.request?.url}
                onChange={event =>
                  handleRequestChange("url", event.target.value)
                }
              />
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>Version：</Col>
            <Col span={14} offset={1}>
              <Select
                value={data?.request?.version}
                onChange={val => handleRequestChange("version", val)}
                style={{ width: "100%" }}
              >
                <Select.Option value="HTTP/1.0">HTTP/1.0</Select.Option>
                <Select.Option value="HTTP/1.1">HTTP/1.1</Select.Option>
                <Select.Option value="HTTP/2.0">HTTP/2.0</Select.Option>
                <Select.Option value="HTTP/0.9">HTTP/0.9</Select.Option>
              </Select>
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>Cookies：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.request?.cookies}
                onChange={event =>
                  handleRequestChange("cookies", event.target.value)
                }
              />
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>自定义头：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.request?.custom_headers}
                onChange={event =>
                  handleRequestChange("custom_headers", event.target.value)
                }
              />
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>POST数据：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.request?.post_text}
                onChange={event =>
                  handleRequestChange("post_text", event.target.value)
                }
              />
            </Col>
          </Row>
        </Tabs.TabPane>
        {/*响应*/}
        <Tabs.TabPane tab="响应" key="2">
          <Row className="run-test-items">
            <Col span={4}>name：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.response?.var["@name"]}
                onChange={event =>
                  handleResponseChange("@name", event.target.value)
                }
              />
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>description：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.response?.var["@description"]}
                onChange={event =>
                  handleResponseChange("@description", event.target.value)
                }
              />
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>source：</Col>
            <Col span={14} offset={1}>
              {/*<Input.TextArea*/}
              {/*  value={data?.response?.var["@source"]}*/}
              {/*  onChange={event =>*/}
              {/*    handleResponseChange("@source", event.target.value)*/}
              {/*  }*/}
              {/*/>*/}
              <Select
                value={data?.response?.var["@source"]}
                onChange={val => handleResponseChange("@source", val)}
                style={{ width: "100%" }}
              >
                <Select.Option value="statusline">statusline</Select.Option>
                <Select.Option value="header">header</Select.Option>
                <Select.Option value="body">body</Select.Option>
              </Select>
            </Col>
          </Row>

          <Row className="run-test-items">
            <Col span={4}>value：</Col>
            <Col span={14} offset={1}>
              <Input.TextArea
                value={data?.response?.var["#text"]}
                onChange={event =>
                  handleResponseChange("#text", event.target.value)
                }
              />
            </Col>
          </Row>
        </Tabs.TabPane>
        {/*匹配条件*/}
        <Tabs.TabPane tab="匹配条件" key="3">
          <Table
            columns={columns}
            dataSource={[test]}
            pagination={false}
            childrenColumnName="children"
            expandedRowKeys={expand}
            onExpand={handleExpand}
            rowKey="id"
          />
        </Tabs.TabPane>
      </Tabs>
    </div>
  );
};

export default TestComponent;
