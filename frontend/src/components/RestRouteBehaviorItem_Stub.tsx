import React, { Component } from "react";
import { MockserviceListRestBehaviorResultStub } from "../api/api";
import { Descriptions } from "antd";

class RestRouteBehaviorItem_Stub extends Component<
  RestRouteBehaviorItem_StubProps,
  RestRouteBehaviorItem_StubState
> {
  constructor(props: RestRouteBehaviorItem_StubProps) {
    super(props);
  }

  render(): React.ReactNode {
    return (
      <Descriptions title="Описание">
        <Descriptions.Item label="URL-параметры">
          {Object.entries(this.props.behavior.url_params!).map(
            ([key, value]) => {
              return (
                <>
                  <span>{`${key}: ${value}`}</span>
                  <br />
                </>
              );
            }
          )}
        </Descriptions.Item>
      </Descriptions>
    );
  }
}

export default RestRouteBehaviorItem_Stub;

interface RestRouteBehaviorItem_StubProps {
  behavior: MockserviceListRestBehaviorResultStub;
}

interface RestRouteBehaviorItem_StubState {}
