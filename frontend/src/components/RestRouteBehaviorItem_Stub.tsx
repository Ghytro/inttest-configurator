import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultStub,
  MockservicesApi,
} from "../api/api";
import { Button, Descriptions, Divider, Modal, Space } from "antd";
import EditRestStubBehaviorModal from "./EditRestStubBehaviorModal";
import { $notify, ENotifyKind } from "../notifier";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";

class RestRouteBehaviorItem_Stub extends Component<
  RestRouteBehaviorItem_StubProps,
  RestRouteBehaviorItem_StubState
> {
  constructor(props: RestRouteBehaviorItem_StubProps) {
    super(props);

    this.state = {
      reqBodyModalOpen: false,
      respBodyModalOpen: false,
      editStubBehaviorModalOpen: false,
    };
  }

  renderKV(obj: { [key: string]: string }): React.ReactNode {
    return (
      <>
        {Object.entries(obj).map(([key, value]) => {
          return (
            <>
              <span>{`${key}: ${value}`}</span>
              <br />
            </>
          );
        })}
      </>
    );
  }

  deleteBehavior() {
    this.props.mockServiceApi
      .deleteRestBehavior(
        this.props.projectId,
        this.props.serviceId,
        this.props.handlerId,
        this.props.behavior.id!
      )
      .then(({ data }) => {
        this.props.refetch();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
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
          <Descriptions.Item label="URL-параметры">
            {this.renderKV(this.props.behavior.url_params!)}
          </Descriptions.Item>
          <Descriptions.Item label="Query-параметры">
            {this.renderKV(this.props.behavior.query!)}
          </Descriptions.Item>
          <Descriptions.Item label="Заголовки">
            {this.renderKV(this.props.behavior.headers!)}
          </Descriptions.Item>
          <Descriptions.Item label="Тело запроса">
            <Button
              type="link"
              onClick={() => this.setState({ reqBodyModalOpen: true })}
            >
              Показать
            </Button>
            <Modal
              open={this.state.reqBodyModalOpen}
              onOk={() => this.setState({ reqBodyModalOpen: false })}
              onCancel={() => this.setState({ reqBodyModalOpen: false })}
            >
              {JSON.stringify(
                this.props.behavior.body! != ""
                  ? JSON.parse(this.props.behavior.body!)
                  : "",
                null,
                "\t"
              )}
            </Modal>
          </Descriptions.Item>
        </Descriptions>

        <Divider />

        <Descriptions title="Описание ответа">
          <Descriptions.Item label="Статус ответа">
            {this.props.behavior.response_status}
          </Descriptions.Item>
          <Descriptions.Item label="Заголовки ответа">
            {this.renderKV(this.props.behavior.response_headers!)}
          </Descriptions.Item>
          <Descriptions.Item label="Тело ответа">
            <Button
              type="link"
              onClick={() => this.setState({ respBodyModalOpen: true })}
            >
              Показать
            </Button>
            <Modal
              open={this.state.respBodyModalOpen}
              onOk={() => this.setState({ respBodyModalOpen: false })}
              onCancel={() => this.setState({ respBodyModalOpen: false })}
            >
              {JSON.stringify(
                this.props.behavior.response_body! != ""
                  ? JSON.parse(this.props.behavior.response_body!)
                  : "",
                null,
                "\t"
              )}
            </Modal>
          </Descriptions.Item>
        </Descriptions>
        <EditRestStubBehaviorModal
          modalType="update"
          mockServiceApi={this.props.mockServiceApi}
          behavior={this.props.behavior}
          handlerId={this.props.handlerId}
          projectId={this.props.projectId}
          serviceId={this.props.serviceId}
          open={this.state.editStubBehaviorModalOpen}
          setClosed={() => {
            this.setState({ editStubBehaviorModalOpen: false });
            this.props.refetch();
          }}
          refetch={() => this.props.refetch()}
        />
      </>
    );
  }
}

export default RestRouteBehaviorItem_Stub;

interface RestRouteBehaviorItem_StubProps {
  behavior: MockserviceListRestBehaviorResultStub;

  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId: number;
  refetch: () => void;
}

interface RestRouteBehaviorItem_StubState {
  reqBodyModalOpen: boolean;
  respBodyModalOpen: boolean;
  editStubBehaviorModalOpen: boolean;
}
