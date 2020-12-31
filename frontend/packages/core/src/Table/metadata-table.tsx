import * as React from "react";
import { useForm } from "react-hook-form";
import styled from "@emotion/styled";
import { DevTool } from "@hookform/devtools";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  Grid as MuiGrid,
  Table as MuiTable,
  TableBody as MuiTableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableRow,
} from "@material-ui/core";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import _ from "lodash";
import type { Schema } from "yup";
import { object } from "yup";

import { useWizardContext } from "../Contexts";
import { TextField } from "../Input/text-field";

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

const TableContainer = styled(MuiTableContainer)(
  {
    borderWidth: "0",
    border: "0",
  },
  props => ({
    maxHeight: props["data-max-height"] || "fit-content",
  })
);

const Table = styled(MuiTable)({
  border: "1px solid rgba(13, 16, 48, 0.12)",
  borderRadius: "4px",
  borderCollapse: "unset",
});

const TableBody = styled(MuiTableBody)({
  "tr:first-of-type > td:first-of-type": {
    borderTopLeftRadius: "3px",
  },
  "tr:first-of-type > td:last-of-type": {
    borderTopRightRadius: "3px",
  },
  "tr:last-of-type > td": {
    borderBottom: "0",
  },
  "tr:last-of-type > td:first-of-type": {
    borderBottomLeftRadius: "3px",
  },
  "tr:last-of-type > td:last-of-type": {
    borderBottomRightRadius: "3px",
  },
});

const TableCell = styled(MuiTableCell)({
  color: "#0D1030",
  fontSize: "14px",
  fontWeight: "normal",
  height: "48px",
  padding: "0 16px",
});

const Grid = styled(MuiGrid)({
  ".MuiFormControl-root": {
    height: "40px",
    width: "100px",
    flexDirection: "row",
  },
  ".textfield-disabled .MuiFormControl-root": {
    width: "41px",
  },
  ".textfield-disabled .MuiInput-input": {
    padding: "0px",
    textAlign: "center",
  },
});

const KeyCellContainer = styled(TableCell)({
  width: "45%",
  background: "rgba(13, 16, 48, 0.03)",
  fontWeight: 500,
});

interface KeyCellProps {
  data: IdentifiableRowData;
}

const KeyCell: React.FC<KeyCellProps> = ({ data }) => {
  let { name } = data;
  if (data.value instanceof Array && data.value.length > 1) {
    name = `${data.name}s`;
  }
  return <KeyCellContainer>{name}</KeyCellContainer>;
};

interface ImmutableRowProps extends KeyCellProps {}

const ImmutableRow: React.FC<ImmutableRowProps> = ({ data }) => {
  let { value } = data;
  if (data.value instanceof Array && data.value.length > 1) {
    value = data.value.join(", ");
  }
  return (
    <TableRow key={data.id}>
      <KeyCell data={data} />
      <TableCell>{value}</TableCell>
    </TableRow>
  );
};

interface MutableRowProps extends ImmutableRowProps {
  onUpdate: (event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => void;
  onReturn: () => void;
  validation: any;
}

const MutableRow: React.FC<MutableRowProps> = ({ data, onUpdate, onReturn, validation }) => {
  const error = validation.errors?.[data.name];

  return (
    <TableRow key={data.id}>
      <KeyCell data={data} />
      <TableCell>
        <Grid container alignItems="center" wrap="nowrap" spacing={2}>
          <div className="textfield-disabled">
            <TextField disabled id={data.id} name={data.name} defaultValue={data.value} />
          </div>
          <ChevronRightIcon />
          <TextField
            id={data.id}
            name={data.name}
            defaultValue={data.value}
            type={data?.input?.type}
            onChange={onUpdate}
            onReturn={onReturn}
            onFocus={onUpdate}
            inputRef={validation.register}
            helperText={error?.message || ""}
            error={!!error || false}
          />
        </Grid>
      </TableCell>
    </TableRow>
  );
};

export interface MetadataTableProps {
  data: RowData[];
  onUpdate?: (id: string, value: unknown) => void;
  maxHeight?: string;
}

export const MetadataTable: React.FC<MetadataTableProps> = ({
  data,
  onUpdate,
  maxHeight,
  children,
}) => {
  const { onSubmit, setOnSubmit } = useWizardContext();
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
    <TableContainer data-max-height={maxHeight}>
      {process.env.REACT_APP_DEBUG_FORMS && onUpdate !== undefined && <DevTool control={control} />}
      <Table>
        <TableBody>
          {rows.map((row: IdentifiableRowData) => {
            return row.input !== undefined && onUpdate ? (
              <MutableRow
                data={row}
                onUpdate={e => onUpdate(e.target.id, e.target.value)}
                onReturn={onSubmit}
                key={row.id}
                validation={validation}
              />
            ) : (
              <ImmutableRow data={row} key={row.id} />
            );
          })}
          {children}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default MetadataTable;
