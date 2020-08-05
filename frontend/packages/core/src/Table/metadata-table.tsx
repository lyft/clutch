import React from "react";
import { useForm } from "react-hook-form";
import { DevTool } from "@hookform/devtools";
import { yupResolver } from '@hookform/resolvers';
import { Table, TableBody, TableCell, TableContainer, TableRow } from "@material-ui/core";
import type { Size } from "@material-ui/core/TableCell";
import _ from "lodash";
import type { Schema } from "yup";
import { object } from "yup";

import { useWizardContext } from "../Contexts";
import TextField from "../Input/text-field";

interface RowData {
  input?: {
    key?: string;
    type?: string;
    validation?: Schema<unknown>;
  };
  name: string;
  value: unknown;
}

interface IdentifiableRowData extends RowData {
  id: string;
}

interface KeyCellProps {
  data: IdentifiableRowData;
  size: Size;
}

const KeyCell: React.FC<KeyCellProps> = ({ data, size }) => {
  let { name } = data;
  if (data.value instanceof Array && data.value.length > 1) {
    name = `${data.name}s`;
  }
  return (
    <TableCell size={size}>
      <strong>{name}</strong>
    </TableCell>
  );
};

interface ViewOnlyRowProps {
  data: IdentifiableRowData;
  size: Size;
}

const ViewOnlyRow: React.FC<ViewOnlyRowProps> = ({ data, size }) => {
  let { value } = data;
  if (data.value instanceof Array && data.value.length > 1) {
    value = data.value.join(", ");
  }
  return (
    <TableRow key={data.id}>
      <KeyCell data={data} size={size} />
      <TableCell size={size}>{value}</TableCell>
    </TableRow>
  );
};

interface EditableRowProps {
  data: IdentifiableRowData;
  onUpdate: (event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => void;
  validation: any;
  size: Size;
}

const EditableRow: React.FC<EditableRowProps> = ({ data, onUpdate, validation, size }) => {
  const error = validation.errors?.[data.name];

  return (
    <TableRow key={data.id}>
      <KeyCell data={data} size={size} />
      <TableCell size={size}>
        <TextField
          id={data.id}
          name={data.name}
          defaultValue={data.value}
          size="small"
          type={data?.input?.type}
          InputProps={{ margin: "dense", color: "secondary", name: data.name }}
          inputProps={data?.input}
          onChange={onUpdate}
          onFocus={onUpdate}
          inputRef={validation.register}
          helperText={error?.message || ""}
          error={!!error || false}
          maxWidth="50%"
        />
      </TableCell>
    </TableRow>
  );
};

const Wrapper: React.FC<{ editable: boolean }> = ({ editable, children }) =>
  editable ? <form>{children}</form> : <>{children}</>;

interface MetadataTableProps {
  data: RowData[];
  onUpdate?: (id: string, value: unknown) => void;
  variant?: Size;
}

const MetadataTable: React.FC<MetadataTableProps> = ({
  data,
  onUpdate,
  children,
  variant,
  ...props
}) => {
  const displayVariant = variant || "medium";
  const { setOnSubmit } = useWizardContext();
  let rows = data;
  if (_.isEmpty(data)) {
    rows = [{ name: "Error", value: "No Data Available" }];
  } else {
    rows = data.map(row => {
      const id = row?.input?.key || _.camelCase(row.name);
      return { ...row, id };
    });
  }

  const validationShape = {};
  rows.forEach(row => {
    if (row?.input?.validation) {
      validationShape[row.name] = row.input.validation;
    }
  });
  const validation = useForm({
    mode: "onChange",
    resolver: yupResolver(object().shape(validationShape)),
  });
  const { control, handleSubmit } = validation;

  React.useEffect(() => {
    setOnSubmit(handleSubmit);
  }, []);

  return (
    <TableContainer {...props}>
      {process.env.REACT_APP_DEBUG_FORMS && onUpdate !== undefined && <DevTool control={control} />}
      <Wrapper editable={onUpdate !== undefined}>
        <Table {...props}>
          <TableBody>
            {rows.map((row: IdentifiableRowData) => {
              return row.input !== undefined && onUpdate ? (
                <EditableRow
                  data={row}
                  onUpdate={e => onUpdate(e.target.id, e.target.value)}
                  key={row.id}
                  validation={validation}
                  size={displayVariant}
                />
              ) : (
                <ViewOnlyRow data={row} key={row.id} size={displayVariant} />
              );
            })}
            {children}
          </TableBody>
        </Table>
      </Wrapper>
    </TableContainer>
  );
};

export default MetadataTable;
