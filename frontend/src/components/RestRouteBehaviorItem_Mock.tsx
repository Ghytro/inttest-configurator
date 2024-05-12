import React, { Component } from "react";
import {
  MockserviceListRestBehaviorResultMock,
  MockservicesApi,
} from "../api/api";
import { Button, Form } from "antd";
import CodeInputModal from "./CodeInputModal";
import { $notify, ENotifyKind } from "../notifier";

class RestRouteBehaviorItem_Mock extends Component<
  RestRouteBehaviorItem_MockProps,
  RestRouteBehaviorItem_MockState
> {
  constructor(props: RestRouteBehaviorItem_MockProps) {
    super(props);

    this.state = {
      codeModalOpen: false,
    };
  }

  updateMockBehavior(code: string) {
    this.props.mockServiceApi
      .updateRestMockBehavior(
        this.props.projectId,
        this.props.serviceId,
        this.props.handlerId,
        this.props.behavior.id!,
        {
          impl: code.split("\n"),
        }
      )
      .then(({ data }) => {})
      .catch((e) => $notify(ENotifyKind.ERROR, e));
  }

  render(): React.ReactNode {
    return (
      <>
        <Form layout="vertical">
          <Form.Item label="Имплементация (Python)">
            <Button
              type="link"
              onClick={() => {
                this.setState({ codeModalOpen: true });
              }}
            >
              Показать
            </Button>
            <CodeInputModal
              open={this.state.codeModalOpen}
              code={(this.props.behavior.impl ?? []).join("\n")}
              setClosed={() => {
                this.props.refetch();
                this.setState({ codeModalOpen: false });
              }}
              submit={(code) => {
                this.updateMockBehavior(code);
                this.props.refetch();
              }}
            />
          </Form.Item>
        </Form>
      </>
    );
  }
}

export default RestRouteBehaviorItem_Mock;

interface RestRouteBehaviorItem_MockProps {
  behavior: MockserviceListRestBehaviorResultMock;
  refetch: () => void;
  mockServiceApi: MockservicesApi;
  projectId: number;
  serviceId: string;
  handlerId: number;
}

interface RestRouteBehaviorItem_MockState {
  codeModalOpen: boolean;
}
