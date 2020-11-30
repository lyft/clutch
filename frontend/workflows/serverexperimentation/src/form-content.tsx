import React from "react";
import { RadioGroup, Select, TextField } from "@clutch-sh/core";
import styled from "styled-components";

const StyledDiv = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  width: 100%;
`;

interface FormProps {
  state: any;
  items: FormItem[];
  register: any;
  errors: any;
}

interface TextFieldProps {
  defaultValue: string | undefined;
}

interface SelectProps {
  options: SelectOption[];
}

interface SelectOption {
  label: string;
  value: string;
}

interface RadioGroupProps {
  options: RadioGroupOption[];
}

interface RadioGroupOption {
  label: string;
  value: string;
}

interface FormItem {
  name: string;
  label: string;
  type: string;
  inputProps: SelectProps | TextFieldProps;
}

const FormContent: React.FC<FormProps> = ({ state, items, register, errors }) => {
  const [data, setData] = state;

  return (
    <StyledDiv>
      {items.map(field => {
        if (["text", "number"].indexOf(field.type) >= 0) {
          const customProps: TextFieldProps = field.inputProps as TextFieldProps;
          return (
            <TextField
              key={field.name}
              name={field.name}
              label={field.label}
              defaultValue={customProps.defaultValue}
              type={field.type}
              onChange={e => {
                const copiedData = { ...data };
                copiedData[field.name] = e.target.value;
                setData(copiedData);
              }}
              inputRef={register}
              error={!!errors[field.name]}
              helperText={errors[field.name] ? errors[field.name].message : ""}
            />
          );
        }
        if (field.type === "radio-group") {
          const customProps: RadioGroupProps = field.inputProps as RadioGroupProps;
          return (
            <RadioGroup
              key={field.name}
              name={field.name}
              label={field.label}
              options={customProps.options}
              onChange={value => {
                const copiedData = { ...data };
                copiedData[field.name] = value;
                setData(copiedData);
              }}
            />
          );
        }
        if (field.type === "select") {
          const customProps: SelectProps = field.inputProps as SelectProps;
          return (
            <div key={field.name} style={{ margin: "15px 0" }}>
              <Select
                name={field.name}
                key={field.name}
                label={field.label}
                options={customProps.options}
                onChange={value => {
                  const copiedData = { ...data };
                  copiedData[field.name] = value;
                  setData(copiedData);
                }}
              />
            </div>
          );
        }

        return <div key="blank" />;
      })}
    </StyledDiv>
  );
};

export { FormContent, FormItem, SelectOption, SelectProps, TextFieldProps };
