import React from "react";
import type { ButtonProps as MuiButtonProps, GridJustification } from "@material-ui/core";
import { Button as MuiButton, emphasize, fade, Grid } from "@material-ui/core";
import styled from "styled-components";

const StyledButton = styled(MuiButton)`
  ${({ theme, ...props }) => `
  margin: 15px;
  background-color: ${
    props["data-destructive"] === "true"
      ? theme.palette.destructive.main
      : fade(theme.palette.secondary.main, 0.3)
  };
  &:hover {
    background-color: ${emphasize(theme.palette.secondary.main, 0.3)};
  };
  `}
`;

interface ButtonProps extends MuiButtonProps {
  text: string;
  destructive?: boolean;
}

const Button: React.FC<ButtonProps> = ({ text, destructive = false, ...props }) => (
  <StyledButton variant="contained" data-destructive={destructive.toString()} {...props}>
    {text}
  </StyledButton>
);

const AdvanceButton: React.FC<ButtonProps> = ({ text, ...props }) => (
  <Button destructive={false} text={text} {...props} />
);
const DestructiveButton: React.FC<ButtonProps> = ({ text, ...props }) => (
  <Button destructive text={text} {...props} />
);

interface ButtonGroupProps {
  buttons: ButtonProps[];
  justify?: GridJustification;
}

const ButtonGroup: React.FC<ButtonGroupProps> = ({ buttons, justify = "center" }) => (
  <Grid container justify={justify}>
    {buttons.map(button => (
      <Button key={button.text} {...button} />
    ))}
  </Grid>
);

export { AdvanceButton, Button, ButtonGroup, DestructiveButton };
