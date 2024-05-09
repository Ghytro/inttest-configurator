import React, { Component } from "react";
import { MockservicesApi } from "../../api/api";
import { Params, useParams } from "react-router-dom";
import { $notify, ENotifyKind } from "../../notifier";
import {
  ServiceCascaderOptVal,
  ServiceSelectionCascader,
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

class ProjectPage extends Component<ProjectPageProps, ProjectPageState> {
  private mockServiceApi: MockservicesApi;

  constructor(props: ProjectPageProps) {
    super(props);

    this.mockServiceApi = new MockservicesApi();
    this.state = {
      loading: false,
      editedServiceData: {
        id: "",
        port: 0,
      },
    };
  }

  componentDidMount(): void {
    this.setState({ loading: true });
    this.mockServiceApi
      .listServices(parseInt(this.props.urlParams.id!))
      .then(({ data }) => {
        for (const [key, value] of Object.entries(data)) {
        }
      })
      .catch((e) => {
        $notify(ENotifyKind.ERROR, e);
      })
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  render(): React.ReactNode {
    return (
      <>
        <Flex gap="middle" vertical={false} justify="space-around">
          <ServiceSelectionCascader
            mockServiceApi={this.mockServiceApi}
            projectId={parseInt(this.props.urlParams.id!)}
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
        </Flex>
        <Divider />
        {(() => {
          switch (this.state.editedServiceType) {
            case restServiceType:
              return (
                <RestServiceEditTab
                  mockServiceApi={this.mockServiceApi}
                  projectId={parseInt(this.props.urlParams.id!)}
                  serviceData={this.state.editedServiceData}
                />
              );
          }
        })()}
        <CreateRestServiceDialog
          mockServiceApi={this.mockServiceApi}
          projectId={parseInt(this.props.urlParams.id!)}
          setClosed={() => this.setState({ createdServiceType: undefined })}
          open={this.state.createdServiceType == restServiceType}
          serviceIdInputInitValue=""
          servicePortInputInitValue=""
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
}
