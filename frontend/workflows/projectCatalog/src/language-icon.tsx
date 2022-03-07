import React from "react";
import { faGolang, faJava, faJs, faPython } from "@fortawesome/free-brands-svg-icons";
import { faGear, faTerminal } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon, FontAwesomeIconProps } from "@fortawesome/react-fontawesome";

interface LanguageIconProps extends Pick<FontAwesomeIconProps, "size"> {
  language: string;
}

const LanguageIcon = ({ language, size = "lg", ...props }: LanguageIconProps) => {
  let icon;
  switch (language) {
    case "python":
      icon = faPython;
      break;
    case "go":
      icon = faGolang;
      break;
    case "bash":
      icon = faTerminal;
      break;
    case "javascript":
      icon = faJs;
      break;
    case "java":
      icon = faJava;
      break;
    default:
      icon = faGear;
  }

  return <FontAwesomeIcon icon={icon} size={size} {...props} />;
};

export default LanguageIcon;
