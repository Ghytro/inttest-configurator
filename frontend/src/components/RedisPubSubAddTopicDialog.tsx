import { Input, Modal } from "antd";
import React, { Component } from "react";
import { Form } from "antd";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

export default class RedisPubSubAddTopicDialog extends Component<
  RedisPubSubAddTopicDialogProps,
  RedisPubSubAddTopicDialogState
> {
  constructor(props: RedisPubSubAddTopicDialogProps) {
    super(props);

    this.state = {
      topicName: props.topicName ?? "",
    };
  }

  sendFormData() {
    const reqForm = {
      topic: this.state.topicName,
    };
    (this.props.modalType == "create"
      ? this.props.mockServiceApi.createRedisPubSubTopic(
          this.props.projectId,
          this.props.brokerId,
          reqForm
        )
      : this.props.mockServiceApi.updateRedisPubSubTopic(
          this.props.projectId,
          this.props.brokerId,
          this.props.topicId!
        )
    )
      .then(({ data }) => {
        this.setState({
          topicName: this.props.topicName ?? "",
        });
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onCancel={() => this.props.setClosed()}
        onOk={() => {
          this.sendFormData();
        }}
      >
        <Form layout="vertical">
          <Form.Item label="Топик">
            <Input
              type="text"
              defaultValue={this.props.topicName}
              onChange={(e) => this.setState({ topicName: e.target.value })}
            />
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

interface RedisPubSubAddTopicDialogProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  brokerId: string;
  modalType: "create" | "update";
  open: boolean;
  setClosed: () => void;
  topicName?: string;
  topicId?: number;
}

interface RedisPubSubAddTopicDialogState {
  topicName: string;
}
