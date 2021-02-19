import React from "react";
import { useForm } from "react-hook-form";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  AvatarIcon,
  Button,
  ButtonGroup,
  client,
  Confirmation,
  Form,
  Link,
  Select,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import type { WizardChild } from "@clutch-sh/wizard";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import styled from "@emotion/styled";
import LockIcon from "@material-ui/icons/Lock";
import LockOpenIcon from "@material-ui/icons/LockOpen";
import * as yup from "yup";

import type { WorkflowProps } from ".";

const ConfirmationMessage = styled.div({
  display: "flex",
  alignItems: "center",
  fontSize: "14px",
});

const schema = yup.object().shape({
  name: yup.string().required("Name is required"),
  description: yup.string().required("Description is required"),
});

const visibilityOptions = {
  PRIVATE: { label: "Private", value: "PRIVATE", startAdornment: <LockIcon /> },
  PUBLIC: { label: "Public", value: "PUBLIC", startAdornment: <LockOpenIcon /> },
};

const RepositoryDetails: React.FC<WizardChild> = () => {
  const { register, errors, handleSubmit } = useForm({
    mode: "onChange",
    reValidateMode: "onChange",
    validationSchema: schema,
  });
  const { onSubmit } = useWizardContext();
  const repositoryData = useDataLayout("repositoryData");
  const repositoryOptions = repositoryData.displayValue()
    .data as IClutch.sourcecontrol.v1.GetRepositoryOptionsResponse;

  return (
    <WizardStep error={repositoryData.error} isLoading={repositoryData.isLoading}>
      <Form onSubmit={handleSubmit(onSubmit)}>
        <Select
          label="Owner"
          name="owner"
          onChange={value => repositoryData.updateData("owner", value)}
          options={repositoryOptions?.availableOwners?.map(owner => {
            return { label: owner.name, startAdornment: <AvatarIcon src={owner.photoUrl} /> };
          })}
        />
        <TextField
          label="Name"
          name="name"
          onChange={e => repositoryData.updateData("name", e.target.value)}
          inputRef={register}
          helperText={errors.name ? errors.name.message : ""}
          error={!!errors.name}
        />
        <TextField
          label="Description"
          name="description"
          onChange={e => repositoryData.updateData("description", e.target.value)}
          inputRef={register}
          error={!!errors.description}
          helperText={errors.description ? errors.description.message : ""}
        />
        <Select
          name="visibility"
          label="Visibility"
          onChange={value =>
            repositoryData.updateData("github_options.parameters.visibility", value)
          }
          defaultOption={0}
          options={repositoryOptions?.visibilityOptions?.map(
            visibility => visibilityOptions?.[visibility]
          )}
        />
        <ButtonGroup>
          <Button text="Continue" type="submit" />
        </ButtonGroup>
      </Form>
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const repoData = useDataLayout("repoData");
  const instance = repoData.displayValue();

  return (
    <WizardStep error={repoData.error} isLoading={repoData.isLoading}>
      <Confirmation action="Repository Creation">
        <ConfirmationMessage>
          Your new repository can be found&nbsp; <Link href={instance.data?.url}>here</Link>
        </ConfirmationMessage>
      </Confirmation>
    </WizardStep>
  );
};

const CreateRepository: React.FC<WorkflowProps> = ({ heading, options }) => {
  const dataLayout = {
    repositoryData: {
      cache: false,
      hydrator: () => {
        return client.post("/v1/sourcecontrol/getRepositoryOptions");
      },
    },
    repoData: {
      deps: ["repositoryData"],
      cache: false,
      hydrator: resourceData => {
        return client.post("/v1/sourcecontrol/createRepository", {
          name: resourceData.name,
          description: resourceData.description,
          owner: resourceData.owner,
          github_options: resourceData.github_options,
        });
      },
    },
  };

  return (
    <Wizard dataLayout={dataLayout} heading={heading}>
      <RepositoryDetails name="Repository Details" />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default CreateRepository;
