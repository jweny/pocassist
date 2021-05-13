import React, { useCallback, useEffect, useState } from "react";
import { Button, Checkbox, Col, Modal, Row } from "antd";
import "./index.less";
import { ColumnProps } from "antd/es/table";
import { CheckboxChangeEvent } from "antd/es/checkbox";

interface CustomColumnsProps {
  className?: string;
  allColumns: ColumnProps<any>[];
  columns: ColumnProps<any>[];
  setColumns: React.Dispatch<React.SetStateAction<ColumnProps<any>[]>>;
}
const CustomColumns: React.FC<CustomColumnsProps> = props => {
  // console.log(props);
  const { className, columns, setColumns, allColumns } = props;
  const [showModal, setShowModal] = useState<boolean>(false);
  const [value, setValue] = useState<string[]>([]);
  const [checkAllState, setCheckAllState] = useState({
    check: false,
    indeterminate: true
  });

  const resetValue = useCallback(() => {
    setValue(columns.map(item => item.dataIndex as string));
    setCheckAllState({
      check: columns.length === allColumns.length,
      indeterminate: columns.length > 0 && columns.length < allColumns.length
    });
  }, [columns, allColumns]);

  useEffect(() => {
    resetValue();
  }, [resetValue]);

  const handleToggleShow = () => {
    setShowModal(prevState => !prevState);
  };

  const handleCancel = () => {
    setShowModal(false);
    resetValue();
  };

  const handleFinish = () => {
    // console.log(value);
    const returnColumns = allColumns.filter(
      item => value.indexOf(item.dataIndex as string) !== -1
    );
    setColumns(returnColumns);
    handleToggleShow();
  };

  const handleCheckAllChange = (e: CheckboxChangeEvent) => {
    const curVal = e.target.checked;

    setCheckAllState({
      indeterminate: false,
      check: curVal
    });
    setValue(curVal ? allColumns.map(item => item.dataIndex as string) : []);
  };

  return (
    <div className={className + " custom-columns-wrap"}>
      <Button type="link" onClick={handleToggleShow}>
        {/*<AppstoreOutlined onClick={handleToggleShow} />*/}
        自定义
      </Button>
      <Modal
        visible={showModal}
        title="自定义表头"
        onCancel={handleCancel}
        onOk={handleFinish}
        width={800}
        wrapClassName="custom-columns-modal"
      >
        <div className="site-checkbox-all-wrapper">
          <Checkbox
            indeterminate={checkAllState.indeterminate}
            onChange={handleCheckAllChange}
            checked={checkAllState.check}
          >
            全选
          </Checkbox>
        </div>
        <Checkbox.Group
          name="columns"
          value={value}
          onChange={val => {
            setValue(val as string[]);
            setCheckAllState({
              check: val.length === allColumns.length,
              indeterminate: val.length > 0 && val.length < allColumns.length
            });
          }}
          style={{ width: "100%" }}
        >
          <Row>
            {allColumns.map(item => {
              return (
                <Col span={6} key={item.dataIndex as string}>
                  <Checkbox value={item.dataIndex}>{item.title}</Checkbox>
                </Col>
              );
            })}
          </Row>
        </Checkbox.Group>
      </Modal>
    </div>
  );
};

export default CustomColumns;
