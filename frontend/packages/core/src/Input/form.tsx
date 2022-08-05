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
});

export { Form, FormRow };
