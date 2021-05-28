import React from "react";
import { Select, TextField } from "@clutch-sh/core";

interface FormFieldsProps {
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
  defaultValue: string;
}

interface SelectOption {
  label: string;
  value: string;
}

interface FormItem {
  name: string;
  label: string;
  type: string;
  inputProps: SelectProps | TextFieldProps;
}

const FormFields: React.FC<FormFieldsProps> = ({ state, items, errors, register }) => {
  const [data, setData] = state;

  return (
    <>
      {items.map(field => {
        if (["text", "multiline-text", "number"].indexOf(field.type) >= 0) {
          const props: TextFieldProps = field.inputProps as TextFieldProps;
          return (
            <TextField
              multiline={field.type === "multiline-text"}
              key={field.name}
              name={field.name}
              label={field.label}
              defaultValue={props.defaultValue}
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
        const props: SelectProps = field.inputProps as SelectProps;
        return (
          <Select
            defaultOption={props.options.map(o => o.value).indexOf(props.defaultValue)}
            name={field.name}
            key={field.name}
            label={field.label}
            options={props.options}
            onChange={value => {
              const copiedData = { ...data };
              copiedData[field.name] = value;
              setData(copiedData);
            }}
          />
        );
      })}
    </>
  );
};

export default FormFields;
