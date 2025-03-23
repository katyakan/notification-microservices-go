export const kafkaConfig = {
  brokers: process.env.KAFKA_BROKERS?.split(",") || ["localhost:9092"],
  clientId: "notification-service",
  groupId: "telegram-notification-group",
  retry: {
    initialRetryTime: 100,
    retries: 8,
  },
  topics: {
    notifications: "notifications",
    deadLetter: "dead-letter",
  },
};
