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
  Radio,
  Row,
  Table,
} from "antd";
import { ColumnProps } from "antd/es/table";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons/lib";
import TestComponent, {
  TestComponentProps
} from "./TestComponent";

interface RunTestProps {
  testData: any; //默认值
  style: any;
  vulId?: number;
}
interface ItemDataProps {
  name?: string;
  value?: string;
  key?: string;
}
const RunTest: React.FC<RunTestProps> = (props, ref) => {
  const { testData, vulId } = props;
  // poc rules、groups 信息
  const [data, setData] = useState<TestComponentProps["data"][]>([]);
  const [groupData, setGroupData] = useState<any>([]);
  // poc set信息
  const [itemData, setItemData] = useState<ItemDataProps[]>([]);
  // 其他poc信息
  const [formData, setFormData] = useState<{
    name: string;
    type: string;
  }>({ name: "", type: "rules" });
  // 运行测试结果
  const [testResult, setTestResult] = useState<any>([]);
  // 展示测试结果Modal框
  const [show, setShow] = useState<boolean>(false);
  // modal框loading
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (!!testData) {
      console.log(testData);
      const itemData = testData.set
        ? Object.keys(testData.set).map(item => {
            return { name: item, value: testData.set[item] };
          })
        : [];

      const valueFormatData = getFormatValue(testData.rules || []);
      const groupFormatData = getFormatGroup(testData.groups || {});
      // 将format之后的rules和groups存入state
      setData(valueFormatData);
      if (testData.hasOwnProperty("rules")) {
        setGroupData([{ name: "rules", data: valueFormatData, id: getId() }]);
      } else {
        setGroupData(groupFormatData);
      }

      setItemData(
        itemData?.map((item: any) => ({ ...item, key: getId() })) || []
      );

      setFormData({
        name: testData.name,
        type: testData.hasOwnProperty("rules") ? "rules" : "groups"
      });
    } else {
      setGroupData([{ name: "rules", data: [{ id: getId() }], id: getId() }]);
    }
  }, [testData]);
  /**
   * 将数据格式化成表格所需格式
   * rules和groups分别处理
   **/
  const getFormatValue = (data: TestComponentProps["data"][]) => {
    return data.map(value => {
      return {
        ...value,
        id: getId()
      };
    });
  };

  const getFormatGroup = (data: any) => {
    return Object.keys(data).map(item => {
      return {
        name: item,
        data: data[item].map((val: any) => ({ ...val, id: getId() })),
        id: getId()
      };
    });
  };

  const handleAddTest = () => {
    setGroupData([
      ...groupData,
      { name: "rules", data: [{ id: getId() }], id: getId() }
    ]);
  };

  const handleTestFinish = (value: TestComponentProps["data"]) => {
    console.log(value);
    const curData = groupData.map((item: any) => {
      if (item.id === value.id) {
        return value;
      }
      return item;
    });
    setGroupData(curData);
  };

  const handleDeleteTest = (id: string) => {
    const curData = groupData.filter((item: any) => {
      // console.log(item.test.id, value.test.id);
      return item.id !== id;
    });
    setGroupData(curData);
  };

  const columns: ColumnProps<ItemDataProps>[] = [
    {
      title: "name",
      dataIndex: "name",
      render: (value, record) => (
        <Input
          value={value}
          onChange={event =>
            handleItemChange(event.target.value, "name", record)
          }
        />
      )
    },
    {
      title: "value",
      dataIndex: "value",
      render: (value, record) => (
        <Input
          value={value}
          onChange={event =>
            handleItemChange(event.target.value, "value", record)
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

  const itemDataFormatter = (data: ItemDataProps[]) => {
    let returnData: any = {};

    data.forEach(item => {
      returnData[item.name as string] = item.value;
    });
    return returnData;
  };
  // ref转发
  const getRunTestData = useCallback(() => {
    // TODO 从runtest取值应该没问题了，rules的增删还没做，目前只是用现有数据修改的，从0添加一个规则还没测过
    console.log(data, "------", groupData, "--------", itemData);
    let json_poc: any = {};
    json_poc.name = formData.name;
    json_poc.set = itemDataFormatter(itemData);
    if (formData.type === "rules") {
      json_poc.rules = groupData[0]?.data?.map((item: any) => {
        const { id, ...rest } = item;
        return rest;
      });
    } else {
      let group: any = {};
      groupData.forEach((item: any) => {
        group[item.name] = item.data;
      });
      json_poc.groups = group;
    }
    return {
      json_poc
    };
  }, [data, itemData, formData, groupData]);

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
          名称：
        </Col>
        <Col span={14}>
          <Input
            value={formData?.name}
            onChange={e => {
              e.persist();
              setFormData(prevState => ({
                ...prevState,
                name: e.target.value
              }));
            }}
          />
        </Col>
      </Row>
      <Row>
        <Col span={3} className="rt-url-label">
          规则类型：
        </Col>
        <Col span={14}>
          <Radio.Group
            onChange={e => {
              setFormData(prevState => ({
                ...prevState,
                type: e.target.value
              }));
            }}
            value={formData?.type}
          >
            <Radio value="rules">rules</Radio>
            <Radio value="groups">groups</Radio>
          </Radio.Group>
        </Col>
      </Row>
      <Button type="link" onClick={handleAddItem}>
        添加变量 <PlusOutlined />
      </Button>
      {formData.type === "groups" && (
        <Button type="link" onClick={handleAddTest}>
          添加rules <PlusOutlined />
        </Button>
      )}
      {/*item列表*/}
      {itemData?.length > 0 && (
        <Table dataSource={itemData} pagination={false} columns={columns} />
      )}
      {/*test列表*/}
      {groupData?.map((item: any, index: number) => {
        return (
          item && (
            <TestComponent
              data={item}
              type={formData.type}
              setData={handleTestFinish}
              delete={handleDeleteTest}
              key={index}
            />
          )
        );
      })}
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
