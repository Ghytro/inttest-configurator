import { Modal } from "antd";
import React, { Component } from "react";
import MonacoEditor from "react-monaco-editor";
import { MockserviceListRestRoutesResult } from "../api/api";

class CodeInputModal extends Component<
  CodeInputModalProps,
  CodeInputModalState
> {
  constructor(props: CodeInputModalProps) {
    super(props);

    this.state = {
      code: props.code ?? "",
    };
  }

  componentDidMount(): void {
    this.setState({ code: this.props.code ?? "" });
  }

  render(): React.ReactNode {
    return (
      <Modal
        open={this.props.open}
        onCancel={() => {
          this.setState({ code: this.props.code ?? "" });
          this.props.setClosed();
        }}
        onOk={() => {
          this.props.submit(this.state.code);
        }}
        width="80%"
        title={this.props.title}
      >
        <MonacoEditor
          width="90%"
          height="600"
          language="python"
          defaultValue={this.props.code}
          value={this.state.code}
          onChange={(val, e) => {
            this.setState({ code: val });
          }}
        />
      </Modal>
    );
  }
}

export default CodeInputModal;

interface CodeInputModalProps {
  open: boolean;
  setClosed: () => void;
  submit: (code: string) => void;
  code?: string;
  title?: string;
}

interface CodeInputModalState {
  code: string;
}

function makeInitialCodeImpl_Http(
  serviceId: string,
  handlerData: MockserviceListRestRoutesResult,
  priority: number
): string {
  const capitalize = (str: string): string => {
    return str.charAt(0).toUpperCase() + str.slice(1);
  };
  let functionName =
    serviceId +
    handlerData
      .route!.split("/")
      .filter((item) => item != "")
      .map((item) => capitalize(item.slice(+(item.charAt(0) == ":"))))
      .join("") +
    `${capitalize(handlerData.method!.toLowerCase())}_${priority}`;

  return `def ${functionName}(url_params, query_params, headers, body):\n\treturn {"response_object": "here"}`;
}

function makeInitialCodeImpl_Broker(
  serviceId: string,
  topic: string,
  priority: number
): string {
  return "";
}

export { makeInitialCodeImpl_Http, makeInitialCodeImpl_Broker };
