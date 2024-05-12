import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultStub,
  MockserviceListRestRoutesResult,
  MockservicesApi,
} from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { Button, Collapse } from "antd";
import RestRouteBehaviors from "./RestRouteBehaviors";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import RestRouteUpdateModal from "./RestRouteUpdateModal";

class RestServiceRoutes extends Component<
  RestServiceRouteProps,
  RestServiceRouteState
> {
  constructor(props: RestServiceRouteProps) {
    super(props);
    this.state = {
      loading: false,
      updRouteModalOpen: false,
    };
  }

  componentDidMount(): void {
    this.fetchAllRoutes();
  }

  fetchAllRoutes() {
    this.setState({ loading: true });
    this.props.mockServiceApi
      .listRestHandlers(this.props.projectId, this.props.serviceId)
      .then(({ data }) => {
        this.setState({ allRoutes: data });
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  deleteRoute(id: number) {
    this.setState({ loading: true });
    this.props.mockServiceApi
      .deleteRestHandler(this.props.projectId, this.props.serviceId, id)
      .then(({ data }) => {
        this.fetchAllRoutes();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  render(): React.ReactNode {
    return (
      <>
        <Collapse
          accordion
          items={this.state.allRoutes?.map((r) => {
            return {
              key: r.id!,
              label: `${r.route} [${r.method}]`,
              children: (
                <>
                  <Button
                    type="default"
                    icon={<EditOutlined />}
                    onClick={() => {
                      this.setState({
                        updRouteModalOpen: true,
                        updRouteId: r.id,
                        updatedHttpMethod: r.method,
                        updatedRoute: r.route,
                      });
                    }}
                  >
                    Редактировать
                  </Button>
                  <Button
                    danger
                    icon={<DeleteOutlined />}
                    onClick={() => {
                      this.deleteRoute(r.id!);
                    }}
                  >
                    Удалить
                  </Button>
                  <RestRouteBehaviors
                    mockServiceApi={this.props.mockServiceApi}
                    projectId={this.props.projectId}
                    serviceId={this.props.serviceId}
                    handlerId={r.id!}
                  />
                </>
              ),
            };
          })}
        />
        <RestRouteUpdateModal
          modalType="update"
          urlInputInit={this.state.updatedRoute}
          httpMethodInputInit={this.state.updatedHttpMethod}
          mockServiceApi={this.props.mockServiceApi}
          open={this.state.updRouteModalOpen}
          projectId={this.props.projectId}
          serviceId={this.props.serviceId}
          setClosed={() =>
            this.setState({
              updRouteModalOpen: false,
              updRouteId: undefined,
              updatedHttpMethod: undefined,
              updatedRoute: undefined,
            })
          }
          handlerId={this.state.updRouteId}
        />
      </>
    );
  }
}

export default RestServiceRoutes;

interface RestServiceRouteProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
}

interface RestServiceRouteState {
  loading: boolean;
  allRoutes?: MockserviceListRestRoutesResult[];

  updRouteId?: number;
  updRouteModalOpen: boolean;
  updatedRoute?: string;
  updatedHttpMethod?: string;
}
