import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultStub,
  MockservicesApi,
} from "../api/api";
import { Button, Divider, Form, Input, Modal, Space } from "antd";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons";
import TextArea from "antd/es/input/TextArea";
import { $notify, ENotifyKind } from "../notifier";

class EditRestStubBehaviorModal extends Component<
  EditRestStubBehaviorModalProps,
  EditRestStubBehaviorModalState
> {
  constructor(props: EditRestStubBehaviorModalProps) {
    super(props);

    this.state = {
      reqUrlParams: this.props.behavior.url_params ?? {},
      reqHeaders: this.props.behavior.headers ?? {},
      reqQueryParams: this.props.behavior.query ?? {},
      addReqUrlParamsInput: "",
      addReqQueryParamsInput: "",
      addReqHeaderInput: "",
      reqBodyInput: this.props.behavior.body ?? "",
      addRespHeaderInput: "",
      respHeaders: this.props.behavior.response_headers ?? {},
      respStatusInput:
        this.props.behavior.response_status == undefined
          ? ""
          : this.props.behavior.response_status.toString(),
      respBodyInput: this.props.behavior.response_body ?? "",
    };
  }

  sendFormData() {
    const reqForm = {
      headers: this.state.reqHeaders,
      query: this.state.reqQueryParams,
      url_params: this.state.reqUrlParams,
      body: this.state.reqBodyInput,

      response_status: parseInt(this.state.respStatusInput),
      response_headers: this.state.respHeaders,
      response_body: this.state.respBodyInput,
    };

    (this.props.modalType == "update"
      ? this.props.mockServiceApi.updateRestStubBehavior(
          this.props.projectId,
          this.props.serviceId,
          this.props.handlerId,
          this.props.behavior.id!,
          reqForm
        )
      : this.props.mockServiceApi.addRestStubBehavior(
          this.props.projectId,
          this.props.serviceId,
          this.props.handlerId,
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
        open={this.props.open}
        onOk={() => this.sendFormData()}
        onCancel={() => this.props.setClosed()}
      >
        <Form layout="vertical">
          <Divider orientation="left">Редактировать запрос</Divider>
          <span>URL-параметры</span>
          {Object.entries(this.state.reqUrlParams).map(([key, value]) => {
            return (
              <Form.Item label={key}>
                <Space>
                  <Input
                    type="text"
                    onChange={(e) => {
                      let urlParamVals = this.state.reqUrlParams;
                      urlParamVals[key] = e.target.value;
                      this.setState({ reqUrlParams: urlParamVals });
                    }}
                    value={value}
                  />
                  <Button
                    icon={<DeleteOutlined />}
                    danger
                    onClick={() => {
                      let urlParamVals = this.state.reqUrlParams;
                      delete urlParamVals[key];
                      this.setState({ reqUrlParams: urlParamVals });
                    }}
                  />
                </Space>
              </Form.Item>
            );
          })}
          <Space>
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
                this.setState({
                  reqUrlParams: reqParams,
                  addReqUrlParamsInput: "",
                });
              }}
            />
          </Space>
          <br />
          <span>Query-параметры</span>
          {Object.entries(this.state.reqQueryParams).map(([key, value]) => {
            return (
              <Form.Item label={key}>
                <Space>
                  <Input
                    type="text"
                    onChange={(e) => {
                      let urlParamVals = this.state.reqQueryParams;
                      urlParamVals[key] = e.target.value;
                      this.setState({ reqQueryParams: urlParamVals });
                    }}
                    value={value}
                  />
                  <Button
                    icon={<DeleteOutlined />}
                    danger
                    onClick={() => {
                      let urlParamVals = this.state.reqQueryParams;
                      delete urlParamVals[key];
                      this.setState({ reqQueryParams: urlParamVals });
                    }}
                  />
                </Space>
              </Form.Item>
            );
          })}
          <Space>
            <Input
              placeholder="Добавить Query-параметр"
              type="text"
              onChange={(e) =>
                this.setState({ addReqQueryParamsInput: e.target.value })
              }
            />
            <Button
              icon={<PlusOutlined />}
              type="default"
              onClick={() => {
                let reqParams = this.state.reqQueryParams;
                reqParams[this.state.addReqQueryParamsInput] = "";
                this.setState({
                  reqQueryParams: reqParams,
                  addReqQueryParamsInput: "",
                });
              }}
            />
          </Space>
          <br />
          <span>Заголовки</span>
          {Object.entries(this.state.reqHeaders).map(([key, value]) => {
            return (
              <Form.Item label={key}>
                <Space>
                  <Input
                    type="text"
                    onChange={(e) => {
                      let headerVals = this.state.reqHeaders;
                      headerVals[key] = e.target.value;
                      this.setState({ reqHeaders: headerVals });
                    }}
                    value={value}
                  />
                  <Button
                    icon={<DeleteOutlined />}
                    danger
                    onClick={() => {
                      let headerVals = this.state.reqHeaders;
                      delete headerVals[key];
                      this.setState({ reqHeaders: headerVals });
                    }}
                  />
                </Space>
              </Form.Item>
            );
          })}
          <Space>
            <Input
              placeholder="Добавить заголовок"
              type="text"
              onChange={(e) =>
                this.setState({ addReqHeaderInput: e.target.value })
              }
            />
            <Button
              icon={<PlusOutlined />}
              type="default"
              onClick={() => {
                let reqParams = this.state.reqHeaders;
                reqParams[this.state.addReqHeaderInput] = "";
                this.setState({
                  reqHeaders: reqParams,
                  addReqHeaderInput: "",
                });
              }}
            />
          </Space>
          <Form.Item label="Тело">
            <TextArea
              onChange={(e) => {
                this.setState({ reqBodyInput: e.target.value });
              }}
              value={this.state.reqBodyInput}
            />
          </Form.Item>

          <Divider orientation="left">Редактировать ответ</Divider>

          <Form.Item label="Статус ответа">
            <Input
              type="text"
              value={this.state.respStatusInput}
              onChange={(e) => {
                this.setState({ respStatusInput: e.target.value });
              }}
            />
          </Form.Item>
          <br />
          <span>Заголовки</span>
          {Object.entries(this.state.respHeaders).map(([key, value]) => {
            return (
              <Form.Item label={key}>
                <Space>
                  <Input
                    type="text"
                    onChange={(e) => {
                      let headers = this.state.respHeaders;
                      headers[key] = e.target.value;
                      this.setState({ respHeaders: headers });
                    }}
                    value={value}
                  />
                  <Button
                    icon={<DeleteOutlined />}
                    danger
                    onClick={() => {
                      let headers = this.state.respHeaders;
                      delete headers[key];
                      this.setState({ respHeaders: headers });
                    }}
                  />
                </Space>
              </Form.Item>
            );
          })}
          <Space>
            <Input
              placeholder="Добавить заголовок"
              type="text"
              onChange={(e) =>
                this.setState({ addRespHeaderInput: e.target.value })
              }
            />
            <Button
              icon={<PlusOutlined />}
              type="default"
              onClick={() => {
                let reqParams = this.state.respHeaders;
                reqParams[this.state.addRespHeaderInput] = "";
                this.setState({
                  respHeaders: reqParams,
                  addRespHeaderInput: "",
                });
              }}
            />
          </Space>
          <Form.Item label="Тело">
            <TextArea
              onChange={(e) => {
                this.setState({ respBodyInput: e.target.value });
              }}
              value={this.state.respBodyInput}
            />
          </Form.Item>
        </Form>
      </Modal>
    );
  }
}

export default EditRestStubBehaviorModal;

interface EditRestStubBehaviorModalProps {
  modalType: "update" | "create";

  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId: number;
  behavior: MockserviceListRestBehaviorResultStub;
  open: boolean;
  setClosed: () => void;
  refetch: () => void;
}

interface EditRestStubBehaviorModalState {
  reqUrlParams: { [key: string]: string };
  reqQueryParams: { [key: string]: string };
  reqHeaders: { [key: string]: string };
  addReqUrlParamsInput: string;
  addReqQueryParamsInput: string;
  addReqHeaderInput: string;
  reqBodyInput: string;

  respStatusInput: string;
  respHeaders: { [key: string]: string };
  addRespHeaderInput: string;
  respBodyInput: string;
}
