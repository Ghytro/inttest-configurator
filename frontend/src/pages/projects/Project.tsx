import React, { Component } from "react";
import { MockservicesApi } from "../../api/api";
import { withRouter } from "../../AppRoutes";
import { Params } from "react-router-dom";
import { $notify, ENotifyKind } from "../../notifier";

const restServiceType = "rest";
const soapServiceType = "soap";
const redisServiceType = "redis-pubsub";

const translateServiceType = {
  restServiceType: "REST API",
  soapServiceType: "SOAP API",
  redisServiceType: "Redis PubSub",
}

class ProjectPage extends Component<ProjectPageProps, ProjectPageState> {
  private mockServiceApi: MockservicesApi;

  constructor(props: ProjectPageProps) {
    super(props);

    this.mockServiceApi = new MockservicesApi();
    this.state = {
      loading: false,
    };
  }

  componentDidMount(): void {
    this.setState({loading: true});
    try {
      const resp = await this.mockServiceApi.listServices(parseInt(this.props.urlParams.id!));
      for (const [key, value] of Object.entries(resp.data)) {

      }
    } catch (e) {
      $notify(ENotifyKind.ERROR, e);
    } finally {
      this.setState({loading: false});
    }
  }

  render(): React.ReactNode {
    return (
        
    );
  }
}

export default withRouter(ProjectPage);

interface ProjectPageProps {
  urlParams: Readonly<Params<string>>;
}

interface ProjectPageState {
  loading: boolean;
  
}
