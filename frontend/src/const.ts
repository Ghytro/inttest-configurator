const restServiceType = "rest";
const soapServiceType = "soap";
const redisServiceType = "redis-pubsub";

const translateServiceType = {
  [restServiceType]: "REST API",
  [soapServiceType]: "SOAP API",
  [redisServiceType]: "Redis PubSub",
};

const stubBehaviorType = "stub";
const mockBehaviorType = "mock";

export {restServiceType, soapServiceType, redisServiceType, translateServiceType, stubBehaviorType, mockBehaviorType};