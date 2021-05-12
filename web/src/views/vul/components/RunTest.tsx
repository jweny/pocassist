import React, {
  useCallback,
  useEffect,
  useImperativeHandle,
  useState
} from "react";
import {
  Button,
  Col,
  Input,
  message,
  Modal,
  Row,
  Spin,
  Table,
  Tabs
} from "antd";
import { ColumnProps } from "antd/es/table";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons/lib";
import TestComponent, {
  CriterionProps,
  TestComponentProps
} from "./TestComponent";
import { sendVul, testVul } from "../../../api/vul";

interface RunTestProps {
  testData: any; //默认值
  style: any;
  vulId?: number;
}
interface ItemDataProps {
  "@name"?: string;
  "@value"?: string;
  key?: string;
}
const RunTest: React.FC<RunTestProps> = (props, ref) => {
  const { testData, vulId } = props;
  // xml_value
  const [data, setData] = useState<TestComponentProps["data"][]>([]);
  // xml_item
  const [itemData, setItemData] = useState<ItemDataProps[]>([]);
  // 其他xml信息
  const [formData, setFormData] = useState<{
    xml_name: string;
    xml_test_url: string;
  }>({ xml_name: "", xml_test_url: "" });
  // 运行测试结果
  const [testResult, setTestResult] = useState<any>([]);
  // 展示测试结果Modal框
  const [show, setShow] = useState<boolean>(false);
  // modal框loading
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (!!testData) {
      // console.log(testData);
      const valueData =
        testData.xml_value && testData.xml_value.test
          ? Array.isArray(testData.xml_value.test)
            ? testData.xml_value.test
            : [testData.xml_value.test]
          : [];
      const itemData = testData.xml_item
        ? Array.isArray(testData.xml_item)
          ? testData.xml_item
          : [testData.xml_item]
        : [];

      const valueFormatData = getFormatValue(valueData);
      setData(valueFormatData);
      // setData(valueData?.map((item: any) => ({ ...item, id: getId() })) || []);

      setItemData(
        itemData?.map((item: any) => ({ ...item.param, key: getId() })) || []
      );

      setFormData({
        xml_name: testData.xml_name,
        xml_test_url: testData.xml_test_url
      });
    }
  }, [testData]);

  function getCriteriaChildren(data: CriterionProps | undefined) {
    if (data?.criteria) {
      const tempCriteria = Array.isArray(data.criteria)
        ? data.criteria
        : [data.criteria];

      return tempCriteria.map(value => {
        // console.log(value);
        const tempChildren: CriterionProps = {
          id: getId(),
          ...value,
          key: "criteria",
          children: getCriteriaChildren(value)
        };
        return tempChildren;
      });
    }
    return Array.isArray(data?.criterion)
      ? data?.criterion?.map(item => ({
          key: "criterion",
          id: getId(),
          ...item
        }))
      : data?.criterion
      ? [
          {
            key: "criterion",
            id: getId(),
            ...(data.criterion as CriterionProps)
          }
        ]
      : [];
  }
  /**
   * 将数据格式化成表格所需格式
   * criteria和criterion用children代替，对象格式的完善成数组，添加id用于检索
   **/
  const getFormatValue = (data: TestComponentProps["data"][]) => {
    return data.map(value => {
      return {
        ...value,
        criteria: {
          ...value.criteria,
          id: getId(),
          key: "criteria",
          children: getCriteriaChildren(value.criteria)
        },
        id: getId()
      };
    });
  };

  const handleAddTest = () => {
    setData(prevState => {
      return [
        ...prevState,
        {
          id: getId(),
          response: {
            var: {
              "@name": "response_code",
              "@source": "statusline",
              "#text": "^.*\\s(\\d\\d\\d)\\s"
            }
          },
          request: {
            method: "GET",
            url: "$(scheme)://$(host):$(port)$(path)/",
            version: "HTTP/1.1"
          },
          criteria: {
            id: getId(),
            key: "criteria",
            "@operator": "AND",
            children: [
              {
                id: getId(),
                key: "criterion",
                "@operator": "equal",
                "@variable": "$(response_code)",
                "@value": "200",
                "@comment": "Response code is 200 OK"
              }
            ]
          }
        }
      ];
    });
  };
  const handleTestFinish = (value: TestComponentProps["data"]) => {
    // console.log(value);
    const curData = data.map(item => {
      // console.log(item.test.id, value.test.id);
      if (item.id === value.id) {
        return value;
      }
      return item;
    });
    setData(curData);
  };

  const handleDeleteTest = (id: string) => {
    const curData = data.filter(item => {
      // console.log(item.test.id, value.test.id);
      return item.id !== id;
    });
    setData(curData);
  };

  const handleRunTest = () => {
    setShow(true);
    setLoading(true);
    testVul(vulId, { target: formData.xml_test_url })
      .then(res => {
        setTestResult(res.data);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  // const handleSend = () => {
  //   setLoading(true);
  //   sendVul(vulId)
  //     .then((res: any) => {
  //       console.log(res);
  //       message.success(res.data);
  //       // setTestResult(res.data);
  //     })
  //     .finally(() => {
  //       setLoading(false);
  //     });
  // };

  const columns: ColumnProps<ItemDataProps>[] = [
    {
      title: "name",
      dataIndex: "@name",
      render: (value, record) => (
        <Input
          value={value}
          onChange={event =>
            handleItemChange(event.target.value, "@name", record)
          }
        />
      )
    },
    {
      title: "value",
      dataIndex: "@value",
      render: (value, record) => (
        <Input
          value={value}
          onChange={event =>
            handleItemChange(event.target.value, "@value", record)
          }
        />
      )
    },
    {
      title: "操作",
      dataIndex: "operation",
      render: (value, record) => {
        return (
          <Button
            type="link"
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteItem(record.key as string)}
          />
        );
      }
    }
  ];
  const handleAddItem = () => {
    setItemData(prevState => [...prevState, { key: getId() }]);
  };
  const handleDeleteItem = (key: string) => {
    // console.log(itemData, key);
    setItemData(prevState => prevState.filter(item => item.key !== key));
  };
  const handleItemChange = (
    val: string,
    type: string,
    record: ItemDataProps
  ) => {
    setItemData(prevState => {
      return prevState.map(item => {
        if (item.key === record.key) {
          return {
            ...item,
            [type]: val
          };
        }
        return item;
      });
    });
  };

  const dataFormatter = (data: TestComponentProps["data"][]) => {
    // 将criteria转换为接口所需格式
    function getCriteriaChildren(value: CriterionProps | undefined): any {
      if (value) {
        const { id, key, children, criteria, criterion, ...rest } = value;
        if (children) {
          return {
            ...rest,
            [children[0].key as string]:
              children.length > 1
                ? children.map(item => getCriteriaChildren(item))
                : getCriteriaChildren(children[0])
          };
        } else {
          return rest;
        }
      } else {
        return null;
      }
    }
    const returnData = data.map(item => {
      const { request, response, criteria } = item;
      const formatCriteria = getCriteriaChildren(criteria);
      return {
        request,
        response,
        criteria: formatCriteria
      };
    });
    return returnData.length > 0
      ? { test: returnData.length === 1 ? returnData[0] : returnData }
      : null;
  };

  const itemDataFormatter = (data: ItemDataProps[]) => {
    const returnData = data.map(item => {
      const { key, ...rest } = item;
      return {
        param: { ...rest }
      };
    });
    return returnData.length > 0 ? returnData : null;
  };
  // ref转发
  const getRunTestData = useCallback(() => {
    // console.log(data, itemData);
    const xml_value = dataFormatter(data);
    const xml_item = itemDataFormatter(itemData);
    return {
      ...formData,
      xml_value,
      xml_item
    };
  }, [data, itemData, formData]);

  useImperativeHandle(
    ref,
    () => {
      return { getRunTestData };
    },
    [getRunTestData]
  );
  return (
    <div style={props.style} className="run-test-wrap">
      <Row>
        <Col span={3} className="rt-url-label">
          测试URL：
        </Col>
        <Col span={14}>
          <Input
            value={formData?.xml_test_url}
            onChange={e => {
              e.persist();
              setFormData(prevState => ({
                ...prevState,
                xml_test_url: e.target.value
              }));
            }}
          />
        </Col>
        <Col span={2} offset={2}>
          <Button type="primary" onClick={handleRunTest} disabled={!vulId}>
            启动测试
          </Button>
        </Col>
        {/*<Col span={2} style={{ marginLeft: 10 }}>*/}
        {/*  <Button type="primary" onClick={handleSend} disabled={!vulId}>*/}
        {/*    发送*/}
        {/*  </Button>*/}
        {/*</Col>*/}
      </Row>
      <Row>
        <Col span={3} className="rt-url-label">
          XML文件名：
        </Col>
        <Col span={14}>
          <Input
            value={formData?.xml_name}
            onChange={e => {
              e.persist();
              setFormData(prevState => ({
                ...prevState,
                xml_name: e.target.value
              }));
            }}
          />
        </Col>
      </Row>
      <Button type="link" onClick={handleAddTest}>
        添加测试 <PlusOutlined />
      </Button>
      <Button type="link" onClick={handleAddItem}>
        添加Item <PlusOutlined />
      </Button>
      {/*item列表*/}
      {itemData?.length > 0 && (
        <Table dataSource={itemData} pagination={false} columns={columns} />
      )}
      {/*test列表*/}
      {data.map((item, index) => {
        return (
          item && (
            <TestComponent
              data={item}
              setData={handleTestFinish}
              delete={handleDeleteTest}
              key={index}
            />
          )
        );
      })}
      <Modal
        visible={show}
        width={document.documentElement.offsetWidth}
        footer={null}
        onCancel={() => {
          setShow(false);
          setTestResult([]);
        }}
        title="测试结果"
        className="test-result-wrap"
      >
        <Spin spinning={loading} size="large">
          <div style={{ minHeight: 500 }}>
            <Tabs defaultActiveKey="0">
              {testResult?.map((item: any, index: number) => {
                return (
                  <Tabs.TabPane tab={`Test ${index + 1}`} key={index}>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        测试目标：
                      </Col>
                      <Col span={18}>
                        <pre>{item?.target}</pre>
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        是否为漏洞：
                      </Col>
                      <Col span={18}>
                        {item?.is_vulerable === 1 ? "是" : "否"}
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        请求头：
                      </Col>
                      <Col span={18}>
                        <pre>{item?.req_header}</pre>
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        请求数据：
                      </Col>
                      <Col span={18}>
                        <pre>{item?.req_data}</pre>
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        响应头：
                      </Col>
                      <Col span={18}>
                        <pre>{item?.resp_header}</pre>
                      </Col>
                    </Row>
                    <Row>
                      <Col
                        span={4}
                        style={{ textAlign: "right", marginRight: 20 }}
                      >
                        响应body：
                      </Col>
                      <Col span={18}>
                        <pre>{item?.resp_body}</pre>
                      </Col>
                    </Row>
                  </Tabs.TabPane>
                );
              })}
            </Tabs>
            {!testResult && (
              <div style={{ height: "100%", textAlign: "center" }}>
                暂无数据
              </div>
            )}
          </div>
        </Spin>
      </Modal>
    </div>
  );
};

// @ts-ignore
export default React.forwardRef(RunTest);

/**
 * 随机生成一个id
 */
export function getId() {
  return new Date().getTime() + "" + Math.floor(Math.random() * 10000);
}
