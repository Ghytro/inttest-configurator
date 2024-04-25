import { Button, Form, GetProp, Input, Modal, Upload } from "antd";
import React, { Component } from "react";
import { ProjectsApi } from "../api/api";
import { FileAddOutlined } from "@ant-design/icons";
import { RcFile, UploadFile, UploadProps } from "antd/es/upload";
import { $notify, ENotifyKind } from "../notifier";

class ProjCreateModal extends Component<
  ProjCreateModalProps,
  ProjCreateModalState
> {
  private projApi: ProjectsApi;
  constructor(props: ProjCreateModalProps) {
    super(props);

    this.projApi = props.projectsApi;
    this.state = {
      projNameInputValue: "",
      projDescInputValue: undefined,
      projUploadedFile: undefined,
      loading: false,
    };
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onCancel={() => {
          this.props.setClosed();
        }}
        onOk={() => {
          this.createProject();
        }}
        confirmLoading={this.state.loading}
      >
        <Form layout="vertical">
          <Form.Item label="Наименование">
            <Input
              onChange={(e) => {
                this.setState({ projNameInputValue: e.target.value });
              }}
            />
          </Form.Item>
          <Form.Item label="Описание">
            <Input
              onChange={(e) => {
                this.setState({
                  projDescInputValue:
                    e.target.value == "" ? undefined : e.target.value,
                });
              }}
            />
          </Form.Item>
          <Form.Item label="Файл проекта">
            <Upload
              beforeUpload={(file) => {
                this.setState({ projUploadedFile: [file] });
                return false;
              }}
              onRemove={(file) => {
                this.setState({ projUploadedFile: [] });
              }}
              fileList={this.state.projUploadedFile}
            >
              <Button icon={<FileAddOutlined />}>
                Прикрепить конфигурацию
              </Button>
            </Upload>
          </Form.Item>
        </Form>
      </Modal>
    );
  }

  createProject() {
    this.setState({ loading: true });
    type FileType = Parameters<GetProp<UploadProps, "beforeUpload">>[0];
    const formFile =
      (this.state.projUploadedFile ?? []).length > 0
        ? (this.state.projUploadedFile![0] as FileType)
        : undefined;
    this.projApi
      .createProject(
        this.state.projNameInputValue,
        formFile,
        this.state.projDescInputValue
      )
      .then(({ data }) => {
        this.setState({
          projNameInputValue: "",
          projDescInputValue: undefined,
          projUploadedFile: undefined,
        });
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e))
      .finally(() => this.setState({ loading: false }));
  }
}

export default ProjCreateModal;

interface ProjCreateModalProps {
  projectsApi: ProjectsApi;
  setClosed: () => void;
  open: boolean;
}

interface ProjCreateModalState {
  projNameInputValue: string;
  projDescInputValue?: string;
  projUploadedFile?: UploadFile[];
  loading: boolean;
}
