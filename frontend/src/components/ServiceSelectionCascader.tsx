import { Cascader, CascaderProps } from "antd";
import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { translateServiceType } from "../const";

class ServiceSelectionCascader extends Component<
  ServiceSelectionCascaderProps,
  ServiceSelectionCascaderState
> {
  constructor(props: ServiceSelectionCascaderProps) {
    super(props);

    this.state = {
      loading: false,
      cascaderState: undefined,
    };
  }

  componentDidMount(): void {
    this.fetchServicesList();
  }

  fetchServicesList() {
    this.setState({ loading: true });
    this.props.mockServiceApi
      .listServices(this.props.projectId)
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

        this.setState({ cascaderState: serviceTypes });
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
      <Cascader
        options={this.state.cascaderState}
        onChange={(values: any[]) => {
          this.props.selectEditedService(values[0], values[1]);
        }}
      />
    );
  }
}

///////// STATE

interface ServiceSelectionCascaderState {
  loading: boolean;
  cascaderState?: ServiceTypeCascaderOpt[];
}

interface ServiceTypeCascaderOpt {
  value: string;
  label: string;
  children?: ServiceCascaderOpt[];
}

interface ServiceCascaderOpt {
  value: ServiceCascaderOptVal;
  label: string;
}

interface ServiceCascaderOptVal {
  id: string;
  port: number;
}

///////// PROPS

interface ServiceSelectionCascaderProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  selectEditedService: (
    serviceType: string,
    serviceData: ServiceCascaderOptVal
  ) => void;
}

export { ServiceSelectionCascader, ServiceCascaderOptVal };
