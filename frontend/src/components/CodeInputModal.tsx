import { Modal } from "antd";
import React, { useState } from "react";
import MonacoEditor from "react-monaco-editor";

const CodeInputModal: React.FC<CodeInputModalProps> = (props) => {
  const [code, setCode] = useState(props.code);

  return (
    <Modal
      open={props.open}
      onCancel={() => {
        setCode(props.code);
        props.setClosed();
      }}
      onOk={() => {
        props.submit(code);
      }}
      width="80%"
    >
      <MonacoEditor
        width="90%"
        height="600"
        language="python"
        value={code}
        onChange={(val, e) => {
          setCode(val);
        }}
      />
    </Modal>
  );
};

export default CodeInputModal;

interface CodeInputModalProps {
  open: boolean;
  setClosed: () => void;
  submit: (code: string) => void;
  code: string;
}
