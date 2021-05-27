import * as React from "react";

import _ from "lodash";

import { Checkbox } from "@clutch-sh/core";


const ProjectSelector = () => {
  // On load, we'll request a list of owned projects and their upstreams and downstreams from the API.
  // The API will contain information about the relationships between projects and upstreams and downstreams.
  // By default, the owned projects will be checked and others will be unchecked.
  // If a project is unchecked, the upstream and downstreams related to it disappear from the list.
  // If a project is rechecked, the checks were preserved.

  const [projects, setProjects] = React.useState({});

  React.useEffect(() => {
    console.log("effect fired");

    setProjects({"users": true, "books": true});

    }, []);

  const upstreams = ["authors", "thumbnails"];
  const downstreams = ["coffee", "shelves"]

  const changeHandler = ({target}) => {
    setProjects({...projects, [target.name]: target.checked});
    console.log("the selected projects were changed", projects);
  }

  return (
    <div>
      My projects
      <div>
      {Object.keys(projects).map(key => <div key={key}>
        <Checkbox name={key} onChange={changeHandler}/> {key}
        </div>)}

      </div>
      Upsteam
      <div>{upstreams.map((v) => <div key={v}> <Checkbox /> {v}</div>)}
      </div>
      Downstream
      <div>{downstreams.map((v) => <div key={v}> <Checkbox /> {v}</div>)}</div>
      </div>
  )
}

const HelloWorld = () => <>Hello from workflow! <ProjectSelector /> </>

export default HelloWorld;
