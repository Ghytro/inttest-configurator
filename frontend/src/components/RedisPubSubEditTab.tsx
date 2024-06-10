import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, FloatButton, Space } from "antd";
import React, { Component } from "react";
import { ServiceCascaderOptVal } from "./ServiceSelectionCascader";
import { MockservicesApi } from "../api/api";
import { projPageUrl } from "../routesEnum";
import { $notify, ENotifyKind } from "../notifier";
import RedisPubSubAddTopicDialog from "./RedisPubSubAddTopicDialog";
import CreateRedisPubSubDialog from "./CreateRedisPubSubDialog";
import RedisPubSubTopics from "./RedisPubSubTopics";

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

  deleteBroker() {
    this.props.mockServiceApi
      .deleteRedisPubSub(this.props.projectId, this.props.brokerData.id)
      .then(({ data }) => {
        window.location.href = projPageUrl(this.props.projectId);
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

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
        <RedisPubSubAddTopicDialog
          open={this.state.addTopicDialogOpen}
          setClosed={() => this.setState({ addTopicDialogOpen: false })}
          mockServiceApi={this.props.mockServiceApi}
          projectId={this.props.projectId}
          modalType="create"
          brokerId={this.props.brokerData.id}
        />
        <CreateRedisPubSubDialog
          modalType="update"
          mockServiceApi={this.props.mockServiceApi}
          open={this.state.brokerEditDialogOpen}
          projectId={this.props.projectId}
          setClosed={() => this.setState({ brokerEditDialogOpen: false })}
          refetch={() => {
            /* todo */
          }}
          brokerIdInputInitValue={this.props.brokerData.id.toString()}
          brokerPortInputInitValue={this.props.brokerData.port.toString()}
        />
        <RedisPubSubTopics
          mockServiceApi={this.props.mockServiceApi}
          projectId={this.props.projectId}
          brokerId={this.props.brokerData.id}
        />
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

interface RedisPubSubEditTabProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  brokerData: ServiceCascaderOptVal;
}

interface RedisPubSubEditTabState {
  brokerEditDialogOpen: boolean;
  addTopicDialogOpen: boolean;
}
