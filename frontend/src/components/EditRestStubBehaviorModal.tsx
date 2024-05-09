import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { Button, Divider, Form, Input, Modal } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import TextArea from "antd/es/input/TextArea";

class EditRestStubBehaviorModal extends Component<
  EditRestStubBehaviorModalProps,
  EditRestStubBehaviorModalState
> {
  constructor(props: EditRestStubBehaviorModalProps) {
    super(props);

    this.state = {
      reqUrlParams: {},
      addReqUrlParamsInput: "",
      reqBodyInput: "",
    };
  }

  updateStubBehavior() {
    this.props.mockServiceApi.updateRestStubBehavior(
      this.props.projectId,
      this.props.serviceId,
      this.props.handlerId,
      this.props.behaviorId,
      {
        body: this.state.reqBodyInput,
        // todo: остальные поля
      }
    );
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onOk={() => this.updateStubBehavior()}
        onCancel={() => this.props.setClosed()}
      >
        <Form layout="vertical">
          <Divider orientation="left">Редактировать запрос</Divider>
          <Form.Item label="URL-параметры">
            {Object.entries(this.state.reqUrlParams).map(([key, value]) => {
              return (
                <Input
                  type="text"
                  onChange={(e) => {
                    let urlParamVals = this.state.reqUrlParams;
                    urlParamVals[key] = e.target.value;
                    this.setState({ reqUrlParams: urlParamVals });
                  }}
                />
              );
            })}
            <Input.Group>
              <Input
                placeholder="Добавить URL-параметр"
                type="text"
                onChange={(e) =>
                  this.setState({ addReqUrlParamsInput: e.target.value })
                }
              />
              <Button
                icon={<PlusOutlined />}
                type="default"
                onClick={() => {
                  let reqParams = this.state.reqUrlParams;
                  reqParams[this.state.addReqUrlParamsInput] = "";
                  this.setState({ reqUrlParams: reqParams });
                }}
              />
            </Input.Group>
            <Form.Item label="Тело запроса">
              <TextArea
                onChange={(e) => {
                  this.setState({ reqBodyInput: e.target.value });
                }}
              />
            </Form.Item>
          </Form.Item>
          <Divider orientation="left">Редактировать ответ</Divider>
        </Form>
      </Modal>
    );
  }
}

export default EditRestStubBehaviorModal;

interface EditRestStubBehaviorModalProps {
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId: number;
  behaviorId: number;
  open: boolean;
  setClosed: () => void;
}

interface EditRestStubBehaviorModalState {
  reqUrlParams: { [key: string]: string };
  addReqUrlParamsInput: string;
  reqBodyInput: string;
}
