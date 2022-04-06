import React from "react";
import { faGolang, faJava, faJs, faPython, faRust } from "@fortawesome/free-brands-svg-icons";
import { faGear, faTerminal } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon, FontAwesomeIconProps } from "@fortawesome/react-fontawesome";

interface LanguageIconProps extends Pick<FontAwesomeIconProps, "size"> {
  language: string;
}

const LanguageIcon = ({ language, size = "lg", ...props }: LanguageIconProps) => {
  let icon;
  let title;
  switch (language) {
    case "python":
      icon = faPython;
      title = "Python";
      break;
    case "go":
      icon = faGolang;
      title = "Golang";
      break;
    case "bash":
      icon = faTerminal;
      title = "Shell";
      break;
    case "javascript":
      icon = faJs;
      title = "JavaScript";
      break;
    case "java":
      icon = faJava;
      title = "Java";
      break;
    case "rust":
      icon = faRust;
      title = "Rust";
      break;
    default:
      icon = faGear;
      title = "None detected";
  }

  return <FontAwesomeIcon title={title} icon={icon} size={size} {...props} />;
};

export default LanguageIcon;
