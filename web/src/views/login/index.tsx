import React, { Fragment } from "react";
import { RouteComponentProps, useHistory } from "react-router-dom";
import { Input, Button, Form, Divider } from "antd";
import "./index.less";

import { MobileOutlined, UnlockOutlined } from "@ant-design/icons/lib";
import { getUserInfos, login } from "../../api/login";
import { setToken, setUserInfo } from "../../utils/auth";

const Login: React.FC<RouteComponentProps> = (props: RouteComponentProps) => {
  const [form] = Form.useForm();
  const history = useHistory();

  const formItemLayout = {
    labelCol: { span: 24 },
    wrapper: { span: 24 }
  };

  const handleFinish = (values: any) => {
    const { captcha_value, ...rest } = values;
    login(rest).then(res => {
      setToken(res?.data?.token);
      getUserInfos().then(res => {
        setUserInfo(res.data);
      });
      history.push("/vul");
    });
  };
  return (
    <div className="apply-login-wrap">
      <div className="apply-login-page">
        <h2 className="login-title">漏洞管理系统</h2>
        <Form
          {...formItemLayout}
          form={form}
          onFinish={handleFinish}
          hideRequiredMark
          size="large"
        >
          <Fragment>
            <Form.Item
              label="账号"
              name="username"
              rules={[{ required: true }]}
            >
              <Input placeholder="请输入账号" prefix={<MobileOutlined />} />
            </Form.Item>
            <Form.Item
              label="密码"
              name="password"
              rules={[{ required: true }]}
            >
              <Input
                placeholder="请输入密码"
                maxLength={40}
                prefix={<UnlockOutlined />}
                type="password"
              />
            </Form.Item>
            {/*<Form.Item*/}
            {/*  label="验证码"*/}
            {/*  name="captcha_value"*/}
            {/*  rules={[{ required: true }]}*/}
            {/*>*/}
            {/*  <Input*/}
            {/*    placeholder="请输入验证码"*/}
            {/*    maxLength={40}*/}
            {/*    prefix={<KeyOutlined />}*/}
            {/*    suffix={*/}
            {/*      <img*/}
            {/*        src={captcha.image}*/}
            {/*        className="login-captcha"*/}
            {/*        onClick={handleGetCaptcha}*/}
            {/*      />*/}
            {/*    }*/}
            {/*  />*/}
            {/*</Form.Item>*/}
          </Fragment>

          <Form.Item wrapperCol={{ span: 24 }}>
            <Button type="primary" htmlType="submit" block>
              登录
            </Button>
          </Form.Item>
        </Form>
        <Divider />
        <div className="login-footer">
          <p>Copyright 2005-2020 360.com 版权所有 360互联网中心</p>
        </div>
      </div>
    </div>
  );
};

export default Login;
