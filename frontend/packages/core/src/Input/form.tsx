import styled from "@emotion/styled";

const Form = styled.form({
  width: "inherit",
  "> *": {
    margin: "8px 0",
  },
});

const FormRow = styled.div({
  display: "flex",
  "> *": {
    margin: "0 8px",
  },
  "> *:first-child": {
    margin: "0 8px 0 0",
  },
  "> *:last-child": {
    margin: "0 0 0 8px",
  },
  /**
   * https://mui.com/material-ui/react-text-field/#helper-text
   * This is used to align items since text fields need an empty help text
   * to stay aligned
   */
  "> *:not(.MuiFormControl-root):not(.MuiAutocomplete-root)": {
    marginTop: "-23px",
  },
});

export { Form, FormRow };
