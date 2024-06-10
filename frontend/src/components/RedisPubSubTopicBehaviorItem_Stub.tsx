import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { Button, Descriptions, Space } from "antd";
import Modal from "antd/es/modal/Modal";
import React, { Component } from "react";
import EditRedisPubSubStubBehaviorModal from "./EditRedisPubSubStubBehaviorModal";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

export default class RedisPubSubTopicBehaviorItem_Stub extends Component<
  RedisPubSubTopicBehaviorItem_StubProps,
  RedisPubSubTopicBehaviorItem_StubState
> {
  constructor(props: RedisPubSubTopicBehaviorItem_StubProps) {
    super(props);

    this.state = {
      editStubBehaviorModalOpen: false,
      showConstPayloadOpen: false,
    };
  }

  async deleteBehavior() {
    try {
      await this.props.mockServiceApi.deleteRedisPubSubGenerator(
        this.props.projectId,
        this.props.brokerId,
        this.props.topicId,
        this.props.generatorId
      );
      this.props.refetch();
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    }
  }

  render(): React.ReactNode {
    return (
      <>
        <Descriptions
          title="Описание запроса"
          extra={
            <Space>
              <Button
                type="default"
                icon={<EditOutlined />}
                onClick={() =>
                  this.setState({ editStubBehaviorModalOpen: true })
                }
              >
                Редактировать
              </Button>
              <Button
                danger
                icon={<DeleteOutlined />}
                onClick={() => this.deleteBehavior()}
              >
                Удалить
              </Button>
            </Space>
          }
        >
          <Descriptions.Item label="Полезная нагрузка">
            <Button
              type="link"
              onClick={() => {
                this.setState({ showConstPayloadOpen: true });
              }}
            >
              Показать
            </Button>
            <Modal
              open={this.state.showConstPayloadOpen}
              onOk={() => this.setState({ showConstPayloadOpen: false })}
              onCancel={() => this.setState({ showConstPayloadOpen: false })}
            >
              {this.props.payload}
            </Modal>
          </Descriptions.Item>
          <Descriptions.Item label="Интервал">
            {this.props.intervalModulus.toString() + this.props.intervalUnits}
          </Descriptions.Item>
          <Descriptions.Item label="Отправить моментально">
            {this.props.sendImmediately ? "да" : "нет"}
          </Descriptions.Item>
        </Descriptions>
        <EditRedisPubSubStubBehaviorModal
          modalType="update"
          brokerId={this.props.brokerId}
          mockServiceApi={this.props.mockServiceApi}
          open={this.state.editStubBehaviorModalOpen}
          projectId={this.props.projectId}
          setClosed={() => {
            this.setState({ editStubBehaviorModalOpen: false });
          }}
          topicId={this.props.topicId}
          generatorId={this.props.generatorId}
        />
      </>
    );
  }
}

interface RedisPubSubTopicBehaviorItem_StubProps {
  intervalModulus: number;
  intervalUnits: string;
  sendImmediately: boolean;
  payload: string;

  brokerId: string;
  mockServiceApi: MockservicesApi;
  projectId: number;
  topicId: number;
  generatorId: number;
  refetch: () => void;
}

interface RedisPubSubTopicBehaviorItem_StubState {
  editStubBehaviorModalOpen: boolean;
  showConstPayloadOpen: boolean;
}
