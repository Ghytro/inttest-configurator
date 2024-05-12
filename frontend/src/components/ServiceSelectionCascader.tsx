import { Cascader } from "antd";
import React, { Component } from "react";

class ServiceSelectionCascader extends Component<
  ServiceSelectionCascaderProps,
  ServiceSelectionCascaderState
> {
  constructor(props: ServiceSelectionCascaderProps) {
    super(props);

    this.state = {
      loading: false,
    };
  }

  render(): React.ReactNode {
    return (
      <Cascader
        options={this.props.cascaderState}
        onChange={(values: any[]) => {
          this.props.selectEditedService(values[0], values[1]);
        }}
        placeholder="Сервис..."
      />
    );
  }
}

///////// STATE

interface ServiceSelectionCascaderState {
  loading: boolean;
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
  cascaderState?: ServiceTypeCascaderOpt[];
  selectEditedService: (
    serviceType: string | undefined,
    serviceData: ServiceCascaderOptVal
  ) => void;
}

export {
  ServiceSelectionCascader,
  ServiceTypeCascaderOpt,
  ServiceCascaderOpt,
  ServiceCascaderOptVal,
};
