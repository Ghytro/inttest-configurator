import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { Checkbox, Input, Select, Space, message, Modal, Form } from "antd";
import { $notify, ENotifyKind } from "../notifier";
import TextArea from "antd/es/input/TextArea";

export default class EditRedisPubSubStubBehaviorModal extends Component<
  EditRedisPubSubStubBehaviorModalProps,
  EditRedisPubSubStubBehaviorModalState
> {
  constructor(props: EditRedisPubSubStubBehaviorModalProps) {
    super(props);

    this.state = {
      stubBody: "",
      intervalModulus: 0,
      intervalUnits: "",
      sendImmediately: false,
    };
  }

  async submitData() {
    const reqBody = {
      interval:
        this.state.intervalModulus.toString() + this.state.intervalUnits,
      send_immediately: this.state.sendImmediately,
      payload: this.state.stubBody,
    };
    try {
      await (this.props.modalType == "create"
        ? this.props.mockServiceApi.createRedisPubSubConstTopicGenerator(
            this.props.projectId,
            this.props.brokerId,
            this.props.topicId,
            reqBody
          )
        : this.props.mockServiceApi.updateRedisPubSubConstTopicGenerator(
            this.props.projectId,
            this.props.brokerId,
            this.props.topicId,
            this.props.generatorId!,
            reqBody
          ));
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    } finally {
      this.props.setClosed();
    }
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onCancel={() => {
          this.props.setClosed();
        }}
        onOk={() => {
          this.submitData();
        }}
      >
        <Form layout="vertical">
          <Form.Item label="Полезная нагрузка">
            <TextArea
              onChange={(e) => {
                this.setState({ stubBody: e.target.value });
              }}
            />
          </Form.Item>
          <Form.Item label="Интервал">
            <Space>
              <Input
                type="text"
                onChange={(e) => {
                  const newVal = e.target.value;
                  if (newVal == "") {
                    this.setState({ intervalModulus: 0 });
                    return;
                  }
                  const numbers = "0123456789";
                  for (let i = 0; i < newVal.length; ++i) {
                    if (!numbers.includes(newVal[i])) {
                      message.error("Интервал должен содержать только числа");
                      return;
                    }
                  }
                  this.setState({ intervalModulus: parseInt(e.target.value) });
                }}
                value={this.state.intervalModulus.toString()}
              />
              <Select
                style={{ width: "75px" }}
                options={["s", "ms", "ns"].map((i) => {
                  return {
                    value: i,
                    label: i,
                  };
                })}
              />
            </Space>
          </Form.Item>
          <Form.Item>
            <Checkbox
              onChange={(e) => {
                this.setState({ sendImmediately: e.target.checked });
              }}
            >
              Отправить моментально при запуске среды выполнения
            </Checkbox>
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

interface EditRedisPubSubStubBehaviorModalProps {
  open: boolean;
  setClosed: () => void;

  mockServiceApi: MockservicesApi;
  modalType: "create" | "update";
  projectId: number;
  brokerId: string;
  topicId: number;
  generatorId?: number;
}

interface EditRedisPubSubStubBehaviorModalState {
  stubBody: string;
  intervalModulus: number;
  intervalUnits: string;
  sendImmediately: boolean;
}
