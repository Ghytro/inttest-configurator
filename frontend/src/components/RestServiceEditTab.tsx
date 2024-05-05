import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { Button, FloatButton } from "antd";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { ServiceCascaderOptVal } from "./ServiceSelectionCascader";
import { projPageUrl } from "../routesEnum";

class RestServiceEditTab extends Component<
  RestServiceEditTabProps,
  RestServiceEditTabState
> {
  constructor(props: RestServiceEditTabProps) {
    super(props);
    this.state = {
      loading: false,
      serviceEditDialogOpen: false,
      addHandlerDialogOpen: false,
    };
  }

  componentDidMount(): void {}

  fetchRoutes() {
    this.setState({ loading: true });
    this.props.mockServiceApi
      .listRestHandlers(this.props.projectId, this.props.serviceData.id)
      .then(({ data }) => {})
      .catch((e) => {
        $notify(ENotifyKind.ERROR, e);
      })
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  deleteService() {
    this.props.mockServiceApi
      .deleteRestService(this.props.serviceData.id, this.props.projectId)
      .then((_) => {
        window.location.href = projPageUrl(this.props.projectId);
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <>
        <Button
          type="default"
          icon={<EditOutlined />}
          onClick={() => {
            this.setState({ serviceEditDialogOpen: true });
          }}
        >
          Редактировать сервис
        </Button>
        <Button
          danger={true}
          icon={<DeleteOutlined />}
          onClick={() => {
            this.deleteService();
          }}
        >
          Удалить сервис
        </Button>
        <FloatButton
          icon={<PlusOutlined />}
          onClick={() => this.setState({ addHandlerDialogOpen: true })}
        />
      </>
    );
  }
}

export default RestServiceEditTab;

interface RestServiceEditTabProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceData: ServiceCascaderOptVal;
}

interface RestServiceEditTabState {
  loading: boolean;
  serviceEditDialogOpen: boolean;
  addHandlerDialogOpen: boolean;
}
