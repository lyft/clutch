import * as React from "react";
import type { FieldValues, UseFormReturn } from "react-hook-form";
import { useForm } from "react-hook-form";
import { DevTool } from "@hookform/devtools";
import { yupResolver } from "@hookform/resolvers/yup";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import {
  alpha,
  Grid as MuiGrid,
  StandardTextFieldProps,
  Table as MuiTable,
  TableBody as MuiTableBody,
  TableCell as MuiTableCell,
  TableContainer as MuiTableContainer,
  TableRow,
  Theme,
} from "@mui/material";
import _ from "lodash";
import type { BaseSchema } from "yup";
import { object } from "yup";

import { useWizardContext } from "../Contexts";
import type { NoteConfig } from "../Feedback";
import { Tooltip } from "../Feedback/tooltip";
import { TextField } from "../Input/text-field";
import styled from "../styled";

interface RowData {
  input?: {
    key?: string;
    type?: string;
    validation?: BaseSchema<unknown>;
    warning?: NoteConfig;
  };
  textFieldLabels?: {
    disabledField: string;
    updatedField: string;
  };
  name: string;
  value: unknown;
  disabledFieldTooltip?: string;
}

interface IdentifiableRowData extends RowData {
  id: string;
}

const TableContainer = styled(MuiTableContainer)<{
  $maxHeight?: MetadataTableProps["maxHeight"];
}>(
  {
    borderWidth: "0",
    border: "0",
  },
  props => ({
    maxHeight: props.$maxHeight || "fit-content",
  })
);

const Table = styled(MuiTable)(({ theme }: { theme: Theme }) => ({
  border: `1px solid ${alpha(theme.palette.secondary[900], 0.12)}`,
  borderRadius: "4px",
  borderCollapse: "unset",
}));

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

const TableCell = styled(MuiTableCell)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.secondary[900],
  fontSize: "14px",
  fontWeight: "normal",
  height: "48px",
  padding: "8px 16px",
}));

const Grid = styled(MuiGrid)({
  display: "flex",
  alignItems: "center",
  ".MuiFormControl-root": {
    flexDirection: "row",
  },
  ".MuiFormControl-root .MuiInputBase-root": {
    height: "40px",
    width: "100px",
    alignSelf: "center",
  },
  ".MuiFormControl-root .MuiFormHelperText-root": {
    flex: 1,
  },
  ".textfield-disabled .MuiInput-input": {
    padding: "0px",
    textAlign: "center",
  },
});

const KeyCellContainer = styled(TableCell)(({ theme }: { theme: Theme }) => ({
  width: "45%",
  background: alpha(theme.palette.secondary[900], 0.03),
  fontWeight: 500,
}));

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
  validation: UseFormReturn<FieldValues, object>;
}

// TODO (maybe): instead of a disabled text field and editable text field, remove disabled field and have a reset icon next to text field
// to reset field to the default value
const MutableRow: React.FC<MutableRowProps> = ({ data, onUpdate, onReturn, validation }) => {
  const error = validation?.formState?.errors?.[data.name];
  const { warning } = data.input;

  // intercept the update callback to prevent updates if there are form errors present
  // based on the validation.
  const updateCallback = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) =>
    error ? () => {} : onUpdate(e);

  // get helper text in case we need info or success cases
  const getHelperText = (): String => error?.message || warning?.text || "";

  // get the color if there is an error or warning
  // adding more cases if we need info or success cases
  const getTextFieldColor = (): StandardTextFieldProps["color"] | undefined => {
    if (error) {
      return "error";
    }
    if (warning) {
      return "warning";
    }
    return undefined;
  };

  const disabledTextFieldComponent = (
    <TextField
      disabled
      id={data.id}
      name={data.name}
      defaultValue={data.value}
      label={data.textFieldLabels?.disabledField}
    />
  );

  return (
    <TableRow key={data.id}>
      <KeyCell data={data} />
      <TableCell>
        <Grid>
          <div className="textfield-disabled">
            {/* // In the case where a disabledFieldTooltip is not provided, the value itself will be the tooltip */}
            <Tooltip title={data.disabledFieldTooltip ?? data.value}>
              {disabledTextFieldComponent}
            </Tooltip>
          </div>
          <ChevronRightIcon />
          <TextField
            id={data.id}
            name={data.name}
            label={data.textFieldLabels?.updatedField}
            defaultValue={data.value}
            type={data?.input?.type}
            onChange={updateCallback}
            onReturn={onReturn}
            onFocus={updateCallback}
            helperText={getHelperText()}
            error={!!error || false}
            color={getTextFieldColor()}
            formRegistration={validation.register}
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
    <TableContainer $maxHeight={maxHeight}>
      {process.env.REACT_APP_DEBUG_FORMS && onUpdate !== undefined && <DevTool control={control} />}
      <Table>
        <TableBody>
          {rows.map((row: IdentifiableRowData) => {
            return row.input !== undefined && onUpdate ? (
              <MutableRow
                data={row}
                onUpdate={e => {
                  onUpdate(e.target.id, e.target.value);
                }}
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
