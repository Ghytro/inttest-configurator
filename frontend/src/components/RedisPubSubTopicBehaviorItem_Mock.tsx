import React, { Component } from "react";
import CodeInputModal from "./CodeInputModal";
import { Button, Descriptions } from "antd";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

export default class RedisPubSubTopicBehaviorItem_Mock extends Component<
  RedisPubSubTopicBehaviorItem_MockProps,
  RedisPubSubTopicBehaviorItem_MockState
> {
  constructor(props: RedisPubSubTopicBehaviorItem_MockProps) {
    super(props);

    this.state = {
      editImplDialogOpen: false,
    };
  }

  async submitData(code: string) {
    try {
      await this.props.mockServiceApi.updateRedisPubSubProgTopicGenerator(
        this.props.projectId,
        this.props.brokerId,
        this.props.topicId,
        this.props.generatorId,
        {
          impl: code.split("\n"),
          // todo
          interval: "0ms",
          send_immediately: true,
        }
      );
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    } finally {
      this.props.refetch();
      this.setState({ editImplDialogOpen: false });
    }
  }

  render(): React.ReactNode {
    return (
      <>
        <Descriptions>
          <Descriptions.Item label="Имплементация">
            <Button
              type="link"
              onClick={() => {
                this.setState({ editImplDialogOpen: true });
              }}
            >
              Редактировать
            </Button>
          </Descriptions.Item>
        </Descriptions>
        <CodeInputModal
          open={this.state.editImplDialogOpen}
          setClosed={() => this.setState({ editImplDialogOpen: false })}
          submit={(code) => {
            this.submitData(code);
          }}
          code={this.props.impl}
        />
      </>
    );
  }
}

interface RedisPubSubTopicBehaviorItem_MockProps {
  impl: string;
  mockServiceApi: MockservicesApi;
  projectId: number;
  brokerId: string;
  topicId: number;
  generatorId: number;
  refetch: () => void;
}

interface RedisPubSubTopicBehaviorItem_MockState {
  editImplDialogOpen: boolean;
}
