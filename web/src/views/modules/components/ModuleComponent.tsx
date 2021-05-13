import React, {
  useCallback,
  useEffect,
  useMemo,
  useState
} from "react";
import {
  Button,
  Form,
  Input,
  message,
  Popconfirm,
  Space,
  Table
} from "antd";
import { PlusOutlined } from "@ant-design/icons/lib";
import CustomColumns from "../../../components/CustomColumns";
import { allProductColumns } from "./columns";
import { ColumnProps } from "antd/es/table";
import "../index.less";
import {
  createProduct,
  deleteProduct,
  getProductList,
  updateProduct
} from "../../../api/webapp";
import AddModal from "./AddModal";
import { getVulBasic } from "../../../api/vul";

export interface QueryProps {
  page: number;
  pagesize: number;
  search_query?: string;
}

interface ModuleComponentProps {
  type: "module" | "product" | "script";
}

const ModuleComponent: React.FC<ModuleComponentProps> = props => {
  const [form] = Form.useForm();
  // 因为要翻译一下script_type，因此将script的columns在当前组件中定义，并通过接口获取type的数组
  const [scriptType, setScriptType] = useState<any[]>([]);

  const allScriptColumns = useMemo(
    () => [
      { title: "模块名称", dataIndex: "name", width: "20%" },
      {
        title: "脚本类型",
        dataIndex: "type_id",
        render: (text: number) => {
          return scriptType.find((item: any) => item.id === text)?.type_name;
        },
        width: "20%"
      },
      { title: "备注", dataIndex: "remarks", width: "30%" }
    ],
    [scriptType]
  );

  const getAllColumns = useCallback(() => {
    switch (props.type) {
      case "product":
        return allProductColumns;
      default:
        return allProductColumns;
    }
  }, [props.type, allScriptColumns]);

  const allColumns = getAllColumns();

  const [columns, setColumns] = useState<ColumnProps<any>[]>([...allColumns]);
  const [loading, setLoading] = useState<boolean>(false);
  const [query, setQuery] = useState<QueryProps>({
    page: 1,
    pagesize: 20
  });
  const [list, setList] = useState<any[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [modalShow, setModalShow] = useState<boolean>(false);
  const [selected, setSelected] = useState<any>(undefined);

  const getTableData = useCallback(
    (params: QueryProps) => {
      setLoading(true);
      let api;
      switch (props.type) {
        case "product":
          api = getProductList;
          break;
        default:
          api = getProductList;
      }
      api(params)
        .then(res => {
          // console.log(res);
          setList(res.data.data);
          setTotal(res.data.total);
        })
        .finally(() => {
          setLoading(false);
        });
    },
    [props.type]
  );
  // 更新列表
  useEffect(() => {
    getTableData(query);
  }, [query, getTableData]);
  // 获取script type
  useEffect(() => {
    if (props.type === "script") {
      getVulBasic().then(res => {
        setScriptType(res.data.ScriptType);
      });
    }
  }, [props.type]);
  // 当allColumns改变时更新columns，解决获取脚本类型之后，列表上不显示脚本类型的问题
  useEffect(() => {
    setColumns([...getAllColumns()]);
  }, [getAllColumns]);

  const realColumns: ColumnProps<any>[] = [
    {
      title: "序号",
      dataIndex: "order",
      render: (value, record, index) => index + 1,
      width: "15%"
    },
    ...columns,
    {
      title: "操作",
      dataIndex: "operation",
      render: (value, record) => {
        return (
          <div>
            <Button
              type="link"
              onClick={() => {
                setSelected(record);
                handleToggleAdd();
              }}
            >
              编辑
            </Button>
            <Popconfirm
              title={`确定删除该${
                props.type === "script"
                  ? "脚本"
                  : props.type === "product"
                  ? "组件"
                  : "模块"
              }吗？`}
              onConfirm={() => handleDelete(record.id)}
              onCancel={() => {}}
            >
              <Button type="link">删除</Button>
            </Popconfirm>
          </div>
        );
      },
      width: "15%"
    }
  ];

  const handleDelete = (id: number) => {
    // console.log(id);
    let api;
    switch (props.type) {
      case "product":
        api = deleteProduct;
        break;
      default:
        api = deleteProduct;
    }
    api(id).then(res => {
      // console.log(res);
      message.success("删除成功");
      getTableData(query);
    });
  };

  const handlePageChange = (page: number, pageSize: number | undefined) => {
    setQuery(prevState => {
      return {
        ...prevState,
        page,
        pagesize: pageSize as number
      };
    });
  };

  const handleFinish = (val: any) => {
    console.log(val);
    setQuery(prevState => {
      return {
        ...prevState,
        page: 1,
        // search_query: JSON.stringify(val)
        ...val
      };
    });
  };

  const handleToggleAdd = () => {
    setModalShow(prevState => {
      // 关闭modal框清除默认值
      if (prevState) {
        setSelected(undefined);
      }
      return !prevState;
    });
  };

  const handleModalFinish = (value: any) => {
    console.log(value);
    let api;
    switch (props.type) {
      case "product":
        api = !!selected ? updateProduct : createProduct;
        break;
      default:
        api = !!selected ? updateProduct : createProduct;
    }
    api(value, selected?.id).then(res => {
      console.log(res);
      getTableData(query);
      handleToggleAdd();
    });
  };

  return (
    <div className="module-manage-wrap">
      <Form form={form} layout="inline" onFinish={handleFinish}>
        <Form.Item label="模糊查询" name="search">
          <Input placeholder="请输入模糊查询" allowClear />
        </Form.Item>
        <Button type="primary" htmlType="submit">
          查询
        </Button>
      </Form>
      <div className="module-table-wrap">
        <div className="module-controller">
          <Space>
            <Button
              icon={<PlusOutlined />}
              type="primary"
              onClick={handleToggleAdd}
            >
              新建
            </Button>
          </Space>
          <CustomColumns
            className="custom-btn"
            allColumns={allColumns}
            columns={columns}
            setColumns={setColumns}
          />
        </div>
        <Table
          columns={realColumns}
          dataSource={list}
          rowKey="id"
          loading={loading}
          pagination={{
            current: query.page,
            pageSize: query.pagesize,
            total: total,
            showTotal: (total: number) => `共${total}条`,
            showQuickJumper: true,
            showSizeChanger: true,
            onChange: handlePageChange
          }}
        />
        <AddModal
          visible={modalShow}
          onCancel={handleToggleAdd}
          onOk={handleModalFinish}
          scriptType={scriptType}
          type={props.type}
          selected={selected}
        />
      </div>
    </div>
  );
};

export default ModuleComponent;
