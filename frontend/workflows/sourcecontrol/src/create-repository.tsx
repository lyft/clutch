import React from "react";
import { useForm } from "react-hook-form";
import { useDataLayoutManager } from "@clutch-sh/data-layout";
import {
  Button,
  CheckboxPanel,
  client,
  Confirmation,
  MetadataTable,
  TextField,
  useWizardContext,
} from "@clutch-sh/core";
import { useDataLayout } from "@clutch-sh/data-layout";
import { Wizard, WizardStep } from "@clutch-sh/wizard";
import type { WizardChild } from "@clutch-sh/wizard";
import { MenuItem, Select } from "@material-ui/core";
import * as yup from "yup";

import type { RepostioryChild, WorkflowProps } from ".";

const REPOSITORY_OPTIONS = {
  // Fill in later, use config to override
};
const schema = yup.object().shape({
  name: yup.string().required("Name is required"),
  description: yup.string().required("Description is required"),
  owner: yup.string().required("Organization is required"),
});

const RepositoryDetails: React.FC<WizardChild> = () => {
  const { register, errors, handleSubmit } = useForm({
    mode: "onChange",
    reValidateMode: "onChange",
    validationSchema: schema,
  });
  const { onSubmit } = useWizardContext();
  const resourceData = useDataLayout("resourceData");

  const [visibility, setVisibility] = React.useState("");
  const handleVisibilityChange = e => {
    setVisibility(e.target.value);
    resourceData.updateData("github_options.parameters.visibility", e.target.value);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <TextField
        label="Name"
        name="name"
        onChange={e => resourceData.updateData("name", e.target.value)}
        inputRef={register}
        helperText={errors.name ? errors.name.message : ""}
        error={!!errors.name}
      />
      <TextField
        label="Description"
        name="description"
        onChange={e => resourceData.updateData("description", e.target.value)}
        inputRef={register}
        error={!!errors.description}
        helperText={errors.description ? errors.description.message : ""}
      />
      <TextField
        label="Organization"
        name="owner"
        onChange={e => resourceData.updateData("owner", e.target.value)}
        inputRef={register}
        error={!!errors.owner}
        helperText={errors.owner ? errors.owner.message : ""}
      />
      <Select
        labelId="Visibility"
        id="visibility"
        value={visibility}
        onChange={handleVisibilityChange}
        defaultValue="PRIVATE"
        displayEmpty
      >
        <MenuItem value="PRIVATE">Private</MenuItem>
        <MenuItem value="PUBLIC">Public</MenuItem>
      </Select>
      <div />
      <Button text="Continue" type="submit" />
    </form>
  );
};

const RepositorySettings: React.FC<RepostioryChild> = ({ options = REPOSITORY_OPTIONS }) => {
  const { onSubmit } = useWizardContext();
  const extraOptionsData = useDataLayout("extraOptionsData");

  return (
    <WizardStep error={extraOptionsData.error} isLoading={extraOptionsData.isLoading}>
      <TextField
        placeholder="Teams"
        onChange={e => extraOptionsData.updateData("teams", e.target.value)}
        onReturn={onSubmit}
      />
      <CheckboxPanel
        header="Options"
        options={options}
        onChange={(state: { [option: string]: boolean }) =>
          extraOptionsData.updateData("applicationsOptions", state)
        }
      />
      <Button text="Create Repository" onClick={onSubmit} />
    </WizardStep>
  );
};

const Confirm: React.FC<WizardChild> = () => {
  const repoData = useDataLayout("repoData");
  const instance = repoData.displayValue();

  return (
    <WizardStep error={repoData.error} isLoading={repoData.isLoading}>
      <Confirmation action="Pull Request" />
      <MetadataTable data={[{ name: "Url", value: instance.data?.url }]} />
    </WizardStep>
  );
};

const CreateRepository: React.FC<WorkflowProps> = ({ heading, options }) => {
  const dataLayout = {
    extraOptionsData: {},
    resourceData: {cache: false},
    repoData: {
      deps: ["resourceData"],
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
      <RepositorySettings name="Repository Settings" options={options} />
      <Confirm name="Confirmation" />
    </Wizard>
  );
};

export default CreateRepository;
