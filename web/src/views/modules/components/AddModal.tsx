import React, { useEffect } from "react";
import { Button, Form, Input, Modal } from "antd";
import { ModalProps } from "antd/es/modal";
import { productFormColumns } from "./columns";

interface AddModuleProps extends ModalProps {
  scriptType?: { name: string; label: string }[];
  type: "script" | "product" | "module";
  selected?: any;
}
const AddModal: React.FC<AddModuleProps> = props => {
  const { type, selected } = props;
  const [form] = Form.useForm();
  const formItemLayout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 16 }
  };

  useEffect(() => {
    form.resetFields();
  }, [selected, form]);

  let columns;
  switch (type) {
    case "product":
      columns = productFormColumns;
      break;
    default:
      columns = productFormColumns;
  }

  const handleFinish = (value: any) => {
    props.onOk && props.onOk(value);
  };

  return (
    <Modal
      {...props}
      forceRender
      title={
        (!!selected ? "编辑" : "新增") +
        (type === "script" ? "脚本" : type === "product" ? "组件" : "模块")
      }
      footer={
        <div>
          <Button type="primary" onClick={form.submit}>
            {!!selected ? "保存" : "添加"}
          </Button>
          <Button onClick={props.onCancel}>取消</Button>
        </div>
      }
    >
      <Form
        {...formItemLayout}
        form={form}
        onFinish={handleFinish}
        initialValues={selected || {}}
      >
        {columns.map(
          (
            {
              name,
              label,
              render = () => <Input placeholder={`请输入${label}`} />,
              ...formProps
            },
            index
          ) => (
            <Form.Item label={label} key={index} name={name} {...formProps}>
              {render()}
            </Form.Item>
          )
        )}
      </Form>
    </Modal>
  );
};

export default AddModal;
