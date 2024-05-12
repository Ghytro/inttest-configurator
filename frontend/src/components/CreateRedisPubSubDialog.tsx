import { Form, Input, Modal, message } from "antd";
import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

class CreateRedisPubSubDialog extends Component<
  CreateRedisPubSubDialogProps,
  CreateRedisPubSubDialogState
> {
  constructor(props: CreateRedisPubSubDialogProps) {
    super(props);

    this.state = {
      brokerIdInputValue: "",
      brokerPortInputValue: "",
    };
  }

  validatePortInputValue(portStr: string) {
    const numbers = "0123456789";
    for (let i = 0; i < portStr.length; ++i) {
      if (!numbers.includes(portStr[i])) {
        throw new Error("порт должен быть валидным числом");
      }
    }
    const portInt = parseInt(portStr);
    if (portInt > 65535) {
      throw new Error("порт должен быть числом от 0 до 65535");
    }
  }

  createBroker() {
    this.props.mockServiceApi
      .createRedisPubSubBroker(this.props.projectId, {
        id: this.state.brokerIdInputValue,
        port: parseInt(this.state.brokerPortInputValue),
      })
      .then(({ data }) => {
        this.props.refetch();
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <Modal
        title="Создать брокер (Redis PubSub)"
        open={this.props.open}
        onCancel={() => {
          this.setState({});
          this.props.setClosed();
        }}
      >
        <Form layout="vertical">
          <Form.Item label="Идентификатор">
            <Input
              value={this.state.brokerIdInputValue}
              onChange={(e) => {
                this.setState({ brokerIdInputValue: e.target.value });
              }}
            />
          </Form.Item>
          <Form.Item label="Порт">
            <Input
              value={this.state.brokerPortInputValue}
              onChange={(e) => {
                try {
                  this.validatePortInputValue(e.target.value);
                } catch (e) {
                  message.error(e.message);
                  return;
                }
                this.setState({ brokerPortInputValue: e.target.value });
              }}
            />
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default CreateRedisPubSubDialog;

interface CreateRedisPubSubDialogProps {
  open: boolean;
  setClosed: () => void;
  refetch: () => void;
  mockServiceApi: MockservicesApi;
  projectId: number;
}

interface CreateRedisPubSubDialogState {
  brokerIdInputValue: string;
  brokerPortInputValue: string;
}
