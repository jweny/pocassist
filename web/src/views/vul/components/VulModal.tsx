import React, { useContext, useEffect, useRef, useState } from "react";
import {
  Alert,
  Button,
  Col,
  Form,
  Input,
  message,
  Modal,
  Row,
  Select,
  Spin
} from "antd";
import { ModalProps } from "antd/es/modal";
import { FormColumnProps } from "./SearchForm";
import {
  createProduct,
  getProductList,
  ProductDataProps
} from "../../../api/webapp";
import VulContext from "../../../store/vul/store";
import {
  createVul,
  getVulDetail,
  getVulList,
  updateVul,
  VulDataProps
} from "../../../api/vul";
import { getUserInfo } from "../../../utils/auth";
import BraftEditor, { BuiltInControlType, ControlType } from "braft-editor";
import "braft-editor/dist/index.css";
import { richFormColumns } from "./columns";
import { PlusOutlined } from "@ant-design/icons/lib";
import AddModal from "../../modules/components/AddModal";

interface AddVulProps extends ModalProps {
  selected?: VulDataProps;
  type?: string;
}
const VulModal: React.FC<AddVulProps> = props => {
  // console.log(props);
  let { selected, type = "vul" } = props;
  const testRef = useRef(null);
  const [step, setStep] = useState<number>(1);
  const [vulData, setVulData] = useState<any>(null);
  // 添加完成后将vul和xml的id暂存，如果xml保存出错或者运行测试之后又要修改前面的信息，通过id修改已有漏洞或xml
  const [vulAddId, setVulAddId] = useState<number | undefined>(undefined);
  const [loading, setLoading] = useState<boolean>(true);
  const { state, dispatch } = useContext(VulContext);
  const [addProduct, setAddProduct] = useState<boolean>(false);

  const [form] = Form.useForm();
  const formItemLayout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 16 }
  };

  useEffect(() => {
    // console.log(selected);
    // 根据selected判断当前是编辑还是新增
    if (!!selected) {
      setVulData(selected);
      setStep(1);
      setLoading(false);
    } else {
      //setTestData(undefined); 测试下需不需要这一行
      setLoading(false);
    }
  }, [selected, type]);

  const formColumns: FormColumnProps[] = [
    {
      name: "webapp",
      label: "影响组件",
      render: () => {
        return (
          <Select
            placeholder="请选择"
            style={{ width: 200 }}
            onSelect={val => {
              if (val === "new") {
                setAddProduct(true);
              }
            }}
          >
            {state.productList?.map(item => {
              return (
                <Select.Option value={item.id as number} key={item.id}>
                  {item.name}
                </Select.Option>
              );
            })}
            <Select.Option value="new">
              新增 <PlusOutlined />
            </Select.Option>
          </Select>
        );
      }
    },
    {
      name: "severity",
      label: "漏洞等级",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 200 }}>
            {state.basic?.VulLevel.map(item => {
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
      name: "category",
      label: "漏洞类型",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 200 }}>
            {state.basic?.VulType.map(item => {
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
      name: "language",
      label: "漏洞语言",
      render: () => {
        return (
          <Select placeholder="请选择" style={{ width: 200 }}>
            {state.basic?.VulLanguage.map(item => {
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
      name: "cve",
      label: "CVE编号"
    },
    {
      name: "cnnvd",
      label: "CNNVD编号"
    }
  ];

  // 富文本编辑器的控制项
  const controls: BuiltInControlType[] = [
    "font-size",
    "line-height",
    "letter-spacing",
    "headings",
    "text-indent",
    "text-align",
    "emoji",
    "hr",
    "separator"
  ];

  const handleStepChange = (val: number) => {
    if (val === 2) {
      form
        .validateFields()
        .then(res => {
          setStep(val);
        })
        .catch(() => {});
    } else {
      setStep(val);
    }
  };

  const handleFinish = (val?: any) => {
    setLoading(true);
    const userInfo = getUserInfo();
    const vulApi = vulData?.id || vulAddId ? updateVul : createVul;

    // @ts-ignore
    const xml = testRef.current?.getRunTestData();

    const finalData: VulDataProps = {
      ...vulData,
      ...val,
      writer_id: userInfo?.id,
      ...xml
    };
    // console.log(finalData);
    vulApi(finalData, vulData?.id || vulAddId)
      .then(res => {
        setVulAddId(res.data.id);
        message.success(`保存成功`);
        getVulList({
          ...state.search_query,
          page: state.page,
          pagesize: state.pagesize
          // search_query: JSON.stringify(state.search_query)
        }).then(res => {
          dispatch({ type: "SET_LIST", payload: res.data.data });
          dispatch({ type: "SET_TOTAL", payload: res.data.total });
        });
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const handleAddProduct = (val: any) => {
    // console.log(val);
    createProduct(val).then(res => {
      setAddProduct(false);
      getProductList({ page: 1, pagesize: 9999 }).then(res => {
        dispatch({ type: "SET_PRODUCT_LIST", payload: res.data.data });
      });
      form.setFieldsValue({ product: res.data.id });
    });
  };


  const height = document.documentElement.offsetHeight;

  return (
    <Modal
      {...props}
      destroyOnClose
      // forceRender
      maskClosable={false}
      onCancel={e => {
        handleStepChange(1);
        props.onCancel && props.onCancel(e);
      }}
      footer={
        <div>
          {step === 1 && (
            <React.Fragment>
              <Button
                type="primary"
                onClick={() => {
                  form
                    .validateFields()
                    .then(data => {
                      const cur = { ...data };
                      richFormColumns.map(item => {
                        cur[item.name as string] =
                          data[item.name as string] &&
                          typeof data[item.name as string] === "object"
                            ? data[item.name as string].toHTML()
                            : data[item.name as string];
                      });
                      setVulData((prev: any) => ({ ...prev, ...cur }));
                      console.log({ ...vulData, ...cur });
                      handleFinish({ ...vulData, ...cur });
                    })
                    .catch(err => {});
                }}
              >
                保存
              </Button>
            </React.Fragment>
          )}

          <Button onClick={props.onCancel}>取消</Button>
        </div>
      }
      bodyStyle={{
        maxHeight: height - 120,
        overflowY: "scroll"
      }}
    >
      <Spin spinning={loading}>
        {/*{step === 1 && (*/}
        <Form
          {...formItemLayout}
          form={form}
          name="vul_detail"
          layout="horizontal"
          initialValues={vulData}
          style={{ display: `${step === 1 ? "flex" : "none"}` }}
        >
          <Form.Item
            name="name_zh"
            label="漏洞名称"
            rules={[{ required: true }]}
            className="form-item-long"
            labelCol={{ span: 3 }}
            wrapperCol={{ span: 20 }}
          >
            <Input placeholder="请输入漏洞名称" />
          </Form.Item>
          {formColumns.map(
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
          <Row style={{ width: "100%", marginBottom: 24 }}>
            <Col offset={3}>
              <Alert message="注意以下为富文本：复制粘贴来的文本先清除格式" />
            </Col>
          </Row>
          {richFormColumns.map(({ name, label, ...formProps }, index) => (
            <Form.Item
              label={label}
              key={index}
              name={name}
              labelCol={{ span: 3 }}
              wrapperCol={{ span: 20 }}
              validateTrigger="onBlur"
              getValueProps={value => {
                return { value: BraftEditor.createEditorState(value) };
              }}
              className="form-item-long"
              {...formProps}
            >
              <BraftEditor
                excludeControls={controls}
                placeholder={`请输入${label}`}
                style={{
                  border: "1px solid #ccc"
                }}
                contentStyle={{
                  height: 200,
                  overflow: "auto"
                }}
              />
            </Form.Item>
          ))}
        </Form>
      </Spin>
      <AddModal
        visible={addProduct}
        onCancel={() => {
          setAddProduct(false);
        }}
        onOk={handleAddProduct}
        type="product"
        selected={undefined}
      />
    </Modal>
  );
};


export default VulModal;
