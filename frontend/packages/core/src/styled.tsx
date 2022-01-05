import type { CreateStyled, StyledOptions } from "@emotion/styled";
import emotion from "@emotion/styled";

const transientOptions: Parameters<CreateStyled>[1] = {
  shouldForwardProp: (propName: string) => !propName.startsWith("$"),
};

const styled = (tag: any, options?: StyledOptions<any>) => {
  return emotion(tag, {...transientOptions, ...(options || {})});
};

export default styled;
