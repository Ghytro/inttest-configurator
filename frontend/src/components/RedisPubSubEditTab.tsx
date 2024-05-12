import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, FloatButton, Space } from "antd";
import React, { Component } from "react";
import CreateRedisPubSubDialog from "./CreateRedisPubSubDialog";

class RedisPubSubEditTab extends Component<
  RedisPubSubEditTabProps,
  RedisPubSubEditTabState
> {
  constructor(props: RedisPubSubEditTabProps) {
    super(props);

    this.state = {
      addTopicDialogOpen: false,
      brokerEditDialogOpen: false,
    };
  }

  deleteBroker() {}

  render(): React.ReactNode {
    return (
      <>
        <Space>
          <Button
            type="default"
            icon={<EditOutlined />}
            onClick={() => {
              this.setState({ brokerEditDialogOpen: true });
            }}
          >
            Редактировать брокер
          </Button>
          <Button
            danger={true}
            icon={<DeleteOutlined />}
            onClick={() => {
              this.deleteBroker();
            }}
          >
            Удалить брокер
          </Button>
        </Space>
        <FloatButton
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => this.setState({ addTopicDialogOpen: true })}
        />
      </>
    );
  }
}

export default RedisPubSubEditTab;

interface RedisPubSubEditTabProps {}

interface RedisPubSubEditTabState {
  brokerEditDialogOpen: boolean;
  addTopicDialogOpen: boolean;
}
