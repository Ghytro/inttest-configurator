import React, { Component } from "react";
import { MockserviceListRestBehaviorResultStub } from "../api/api";
import { Button, Descriptions, Divider, Modal } from "antd";
import EditRestStubBehaviorModal from "./EditRestStubBehaviorModal";

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

  render(): React.ReactNode {
    return (
      <>
        <Descriptions
          title="Описание запроса"
          extra={
            <Button
              type="primary"
              onClick={() => this.setState({ editStubBehaviorModalOpen: true })}
            >
              Редактировать
            </Button>
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
                JSON.parse(this.props.behavior.body!),
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
                JSON.parse(this.props.behavior.response_body!),
                null,
                "\t"
              )}
            </Modal>
          </Descriptions.Item>
        </Descriptions>
      </>
    );
  }
}

export default RestRouteBehaviorItem_Stub;

interface RestRouteBehaviorItem_StubProps {
  behavior: MockserviceListRestBehaviorResultStub;
}

interface RestRouteBehaviorItem_StubState {
  reqBodyModalOpen: boolean;
  respBodyModalOpen: boolean;
  editStubBehaviorModalOpen: boolean;
}
