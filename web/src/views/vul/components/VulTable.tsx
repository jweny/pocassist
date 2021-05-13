import React, { useContext, useEffect, useState } from "react";
import { Button, message, Popconfirm, Space, Table } from "antd";
import { ImportOutlined, PlusOutlined } from "@ant-design/icons/lib";
import { ColumnProps } from "antd/es/table";
import VulContext from "../../../store/vul/store";
import CustomColumns from "../../../components/CustomColumns";
import { getAllColumns } from "./columns";
import VulModal from "./VulModal";
import {
  deleteVul,
  getVulDetail,
  getVulList,
  VulDataProps
} from "../../../api/vul";
import { VulComponentProps } from "./SearchForm";

const VulTable: React.FC<VulComponentProps> = props => {
  const allColumns = getAllColumns();

  const [loading, setLoading] = useState<boolean>(false);
  const [columns, setColumns] = useState<ColumnProps<VulDataProps>[]>([
    ...allColumns.slice(0, 4)
  ]);
  const [addShow, setAddShow] = useState<boolean>(false);
  const [selected, setSelected] = useState<VulDataProps | undefined>(undefined);

  const { state, dispatch } = useContext(VulContext);

  const getListData = () => {
    setLoading(true);
    getVulList({
      ...state.search_query,
      page: state.page,
      pagesize: state.pagesize
      // search_query: JSON.stringify(state.search_query)
    })
      .then(res => {
        dispatch({ type: "SET_LIST", payload: res.data.data });
        dispatch({ type: "SET_TOTAL", payload: res.data.total });
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    getListData();
  }, [state.page, state.pagesize, state.search_query]);

  const realColumns: ColumnProps<VulDataProps>[] = [
    {
      title: "序号",
      dataIndex: "order",
      render: (value, record, index) => index + 1
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
                setLoading(true);
                getVulDetail(record.id as number)
                  .then(res => {
                    setSelected(res.data || record);
                    handleToggleAdd();
                  })
                  .finally(() => {
                    setLoading(false);
                  });
              }}
            >
              详情
            </Button>
            <Popconfirm
              title="确定删除该条漏洞吗？"
              onCancel={() => {}}
              onConfirm={() => {
                handleDeleteVul(record.id as number);
              }}
            >
              <Button type="link">删除</Button>
            </Popconfirm>
          </div>
        );
      }
    }
  ];

  const handlePageChange = (page: number, pageSize: number | undefined) => {
    // console.log(page, pageSize);
    dispatch({ type: "SET_PAGINATION", payload: { page, pagesize: pageSize } });
  };

  const handleToggleAdd = () => {
    setAddShow(prevState => {
      // 关闭modal框清除默认值
      if (prevState) {
        setSelected(undefined);
      }
      return !prevState;
    });
  };

  const handleDeleteVul = (id: number) => {
    setLoading(true);
    deleteVul(id)
      .then(res => {
        message.success("删除成功");
        getListData();
      })
      .finally(() => {
        setLoading(false);
      });
  };
  const width = document.documentElement.offsetWidth;
  return (
    <div className="vul-manage-table">
      <div className="vul-controller">
        <Space>
          <Button
            icon={<PlusOutlined />}
            type="primary"
            onClick={handleToggleAdd}
          >
            新建描述
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
        dataSource={state.list}
        rowKey="id"
        loading={loading}
        pagination={{
          current: state.page,
          pageSize: state.pagesize,
          total: state.total,
          showTotal: (total: number) => `共${total}条`,
          showQuickJumper: true,
          showSizeChanger: true,
          onChange: handlePageChange
        }}
      />
      {addShow && (
        <VulModal
          visible={addShow}
          onCancel={handleToggleAdd}
          title={!!selected ? "漏洞描述详情" : "新建漏洞描述"}
          width={width}
          wrapClassName="add-vul-wrap"
          selected={selected}
          type="vul"
        />
      )}
    </div>
  );
};

export default VulTable;
