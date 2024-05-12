import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultMock,
  MockserviceListRestBehaviorResultStub,
  MockservicesApi,
} from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { Button, Collapse, Space, Table } from "antd";
import RestRouteBehaviorItem_Mock from "./RestRouteBehaviorItem_Mock";
import { mockBehaviorType, stubBehaviorType } from "../const";
import RestRouteBehaviorItem_Stub from "./RestRouteBehaviorItem_Stub";
import { PlusOutlined } from "@ant-design/icons";

class RestRouteBehaviors extends Component<
  RestRouteBehaviorsProps,
  RestRouteBehaviorsState
> {
  constructor(props: RestRouteBehaviorsProps) {
    super(props);

    this.state = {
      loading: false,
      behaviorsModel: [],
    };
  }

  componentDidMount(): void {
    this.fetchBehaviors();
  }

  fetchBehaviors() {
    this.setState({ loading: true });
    this.props.mockServiceApi
      .listRestBehaviors(
        this.props.projectId,
        this.props.serviceId,
        this.props.handlerId
      )
      .then(({ data }) => {
        const anyBehaviors = data
          .mocks!.map((m): BehaviorCommon => {
            return {
              type: mockBehaviorType,
              behavior: m,
            };
          })
          .concat(
            data.stubs!.map((s): BehaviorCommon => {
              return {
                type: stubBehaviorType,
                behavior: s,
              };
            })
          )
          .sort((a, b) => {
            return a.behavior.priority - b.behavior.priority;
          });
        this.setState({ behaviorsModel: anyBehaviors });
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => {
        this.setState({ loading: false });
      });
  }

  render(): React.ReactNode {
    return (
      <>
        <Space>
          <Button icon={<PlusOutlined />} onClick={() => {}}>
            Добавить stub
          </Button>
          <Button icon={<PlusOutlined />} onClick={() => {}}>
            Добавить mock
          </Button>
        </Space>
        <Collapse
          items={this.state.behaviorsModel.map((b) => {
            switch (b.type) {
              case mockBehaviorType: {
                const behavior =
                  b.behavior as MockserviceListRestBehaviorResultMock;
                return {
                  key: behavior.id,
                  label: `Поведение №${behavior.priority! + 1} (${b.type})`,
                  children: (
                    <RestRouteBehaviorItem_Mock
                      mockServiceApi={this.props.mockServiceApi}
                      projectId={this.props.projectId}
                      serviceId={this.props.serviceId}
                      handlerId={this.props.handlerId}
                      refetch={() => {
                        this.fetchBehaviors();
                      }}
                      behavior={behavior}
                    />
                  ),
                };
              }
              case stubBehaviorType: {
                const behavior =
                  b.behavior as MockserviceListRestBehaviorResultStub;
                return {
                  key: behavior.id,
                  label: `Поведение №${behavior.priority! + 1} (${b.type})`,
                  children: (
                    <RestRouteBehaviorItem_Stub
                      behavior={behavior}
                      projectId={this.props.projectId}
                      handlerId={this.props.handlerId}
                      mockServiceApi={this.props.mockServiceApi}
                      serviceId={this.props.serviceId}
                      refetch={() => {
                        this.fetchBehaviors();
                      }}
                    />
                  ),
                };
              }
              default: {
                return {
                  key: -1,
                  label: "",
                  children: <></>,
                };
              }
            }
          })}
        ></Collapse>
      </>
    );
  }
}

export default RestRouteBehaviors;

interface RestRouteBehaviorsProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId: number;
}

interface RestRouteBehaviorsState {
  loading: boolean;
  behaviorsModel: BehaviorCommon[];
}

interface BehaviorCommon {
  type: string;
  behavior: any;
}
