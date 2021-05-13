import React from "react";
import { Input } from "antd";
import { MinusCircleOutlined, PlusCircleOutlined } from "@ant-design/icons/lib";
import { InputProps } from "antd/es/input";

export interface InputGroupProps extends InputProps {
  showDeleteButton?: boolean;
  showAddButton?: boolean;
  value?: any;
  onChange?: any;
  handleAdd?: any;
  handleDelete?: any;
  handleBlur?: any;
  id?: string;
}

const InputGroup: React.FC<InputGroupProps> = props => {
  const handleChange = (val: any) => {
    props.onChange?.({ ...props.value, ...val });
  };

  return (
    <div className={props.className}>
      <Input.Group compact style={{ width: "92%", float: "left" }}>
        <Input
          style={{ width: "20%" }}
          name="key"
          value={props?.value?.key}
          onChange={e => {
            handleChange({ key: e.target.value });
          }}
          onBlur={props.handleBlur}
        />
        <Input
          style={{ width: "80%" }}
          name="value"
          value={props?.value?.value}
          onChange={e => {
            handleChange({ value: e.target.value });
          }}
          onBlur={props.handleBlur}
        />
      </Input.Group>
      <PlusCircleOutlined
        className="dynamic-delete-button color-grey"
        type="plus-circle-o"
        onClick={props.handleAdd}
      />
      {props.showDeleteButton && (
        <MinusCircleOutlined
          className="dynamic-delete-button color-red"
          type="minus-circle-o"
          style={{ color: "#ff6868" }}
          onClick={() => props.handleDelete?.(props.value.id)}
        />
      )}
    </div>
  );
};

InputGroup.defaultProps = {
  showAddButton: true,
  showDeleteButton: true
};

export default InputGroup;
