import React, { Component } from "react";
import { ProjectsApi } from "../../api/api";
import { $notify, ENotifyKind } from "../../notifier";
import { FloatButton, List } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import ProjCreateModal from "../../components/ProjCreateModal";

class ProjectsPage extends Component<any, ProjectsPageState> {
  private projectApi: ProjectsApi;

  constructor(props: any) {
    super(props);

    this.projectApi = new ProjectsApi();
    this.state = {
      loading: true,
      projectsList: [],
      projCreateModalOpen: false,
    };
    this.fetchProjectList();
  }

  fetchProjectList() {
    this.setState({ loading: true });
    this.projectApi
      .listProjects()
      .then(({ data }) => {
        this.setState({
          projectsList: data.map(
            (resp) =>
              ({
                name: resp.name,
                description: resp.desc,
              } as ViewedProject)
          ),
        });
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  render(): React.ReactNode {
    const projectsList =
      this.state.projectsList.length > 0 ? (
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
      ) : (
        <div>Нет проектов</div>
      );

    return (
      <>
        {projectsList}
        <FloatButton
          icon={<PlusOutlined />}
          type="primary"
          onClick={() => {
            this.setState({ projCreateModalOpen: true });
          }}
        />
        <ProjCreateModal
          projectsApi={this.projectApi}
          open={this.state.projCreateModalOpen}
          setClosed={() => {
            this.setState({ projCreateModalOpen: false });
            this.fetchProjectList();
          }}
        />
      </>
    );
  }
}

interface ProjectsPageState {
  loading: boolean;
  projectsList: ViewedProject[];
  projCreateModalOpen: boolean;
}

interface ViewedProject {
  name: string;
  description: string;
}

export default ProjectsPage;
