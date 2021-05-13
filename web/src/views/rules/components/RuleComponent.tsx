import React, { useContext, useEffect, useState } from "react";
import { Button, Col, Input, message, Row, Select, Switch, Tabs } from "antd";
import { getId } from "./RunTest";
import InputGroup from "./InputGroup";

interface RuleComponentProps {
  data: any;
  setData: any;
}
const RuleComponent: React.FC<RuleComponentProps> = props => {
  console.log(props);
  const { data } = props;
  const [headerList, setHeaderList] = useState<any[]>([{ id: getId() }]);

  useEffect(() => {
    if (data.headers) {
      if (Object.keys(data.headers).length > 0) {
        let headerListNext = Object.keys(data.headers).map(item => ({
          key: item,
          value: data.headers[item],
          id: getId()
        }));
        setHeaderList(headerListNext);
      } else {
        setHeaderList([{ id: getId() }]);
      }
    }
  }, [data]);

  const handleAdd = () => {
    setHeaderList(prevState => {
      return [...prevState, { id: getId() }];
    });
  };

  const handleDelete = (id: string) => {
    setHeaderList(prevState => {
      return prevState.filter(item => item.id !== id);
    });
  };

  const handleChange = (val: any) => {
    console.log(val);
    setHeaderList(prevState => {
      return prevState.map(item => {
        if (item.id === val.id) {
          return val;
        } else {
          return item;
        }
      });
    });
  };

  const handleHeaderSubmit = () => {
    let newVal: any = {};
    headerList.forEach(header => {
      console.log(header);
      if (header.key) {
        newVal[header.key] = header.value;
      }
    });
    console.log(newVal);
    props.setData(data.id, "headers", newVal);
    message.success("保存请求头成功");
  };

  return (
    <div>
      <Row className="run-test-items">
        <Col span={4}>请求方法 method：</Col>
        <Col span={14} offset={1}>
          <Select
            value={data?.method}
            onChange={val => props.setData(data.id, "method", val)}
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
        <Col span={4}>请求路径 path：</Col>
        <Col span={14} offset={1}>
          <Input
            value={data?.path}
            onChange={event =>
              props.setData(data.id, "path", event.target.value)
            }
          />
        </Col>
      </Row>

      <Row className="run-test-items">
        <Col span={4}>请求头 headers：</Col>
        <Col span={14} offset={1}>
          {headerList.map(item => {
            return (
              <InputGroup
                key={item.id}
                value={item}
                handleAdd={handleAdd}
                handleDelete={handleDelete}
                onChange={(val: any) => handleChange(val)}
                showDeleteButton={headerList.length > 1}
                className="header-input-group"
              />
            );
          })}
          {/*<Input.TextArea*/}
          {/*  value={data?.headers}*/}
          {/*  onChange={event =>*/}
          {/*    props.setData(data.id, "headers", event.target.value)*/}
          {/*  }*/}
          {/*/>*/}
          <Button onClick={handleHeaderSubmit} type="primary" ghost>
            保存请求头
          </Button>
          <span style={{ color: "red", marginLeft: "10px" }}>
            注意：填写完请求头后请点击该按钮，否则不会保存请求头！
          </span>
        </Col>
      </Row>

      <Row className="run-test-items">
        <Col span={4}>请求体 body：</Col>
        <Col span={14} offset={1}>
          <Input.TextArea
            value={data?.body}
            onChange={event =>
              props.setData(data.id, "body", event.target.value)
            }
          />
        </Col>
      </Row>

      <Row className="run-test-items">
        <Col span={4}>允许跳转 follow_redirects：</Col>
        <Col span={14} offset={1}>
          <Switch
            checked={data?.follow_redirects}
            onChange={checked =>
              props.setData(data.id, "follow_redirects", checked)
            }
          />
        </Col>
      </Row>

      <Row className="run-test-items">
        <Col span={4}>提取规则 search：</Col>
        <Col span={14} offset={1}>
          <Input.TextArea
            value={data?.search}
            onChange={event =>
              props.setData(data.id, "search", event.target.value)
            }
          />
        </Col>
      </Row>

      <Row className="run-test-items">
        <Col span={4}>表达式 expression：</Col>
        <Col span={14} offset={1}>
          <Input.TextArea
            value={data?.expression}
            onChange={event =>
              props.setData(data.id, "expression", event.target.value)
            }
          />
        </Col>
      </Row>
    </div>
  );
};

export default RuleComponent;
