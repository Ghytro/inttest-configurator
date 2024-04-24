import React, { Component } from "react";
import { ProjectsApi } from "../../api/api";
import { $notify, ENotifyKind } from "../../notifier";
import { List } from "antd";

class ProjectsPage extends Component<any, ProjectsPageState> {
  private projectApi: ProjectsApi;
  constructor(props: any) {
    super(props);

    this.projectApi = new ProjectsApi();
    this.setState({loading: true});
    this.projectApi.listProjects()
      .then(({data}) => {
        this.setState({
          projectsList: data.map((resp) => ({
            name: resp.name,
            description: resp.desc,
          } as ViewedProject))
        });
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {
        this.setState({loading: false});
      });
  }

  render(): React.ReactNode {
    return (
      <>
        <List
          loading={this.state.loading}
          size="large"
          dataSource={this.state.projectsList}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                title={item.name}
                description={item.description}
              />
            </List.Item>
          )}
        />
      </>
    );
  }
}

interface ProjectsPageState {
  loading: boolean;
  projectsList: ViewedProject[];
}

interface ViewedProject {
  name: string;
  description: string;
}

export default ProjectsPage;
