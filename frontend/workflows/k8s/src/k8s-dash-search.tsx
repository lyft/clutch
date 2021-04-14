import React from "react";
import { useForm } from "react-hook-form";
import { Form, FormRow, IconButton, Paper, TextField } from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import { yupResolver } from "@hookform/resolvers/yup";
import SearchIcon from "@material-ui/icons/Search";
import * as yup from "yup";

const Container = styled.div({
  margin: "32px 0",
});

const schema = yup.object().shape({
  namespace: yup.string().required("Namespace is required"),
  clientset: yup.string().required("Clientset is required"),
});

const K8sDashSearch = ({ onSubmit }) => {
  const { errors, handleSubmit, register } = useForm({
    mode: "onChange",
    reValidateMode: "onChange",
    resolver: yupResolver(schema),
  });
  const inputData = useDataLayout("inputData");
  const podListData = useDataLayout("podListData");

  const submitHandler = v => {
    inputData.updateData("namespace", v.namespace);
    inputData.updateData("clientset", v.clientset);
    podListData.hydrate();
    onSubmit();
  };

  return (
    <Container>
      <Paper>
        <Form onSubmit={handleSubmit(submitHandler)}>
          <FormRow>
            <TextField
              defaultValue={inputData.displayValue()?.namespace}
              placeholder="Namespace"
              name="namespace"
              error={!!errors?.namespace}
              helperText={errors?.namespace?.message}
              inputRef={register}
            />
            <TextField
              defaultValue={inputData.displayValue()?.clientset}
              placeholder="Clientset"
              name="clientset"
              error={!!errors?.clientset}
              helperText={errors?.clientset?.message}
              inputRef={register}
            />
            <IconButton type="submit">
              <SearchIcon />
            </IconButton>
          </FormRow>
        </Form>
      </Paper>
    </Container>
  );
};

export default K8sDashSearch;
