import type { TypographyOptions } from "@material-ui/core/styles/createTypography";

const REGULAR = 400;
const MEDIUM = 500;
const BOLD = 700;

const Typography: TypographyOptions = {
  fontFamily: "'Roboto', 'Helvetica', 'Arial', sans-serif",
  h1: {
    fontWeight: MEDIUM,
    fontSize: "36px",
    lineHeight: "44px",
  },
  h2: {
    fontWeight: BOLD,
    fontSize: "26px",
    lineHeight: "32px",
  },
  h3: {
    fontWeight: BOLD,
    fontSize: "22px",
    lineHeight: "28px",
  },
  h4: {
    fontWeight: BOLD,
    fontSize: "20px",
    lineHeight: "24px",
  },
  h5: {
    fontWeight: BOLD,
    fontSize: "16px",
    lineHeight: "20px",
  },
  h6: {
    fontWeight: BOLD,
    fontSize: "14px",
    lineHeight: "24px",
  },
  subtitle1: {
    fontWeight: MEDIUM,
    fontSize: "20px",
    lineHeight: "24px",
  },
  subtitle2: {
    fontWeight: MEDIUM,
    fontSize: "16px",
    lineHeight: "20px",
  },
  body1: {
    fontWeight: REGULAR,
    fontSize: "20px",
    lineHeight: "26px",
  },
  body2: {
    fontWeight: REGULAR,
    fontSize: "16px",
    lineHeight: "22px",
  },
  // TODO: add button, caption, and overline
  // button: {
  //   fontWeight: ,
  //   fontSize: "",
  //   lineHeight: ,
  // },
  // caption: {
  //   fontWeight: ,
  //   fontSize: "",
  //   lineHeight: ,
  // },
  // overline: {
  //   fontWeight: ,
  //   fontSize: "",
  //   lineHeight: ,
  // },
};
// TODO: remaining styles to add
// const STYLE_MAP = {
//   h1: {
//     size: "36",
//     weight: MEDIUM,
//     lineHeight: "44",
//   },
//   h2: {
//     size: "26",
//     weight: BOLD,
//     lineHeight: "32",
//   },
//   h3: {
//     size: "22",
//     weight: BOLD,
//     lineHeight: "28",
//   },
//   h4: {
//     size: "20",
//     weight: BOLD,
//     lineHeight: "24",
//   },
//   h5: {
//     size: "16",
//     weight: BOLD,
//     lineHeight: "20",
//   },
//   h6: {
//     size: "14",
//     weight: BOLD,
//     lineHeight: "18",
//   },
//   subtitle1: {
//     size: "20",
//     weight: MEDIUM,
//     lineHeight: "24",
//   },
//   subtitle2: {
//     size: "16",
//     weight: MEDIUM,
//     lineHeight: "20",
//   },
//   subtitle3: {
//     size: "14",
//     weight: MEDIUM,
//     lineHeight: "18",
//   },
//   body1: {
//     size: "20",
//     weight: REGULAR,
//     lineHeight: "26",
//   },
//   body2: {
//     size: "16",
//     weight: REGULAR,
//     lineHeight: "22",
//   },
//   body3: {
//     size: "14",
//     weight: REGULAR,
//     lineHeight: "18",
//   },
//   body4: {
//     size: "12",
//     weight: REGULAR,
//     lineHeight: "16",
//   },
//   caption1: {
//     size: "16",
//     weight: BOLD,
//     lineHeight: "20",
//     props: {
//       textTransform: "uppercase",
//     },
//   },
//   caption2: {
//     size: "12",
//     weight: BOLD,
//     lineHeight: "16",
//     props: {
//       textTransform: "uppercase",
//     },
//   },
//   overline: {
//     size: "10",
//     weight: REGULAR,
//     lineHeight: "10",
//     props: {
//       textTransform: "uppercase",
//       letterSpacing: "1.5px",
//     },
//   },
//   input: {
//     size: "14",
//     weight: REGULAR,
//     lineHeight: "18",
//   },
// } as StyleMapProps;

export default Typography;
