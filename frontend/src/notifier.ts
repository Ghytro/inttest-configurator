import { message } from "antd";
import { TypeOpen } from "antd/es/message/interface";

enum ENotifyKind {
    INFO,
    WARNING,
    ERROR,
}

function $notify(kind: ENotifyKind, e: any) {
    let notifier: TypeOpen = message.info;
    switch (kind) {
    case ENotifyKind.INFO:
        notifier = message.info
    case ENotifyKind.WARNING:
        notifier = message.warning;
    case ENotifyKind.ERROR:
        notifier = message.error;
    }
    notifier(e.response.data.message);
}

export { $notify, ENotifyKind };