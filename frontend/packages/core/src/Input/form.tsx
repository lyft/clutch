import * as React from "react";
import styled from "@emotion/styled";

const StyledForm = styled.form(
  {
    width: "inherit",
    display: "flex",
  },
  props => ({
    "> *": {
      margin: props["data-direction"] === "row" ? "0 8px" : "8px 0",
    },
    flexDirection: props["data-direction"],
  })
);

export interface FormProps extends React.FormHTMLAttributes<HTMLFormElement> {
  direction?: "row" | "column";
}

const Form = ({ children, direction = "column", ...props }: FormProps) => (
  <StyledForm data-direction={direction} {...props}>
    {children}
  </StyledForm>
);
export default Form;
