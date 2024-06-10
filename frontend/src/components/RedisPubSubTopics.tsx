import { Button, Collapse, Space } from "antd";
import React, { Component } from "react";
import {
  MockserviceListBrokerTopicsResultItem,
  MockservicesApi,
} from "../api/api";
import { $notify, ENotifyKind } from "../notifier";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import RedisPubSubTopicBehaviors from "./RedisPubSubTopicBehaviors";

export default class RedisPubSubTopics extends Component<
  RedisPubSubTopicsProps,
  RedisPubSubTopicsState
> {
  constructor(props: RedisPubSubTopicsProps) {
    super(props);
    this.state = {
      allTopics: [],
    };
  }

  componentDidMount(): void {
    this.fetchAllTopics();
  }

  async fetchAllTopics() {
    try {
      const resp = await this.props.mockServiceApi.listRedisPubSubTopics(
        this.props.projectId,
        this.props.brokerId
      );
      this.setState({ allTopics: resp.data });
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  async deleteTopic(id: number) {
    try {
      await this.props.mockServiceApi.deleteRedisPubTopic(
        this.props.projectId,
        this.props.brokerId,
        id
      );
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  render(): React.ReactNode {
    return (
      <>
        <Collapse
          accordion
          items={this.state.allTopics.map((t) => {
            return {
              key: t.id,
              label: t.topic,
              children: (
                <>
                  <Space>
                    <Button
                      type="default"
                      icon={<EditOutlined />}
                      onClick={() => {}}
                    >
                      Редактировать
                    </Button>
                    <Button
                      danger
                      icon={<DeleteOutlined />}
                      onClick={() => {
                        this.deleteTopic(t.id!);
                      }}
                    >
                      Удалить
                    </Button>
                  </Space>
                  <RedisPubSubTopicBehaviors
                    mockServiceApi={this.props.mockServiceApi}
                    projectId={this.props.projectId}
                    brokerId={this.props.brokerId}
                    topicId={t.id!}
                    topicData={t}
                  />
                </>
              ),
            };
          })}
        />
      </>
    );
  }
}

interface RedisPubSubTopicsProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  brokerId: string;
}

interface RedisPubSubTopicsState {
  allTopics: MockserviceListBrokerTopicsResultItem[];
}
