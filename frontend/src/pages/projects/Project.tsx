import React, { Component } from "react";
import { MockservicesApi, ProjectsApi } from "../../api/api";
import { Params, useParams } from "react-router-dom";
import { $notify, ENotifyKind } from "../../notifier";
import {
  ServiceCascaderOpt,
  ServiceCascaderOptVal,
  ServiceSelectionCascader,
  ServiceTypeCascaderOpt,
} from "../../components/ServiceSelectionCascader";
import { Button, Divider, Dropdown, Flex, Space } from "antd";
import {
  redisServiceType,
  restServiceType,
  soapServiceType,
  translateServiceType,
} from "../../const";
import { ItemType } from "antd/es/menu/hooks/useItems";
import CreateRestServiceDialog from "../../components/CreateRestServiceDialog";
import RestServiceEditTab from "../../components/RestServiceEditTab";
import { DownOutlined } from "@ant-design/icons";
import CreateRedisPubSubDialog from "../../components/CreateRedisPubSubDialog";
import RedisPubSubEditTab from "../../components/RedisPubSubEditTab";

class ProjectPage extends Component<ProjectPageProps, ProjectPageState> {
  private mockServiceApi: MockservicesApi;
  private projectsApi: ProjectsApi;

  constructor(props: ProjectPageProps) {
    super(props);

    this.mockServiceApi = new MockservicesApi();
    this.projectsApi = new ProjectsApi();
    this.state = {
      loading: false,
      editedServiceData: {
        id: "",
        port: 0,
      },
    };
  }

  fetchServicesList() {
    this.setState({ loading: true });
    this.mockServiceApi
      .listServices(parseInt(this.props.urlParams.id!))
      .then(({ data }) => {
        const serviceTypes: ServiceTypeCascaderOpt[] = Object.entries(data).map(
          ([key, value]) => ({
            value: key,
            label: translateServiceType[key],
            children: value.map((v) => ({
              value: {
                id: v.id!,
                port: v.port!,
              },
              label: v.id!,
            })),
          })
        );

        this.setState({ serviceCascaderState: serviceTypes });
      })
      .catch((e) => {
        $notify(ENotifyKind.ERROR, e);
      })
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  async unloadProject() {
    try {
      const resp = await this.projectsApi.unloadProject(
        parseInt(this.props.urlParams.id!)
      );
      const cd = resp.headers["content-disposition"];
      let fileName = cd.split("filename=")[1].replace(/["']/g, "");
      fileName = decodeURIComponent(fileName).replace(/\+/g, " ");
      var file = new Blob([JSON.stringify(resp.data, null, 4)], {
        type: "text/json",
      });
      let url = URL.createObjectURL(file);
      this.openFile(url, fileName);
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  openFile(url: string, fileName = "configuration.json"): void {
    const a = document.createElement("a");
    a.href = url;
    a.download = fileName;
    document.body.append(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
  }

  componentDidMount(): void {
    this.fetchServicesList();
  }

  render(): React.ReactNode {
    return (
      <>
        <Flex gap="middle" vertical={false} justify="space-around">
          <ServiceSelectionCascader
            cascaderState={this.state.serviceCascaderState}
            selectEditedService={(
              serviceType: string | undefined,
              serviceData: ServiceCascaderOptVal
            ) => {
              this.setState({
                editedServiceType: serviceType,
                editedServiceData: serviceData,
              });
            }}
          />
          <Dropdown
            menu={{
              items: [restServiceType, soapServiceType, redisServiceType].map(
                (t): ItemType => {
                  return {
                    key: t,
                    label: (
                      <Button
                        type="link"
                        onClick={() => {
                          this.setState({ createdServiceType: t });
                        }}
                      >
                        {translateServiceType[t]}
                      </Button>
                    ),
                  };
                }
              ),
            }}
          >
            <a onClick={(e) => e.preventDefault()}>
              <Space>
                Создать сервис...
                <DownOutlined />
              </Space>
            </a>
          </Dropdown>
          <Button
            onClick={() => {
              this.unloadProject();
            }}
          >
            Выгрузить проект
          </Button>
        </Flex>
        <Divider />

        {(() => {
          switch (this.state.editedServiceType) {
            case restServiceType:
              return (
                <RestServiceEditTab
                  key={this.state.editedServiceData.id}
                  mockServiceApi={this.mockServiceApi}
                  projectId={parseInt(this.props.urlParams.id!)}
                  serviceData={this.state.editedServiceData}
                />
              );
            case redisServiceType:
              return (
                <RedisPubSubEditTab
                  key={this.state.editedServiceData.id}
                  mockServiceApi={this.mockServiceApi}
                  brokerData={this.state.editedServiceData}
                  projectId={parseInt(this.props.urlParams.id!)}
                />
              );
          }
        })()}

        <CreateRestServiceDialog
          mockServiceApi={this.mockServiceApi}
          projectId={parseInt(this.props.urlParams.id!)}
          setClosed={() => this.setState({ createdServiceType: undefined })}
          open={this.state.createdServiceType == restServiceType}
          refetch={() => this.fetchServicesList()}
          serviceIdInputInitValue=""
          servicePortInputInitValue=""
        />
        <CreateRedisPubSubDialog
          modalType="create"
          open={this.state.createdServiceType == redisServiceType}
          setClosed={() => this.setState({ createdServiceType: undefined })}
          mockServiceApi={this.mockServiceApi}
          projectId={parseInt(this.props.urlParams.id!)}
          refetch={() => this.fetchServicesList()}
        />
      </>
    );
  }
}

const ProjectPageFC = (props) => {
  const params = useParams();
  return <ProjectPage {...props} urlParams={params} />;
};

export default ProjectPageFC;

interface ProjectPageProps {
  urlParams: Readonly<Params<string>>;
}

interface ProjectPageState {
  loading: boolean;
  createdServiceType?: string;
  editedServiceData: ServiceCascaderOptVal;
  editedServiceType?: string;
  serviceCascaderState?: ServiceTypeCascaderOpt[];
}
