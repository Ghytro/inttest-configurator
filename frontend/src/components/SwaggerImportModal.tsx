import { FileAddOutlined } from "@ant-design/icons";
import { Button, GetProp, Modal, Upload, UploadFile, UploadProps } from "antd";
import React, { Component } from "react";
import { MockservicesApi } from "../api/api";
import { $notify, ENotifyKind } from "../notifier";

export default class SwaggerImportModal extends Component<
  SwaggerImportModalProps,
  SwaggerImportModalState
> {
  constructor(props: SwaggerImportModalProps) {
    super(props);
    this.state = {};
  }
  submitData() {
    type FileType = Parameters<GetProp<UploadProps, "beforeUpload">>[0];
    const formFile =
      (this.state.swaggerUploadedFile ?? []).length > 0
        ? (this.state.swaggerUploadedFile![0] as FileType)
        : undefined;
    this.props.mockServiceApi
      .importSwagger(this.props.projectId, this.props.serviceId, formFile)
      .then(({ data }) => {
        this.setState({ swaggerUploadedFile: undefined });
        this.props.setClosed();
      })
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <Modal
        onCancel={() => {
          this.props.setClosed();
        }}
        onOk={() => {
          this.submitData();
        }}
        open={this.props.open}
      >
        <Upload
          beforeUpload={(file) => {
            this.setState({ swaggerUploadedFile: [file] });
            return false;
          }}
          onRemove={(file) => {
            this.setState({ swaggerUploadedFile: [] });
          }}
          fileList={this.state.swaggerUploadedFile}
        >
          <Button icon={<FileAddOutlined />}>Прикрепить конфигурацию</Button>
        </Upload>
      </Modal>
    );
  }
}

interface SwaggerImportModalProps {
  open: boolean;
  setClosed: () => void;
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
}

interface SwaggerImportModalState {
  swaggerUploadedFile?: UploadFile[];
}
