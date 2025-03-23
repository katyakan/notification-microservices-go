export interface KafkaMessage {
  id: string;
  type: string;
  payload: any;
  timestamp: number;
}

export interface NotificationMessage {
  chatId: number;
  text: string;
  messageId?: string;
}
