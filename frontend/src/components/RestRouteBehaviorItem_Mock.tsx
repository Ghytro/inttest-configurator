import React, { Component } from "react";
import { MockserviceListRestBehaviorResultMock } from "../api/api";

class RestRouteBehaviorItem_Mock extends Component<
  RestRouteBehaviorItem_MockProps,
  RestRouteBehaviorItem_MockState
> {
  constructor(props: RestRouteBehaviorItem_MockProps) {
    super(props);
  }

  render(): React.ReactNode {}
}

export default RestRouteBehaviorItem_Mock;

interface RestRouteBehaviorItem_MockProps {
  behavior: MockserviceListRestBehaviorResultMock;
}

interface RestRouteBehaviorItem_MockState {}
