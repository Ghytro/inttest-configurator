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
      brokerIdInputValue: this.props.brokerIdInputInitValue ?? "",
      brokerPortInputValue: this.props.brokerPortInputInitValue ?? "",
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

  sendFormData() {
    // todo: поправить баг, запомнить старый id
    const reqForm = {
      id: this.state.brokerIdInputValue,
      port: parseInt(this.state.brokerPortInputValue),
    };
    (this.props.modalType == "create"
      ? this.props.mockServiceApi.createRedisPubSubBroker(
          this.props.projectId,
          reqForm
        )
      : this.props.mockServiceApi.updateRedisPubSub(
          this.props.projectId,
          this.state.brokerIdInputValue,
          reqForm
        )
    )
      .then(({ data }) => {
        this.props.refetch();
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <Modal
        title={
          this.props.modalType == "create"
            ? "Создать"
            : "Редактировать" + " брокер (Redis PubSub)"
        }
        open={this.props.open}
        onCancel={() => {
          this.setState({
            brokerIdInputValue: this.props.brokerIdInputInitValue ?? "",
            brokerPortInputValue: this.props.brokerPortInputInitValue ?? "",
          });
          this.props.setClosed();
        }}
        onOk={() => this.sendFormData()}
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
  modalType: "create" | "update";
  open: boolean;
  setClosed: () => void;
  refetch: () => void;
  mockServiceApi: MockservicesApi;
  projectId: number;

  brokerIdInputInitValue?: string;
  brokerPortInputInitValue?: string;
}

interface CreateRedisPubSubDialogState {
  brokerIdInputValue: string;
  brokerPortInputValue: string;
}
