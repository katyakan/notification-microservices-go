export interface BaseKafkaMessage<T = unknown> {
  id: string;
  type: string;
  payload: T;
  timestamp: number;
}

export interface NotificationMessage {
  chatId: number;
  text: string;
  messageId?: string;
}

export type KafkaMessage =
  | (BaseKafkaMessage<NotificationMessage> & { type: "notification" })
  | (BaseKafkaMessage<Record<string, unknown>> & { type: string });
