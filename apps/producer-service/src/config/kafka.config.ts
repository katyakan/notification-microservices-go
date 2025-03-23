export const kafkaConfig = {
  brokers: process.env.KAFKA_BROKERS?.split(",") || ["localhost:9092"],
  clientId: "producer-service",
  retry: {
    initialRetryTime: 100,
    retries: 8,
  },
  topics: {
    notifications: "notifications",
    deadLetter: "dead-letter",
  },
};
