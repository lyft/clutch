import React from "react";
import type { FieldValues, UseFormRegister } from "react-hook-form";
import { Select, TextField } from "@clutch-sh/core";

type FieldType = "title" | "text" | "number" | "select";

interface FormProps {
  state: any;
  items: FormItem[];
  register: UseFormRegister<FieldValues>;
  errors: any;
}

interface TextFieldProps {
  defaultValue: string | undefined;
}

interface SelectProps {
  options: SelectOption[];
  defaultValue: string;
}

interface SelectOption {
  label: string;
  value: string;
}

export interface FormItem {
  name?: string;
  label: string;
  type: FieldType;
  validation?: any;
  visible?: boolean;
  inputProps?: SelectProps | TextFieldProps;
}

const FormFields: React.FC<FormProps> = ({ state, items, register, errors }) => {
  const [data, setData] = state;

  return (
    <>
      {items.map(field => {
        if (field.type === "title") {
          return (
            <h3 key={field.label} style={{ textAlign: "center" }}>
              {field.label}
            </h3>
          );
        }
        if (["text", "number"].indexOf(field.type) >= 0) {
          const customProps: TextFieldProps = field.inputProps as TextFieldProps;

          return (
            <TextField
              key={field.label}
              name={field.name}
              label={field.label}
              defaultValue={customProps.defaultValue}
              type={field.type}
              onChange={e => {
                const copiedData = { ...data };
                copiedData[field.name] = e.target.value;
                setData(copiedData);
              }}
              error={!!errors[field.name]}
              helperText={errors[field.name] ? errors[field.name].message : ""}
              formRegistration={register}
            />
          );
        }
        if (field.type === "select") {
          const customProps: SelectProps = field.inputProps as SelectProps;
          return (
            <Select
              key={field.label}
              name={field.name}
              label={field.label}
              options={customProps.options}
              defaultOption={customProps.options
                .map(o => o.value)
                .indexOf(customProps.defaultValue)}
              onChange={value => {
                const copiedData = { ...data };
                copiedData[field.name] = value;
                setData(copiedData);
              }}
            />
          );
        }

        return <div key="blank" />;
      })}
    </>
  );
};

export default FormFields;
