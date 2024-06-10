import { Form, Input, Modal, message } from "antd";
import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

class CreateRestServiceDialog extends Component<
  CreateRestServiceDialogProps,
  CreateRestServiceDialogState
> {
  constructor(props: CreateRestServiceDialogProps) {
    super(props);
    this.state = {
      serviceIdInputValue: this.props.serviceIdInputInitValue,
      servicePortInputValue: this.props.servicePortInputInitValue,
    };
  }

  createService() {
    this.props.mockServiceApi
      .createRestService(this.props.projectId, {
        id: this.state.serviceIdInputValue,
        port: parseInt(this.state.servicePortInputValue),
      })
      .then(({ data }) => {
        this.setState({
          serviceIdInputValue: this.props.serviceIdInputInitValue,
          servicePortInputValue: this.props.servicePortInputInitValue,
        });
        this.props.refetch();
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
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

  render(): React.ReactNode {
    return (
      <Modal
        title="Создать REST-сервис"
        open={this.props.open}
        onOk={() => {
          this.createService();
        }}
        onCancel={() => {
          this.props.setClosed();
        }}
      >
        <Form layout="vertical">
          <Form.Item label="Идентификатор">
            <Input
              value={this.state.serviceIdInputValue}
              onChange={(e) => {
                this.setState({ serviceIdInputValue: e.target.value });
              }}
            />
          </Form.Item>
          <Form.Item label="Порт">
            <Input
              value={this.state.servicePortInputValue}
              onChange={(e) => {
                try {
                  this.validatePortInputValue(e.target.value);
                } catch (e) {
                  message.error(e.message);
                  return;
                }
                this.setState({ servicePortInputValue: e.target.value });
              }}
            />
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default CreateRestServiceDialog;

interface CreateRestServiceDialogProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  open: boolean;
  setClosed: () => void;
  refetch: () => void;

  serviceIdInputInitValue: string;
  servicePortInputInitValue: string;
}

interface CreateRestServiceDialogState {
  serviceIdInputValue: string;
  servicePortInputValue: string;
}
