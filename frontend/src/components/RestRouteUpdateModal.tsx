import { Form, Input, Modal, Select } from "antd";
import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

class RestRouteUpdateModal extends Component<
  RestRouteUpdateModalProps,
  RestRouteUpdateModalState
> {
  constructor(props: RestRouteUpdateModalProps) {
    super(props);

    this.state = {
      httpMethodInput: this.props.httpMethodInputInit ?? "",
      urlInput: this.props.urlInputInit ?? "",
    };
  }

  submitData() {
    const form = {
      method: this.state.httpMethodInput,
      route: this.state.urlInput,
    };
    this.props.modalType == "update"
      ? this.props.mockServiceApi.updateRestHandler(
          this.props.projectId,
          this.props.serviceId,
          this.props.handlerId ?? 0,
          form
        )
      : this.props.mockServiceApi
          .addRestHandler(this.props.projectId, this.props.serviceId, form)
          .then(({ data }) => {
            this.props.setClosed();
          })
          .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  componentDidMount(): void {
    this.setState({
      httpMethodInput: this.props.httpMethodInputInit ?? "",
      urlInput: this.props.urlInputInit ?? "",
    });
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onOk={() => {
          this.submitData();
        }}
        onCancel={() => this.props.setClosed()}
      >
        <Form layout="vertical">
          <Form.Item label="URL">
            <Input
              type="text"
              defaultValue={this.props.urlInputInit}
              onChange={(e) => this.setState({ urlInput: e.target.value })}
            />
          </Form.Item>
          <Form.Item label="HTTP-метод">
            <Select
              defaultValue={this.props.httpMethodInputInit}
              options={["GET", "POST", "PUT", "PATCH", "DELETE"].map(
                (method) => {
                  return {
                    value: method,
                    label: method,
                  };
                }
              )}
              onChange={(value) => this.setState({ httpMethodInput: value })}
            />
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default RestRouteUpdateModal;

interface RestRouteUpdateModalProps {
  modalType: "update" | "create";
  setClosed: () => void;
  open: boolean;
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId?: number;

  urlInputInit?: string;
  httpMethodInputInit?: string;
}

interface RestRouteUpdateModalState {
  urlInput: string;
  httpMethodInput: string;
}
