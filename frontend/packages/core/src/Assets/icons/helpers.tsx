import { styled } from "../../Utils";

const StyledSvg: React.ElementType = styled("svg")((props: { hoverFill: string }) => ({
  path: {
    "&:hover": {
      fill: props.hoverFill,
    },
  },
}));

export default StyledSvg;
