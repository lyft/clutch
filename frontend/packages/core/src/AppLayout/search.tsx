import React from "react";
import { useNavigate } from "react-router-dom";
import styled from "@emotion/styled";
import {
  Box,
  ClickAwayListener,
  Grid,
  Icon,
  IconButton,
  Popper as MuiPopper,
  TextField,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import SearchIcon from "@material-ui/icons/Search";
import type { AutocompleteRenderInputParams } from "@material-ui/lab/Autocomplete";
import Autocomplete from "@material-ui/lab/Autocomplete";
import type { FilterOptionsState } from "@material-ui/lab/useAutocomplete";
import _ from "lodash";

import { useAppContext } from "../Contexts";

import type { SearchIndex } from "./utils";
import { searchIndexes } from "./utils";

const hotKey = "/";

const Container = styled(Box)({
  borderRadius: "4px",
  marginRight: "8px",
  backgroundColor: "#ffffff",
});

const InputField = styled(TextField)({
  maxWidth: "504px",
  minWidth: "504px",
  "@media screen and (max-width: 800px)": {
    minWidth: "125px",
  },
  ".MuiAutocomplete-input": {
    color: "#0d1030",
  },
});

const ResultLabel = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontWeight: 500,
  fontSize: "14px",
});

const SearchIconButton = styled(IconButton)({
  color: "#ffffff",
  padding: "8px",
  "&:hover": {
    background: "#2d3db4",
  },
  "&:active": {
    background: "#2938a5",
  },
});

const StyledCloseIcon = styled(Icon)({
  color: "#0c0b31",
});

const Popper = styled(MuiPopper)({
  ".MuiAutocomplete-option": {
    "&:hover": {
      background: "linear-gradient(0deg, rgba(53, 72, 212, 0.1), rgba(53, 72, 212, 0.1)), #ffffff",
    },
  },
});

const renderPopper = props => {
  return <Popper {...props} />;
};

const CustomCloseIcon: React.FC = () => {
  return (
    <StyledCloseIcon>
      <CloseIcon fontSize="small" />
    </StyledCloseIcon>
  );
};

const Input = (params: AutocompleteRenderInputParams): React.ReactNode => {
  const searchRef = React.useRef();
  const handleKeyPress = (event: KeyboardEvent) => {
    if (searchRef.current !== undefined) {
      if (event.key === hotKey && (event.target as Node).nodeName !== "INPUT") {
        // @ts-ignore
        searchRef.current.focus();
      } else if (event.key === "Escape") {
        // @ts-ignore
        searchRef.current.blur();
      }
    }
  };

  React.useLayoutEffect(() => {
    window.addEventListener("keydown", handleKeyPress);
  }, []);

  return (
    <InputField
      {...params}
      placeholder="Search..."
      variant="outlined"
      fullWidth
      inputRef={searchRef}
    />
  );
};

interface ResultProps {
  option: SearchIndex;
  handleSelection: () => void;
}

const Result: React.FC<ResultProps> = ({ option, handleSelection }) => (
  <Grid container alignItems="center" onClick={handleSelection}>
    <Grid item xs>
      <ResultLabel>{option.label}</ResultLabel>
    </Grid>
  </Grid>
);

const filterResults = (searchOptions: SearchIndex[], state: FilterOptionsState<SearchIndex>) => {
  return _.filter(searchOptions, o => {
    return o.label.toLowerCase().includes(state.inputValue.toLowerCase());
  });
};

const SearchField: React.FC = () => {
  const { workflows } = useAppContext();
  const navigate = useNavigate();
  const options = searchIndexes(workflows);
  const [inputValue, setInputValue] = React.useState("");
  const [showOptions, setShowOptions] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  const renderResult = (option: SearchIndex) => {
    const handleSelection = () => {
      navigate(option.path);
    };

    return <Result option={option} handleSelection={handleSelection} />;
  };

  const onInputChange = (__: React.ChangeEvent<{}>, value: string) => {
    if (value === "") {
      setShowOptions(false);
      setInputValue("");
    } else if (value !== hotKey) {
      setShowOptions(true);
      setInputValue(value);
    }
    const option = _.find(options, o => {
      return o.label === value;
    });
    if (option !== undefined) {
      setShowOptions(false);
      setInputValue("");
    }
  };

  const onOptionsOpen = () => {
    if (inputValue !== "") {
      setShowOptions(true);
    }
  };

  const onOptionsClose = (event: any) => {
    setShowOptions(false);
    const option = _.find(options, o => {
      return o.label === event.target.value;
    });
    if (option !== undefined) {
      navigate(option.path);
    }
    setInputValue("");
  };

  const handleOpen = () => {
    setOpen(!open);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Grid container alignItems="center">
      {open ? (
        <Container>
          <ClickAwayListener onClickAway={handleClose}>
            <Autocomplete
              autoComplete
              freeSolo
              selectOnFocus
              size="small"
              inputValue={inputValue}
              renderInput={Input}
              renderOption={renderResult}
              onInputChange={onInputChange}
              open={showOptions}
              onOpen={onOptionsOpen}
              onClose={onOptionsClose}
              options={options}
              filterOptions={filterResults}
              getOptionLabel={x => x.label}
              PopperComponent={renderPopper}
              closeIcon={<CustomCloseIcon />}
            />
          </ClickAwayListener>
        </Container>
      ) : (
        <SearchIconButton onClick={handleOpen}>
          <SearchIcon />
        </SearchIconButton>
      )}
    </Grid>
  );
};

export default SearchField;
