import React, { Component } from "react";
import {
  MockserviceListBrokerTopicsResultItem,
  MockserviceListTopicGeneratorsMockItem,
  MockserviceListTopicGeneratorsStubItem,
  MockservicesApi,
} from "../api/api";
import { mockBehaviorType, stubBehaviorType } from "../const";
import CodeInputModal, { makeInitialCodeImpl_Broker } from "./CodeInputModal";
import { $notify, ENotifyKind } from "../notifier";
import { Button, Collapse, Space } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import RedisPubSubTopicBehaviorItem_Mock from "./RedisPubSubTopicBehaviorItem_Mock";
import RedisPubSubTopicBehaviorItem_Stub from "./RedisPubSubTopicBehaviorItem_Stub";
import EditRedisPubSubStubBehaviorModal from "./EditRedisPubSubStubBehaviorModal";

export default class RedisPubSubTopicBehaviors extends Component<
  RedisPubSubTopicBehaviorsProps,
  RedisPubSubTopicBehaviorsState
> {
  constructor(props: RedisPubSubTopicBehaviorsProps) {
    super(props);
    this.state = {
      addConstModalOpen: false,
      addProgModalOpen: false,
      generatorsModel: [],
      initialMockCode: "",
    };
  }

  componentDidMount(): void {
    this.fetchBehaviors();
  }

  async fetchBehaviors() {
    try {
      const resp =
        await this.props.mockServiceApi.listRedisPubSubTopicGenerators(
          this.props.projectId,
          this.props.brokerId,
          this.props.topicId
        );
      const allGenerators = resp.data
        .mocks!.map((m): GeneratorCommon => {
          return {
            type: mockBehaviorType,
            gen: m,
          };
        })
        .concat(
          resp.data.stubs!.map((s): GeneratorCommon => {
            return {
              type: stubBehaviorType,
              gen: s,
            };
          })
        )
        .sort((a, b) => {
          return a.gen.priority - b.gen.priority;
        });
      this.setState({
        generatorsModel: allGenerators,
        initialMockCode: makeInitialCodeImpl_Broker(
          this.props.brokerId,
          this.props.topicData.topic!,
          allGenerators.length
        ),
      });
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  async createMockBehavior(code: string) {
    try {
      await this.props.mockServiceApi.createRedisPubSubProgTopicGenerator(
        this.props.projectId,
        this.props.brokerId,
        this.props.topicId,
        {
          impl: code.split("\n"),
          // todo
          send_immediately: true,
          interval: "0ms",
        }
      );
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  render(): React.ReactNode {
    return (
      <>
        <Space>
          <Button
            icon={<PlusOutlined />}
            onClick={() => {
              this.setState({ addConstModalOpen: true });
            }}
          >
            Добавить константный
          </Button>
          <Button
            icon={<PlusOutlined />}
            onClick={() => {
              this.setState({ addProgModalOpen: true });
            }}
          >
            Добавить программируемый
          </Button>
        </Space>
        <Collapse
          items={this.state.generatorsModel.map((generator) => {
            switch (generator.type) {
              case mockBehaviorType: {
                const gen =
                  generator.gen as MockserviceListTopicGeneratorsMockItem;
                return {
                  key: gen.id,
                  label: `Генератор №${gen.priority! + 1} (программируемый)`,
                  children: (
                    <RedisPubSubTopicBehaviorItem_Mock
                      brokerId={this.props.brokerId}
                      generatorId={gen.id!}
                      impl={(gen.impl ?? []).join("\n")}
                      mockServiceApi={this.props.mockServiceApi}
                      projectId={this.props.projectId}
                      refetch={() => {
                        this.fetchBehaviors();
                      }}
                      topicId={this.props.topicId}
                    />
                  ),
                };
              }
              case stubBehaviorType: {
                const gen =
                  generator.gen as MockserviceListTopicGeneratorsStubItem;
                return {
                  key: gen.id,
                  label: `Генератор №${gen.priority! + 1} (константный)`,
                  children: (
                    <RedisPubSubTopicBehaviorItem_Stub
                      refetch={() => this.fetchBehaviors()}
                      brokerId={this.props.brokerId}
                      projectId={this.props.projectId}
                      generatorId={gen.id!}
                      intervalModulus={parseInt(
                        gen.interval!.slice(
                          0,
                          gen.interval!.search("/[^A-Za-z]/")
                        )
                      )}
                      intervalUnits={gen.interval!.slice(
                        gen.interval!.search("/[^A-Za-z]/") - 1
                      )}
                      mockServiceApi={this.props.mockServiceApi}
                      payload={gen.payload!}
                      sendImmediately={gen.send_immediately!}
                      topicId={this.props.topicId}
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
        />
        <EditRedisPubSubStubBehaviorModal
          mockServiceApi={this.props.mockServiceApi}
          brokerId={this.props.brokerId}
          projectId={this.props.projectId}
          topicId={this.props.topicId}
          open={this.state.addConstModalOpen}
          setClosed={() => {
            this.setState({ addConstModalOpen: false });
          }}
          modalType="create"
        />
        <CodeInputModal
          open={this.state.addProgModalOpen}
          setClosed={() => {
            this.setState({ addProgModalOpen: false });
          }}
          submit={(code) => {
            this.createMockBehavior(code);
          }}
        />
      </>
    );
  }
}

interface RedisPubSubTopicBehaviorsProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  brokerId: string;
  topicId: number;
  topicData: MockserviceListBrokerTopicsResultItem;
}

interface RedisPubSubTopicBehaviorsState {
  generatorsModel: GeneratorCommon[];
  initialMockCode: string;
  addConstModalOpen: boolean;
  addProgModalOpen: boolean;
}

interface GeneratorCommon {
  type: string;
  gen: any;
}
