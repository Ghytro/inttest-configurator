import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultStub,
  MockserviceListRestRoutesResult,
  MockservicesApi,
} from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { Collapse } from "antd";
import RestRouteBehaviors from "./RestRouteBehaviors";

class RestServiceRoutes extends Component<
  RestServiceRouteProps,
  RestServiceRouteState
> {
  constructor(props: RestServiceRouteProps) {
    super(props);
    this.state = {
      loading: false,
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

  render(): React.ReactNode {
    return (
      <Collapse
        accordion
        items={this.state.allRoutes?.map((r) => {
          return {
            key: r.id!,
            label: `${r.route} [${r.method}]`,
            children: (
              <RestRouteBehaviors
                mockServiceApi={this.props.mockServiceApi}
                projectId={this.props.projectId}
                serviceId={this.props.serviceId}
                handlerId={r.id!}
              />
            ),
          };
        })}
      />
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
}
