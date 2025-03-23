import { Injectable, OnModuleInit, OnModuleDestroy } from "@nestjs/common";
import { Kafka, Consumer, EachMessagePayload } from "kafkajs";
import { kafkaConfig } from "../config/kafka.config";
import { KafkaMessage, NotificationMessage, BaseKafkaMessage } from "@libs/shared";

@Injectable()
export class KafkaService implements OnModuleInit, OnModuleDestroy {
  private readonly kafka: Kafka;
  private readonly consumer: Consumer;

  constructor() {
    this.kafka = new Kafka({
      clientId: kafkaConfig.clientId,
      brokers: kafkaConfig.brokers,
      retry: kafkaConfig.retry,
    });

    this.consumer = this.kafka.consumer({ groupId: kafkaConfig.groupId });
  }

  async onModuleInit() {
    await this.consumer.connect();
    await this.consumer.subscribe({
      topic: kafkaConfig.topics.notifications,
      fromBeginning: true,
    });

    await this.consumer.run({
      eachMessage: async (payload: EachMessagePayload) => {
        try {
          if (!payload.message.value) {
            throw new Error("Empty message value");
          }

          const rawMessage = JSON.parse(payload.message.value.toString()) as unknown;
          if (!this.isValidKafkaMessage(rawMessage)) {
            throw new Error("Invalid message format");
          }

          await this.processMessage(rawMessage);
        } catch (error) {
          console.error("Error processing message:", error);
          await this.handleDeadLetter(payload);
        }
      },
    });
  }

  async onModuleDestroy() {
    await this.consumer.disconnect();
  }

  private isValidKafkaMessage(message: unknown): message is KafkaMessage {
    if (typeof message !== "object" || message === null) {
      return false;
    }

    const msg = message as Record<string, unknown>;
    return (
      typeof msg.id === "string" &&
      typeof msg.type === "string" &&
      typeof msg.timestamp === "number" &&
      msg.payload !== undefined &&
      typeof msg.payload === "object"
    );
  }

  private isNotificationMessage(
    message: KafkaMessage,
  ): message is BaseKafkaMessage<NotificationMessage> & { type: "notification" } {
    return (
      message.type === "notification" &&
      typeof message.payload === "object" &&
      message.payload !== null &&
      typeof (message.payload as NotificationMessage).chatId === "number" &&
      typeof (message.payload as NotificationMessage).text === "string"
    );
  }

  private async processMessage(message: KafkaMessage) {
    console.log("Processing message:", {
      id: message.id,
      type: message.type,
      timestamp: message.timestamp,
    });

    if (this.isNotificationMessage(message)) {
      await this.handleNotification(message.payload);
    } else {
      console.warn(`Unknown message type: ${message.type}`);
    }
  }

  private async handleNotification(payload: NotificationMessage) {
    const { chatId, text } = payload;
    console.log(`Sending notification to chat ${chatId}: ${text}`);
    await Promise.resolve();
    // TODO: Реализовать реальную отправку уведомления
    // await notificationClient.send({ chatId, text });
  }

  private async handleDeadLetter(payload: EachMessagePayload) {
    const deadLetterProducer = this.kafka.producer();
    await deadLetterProducer.connect();

    try {
      await deadLetterProducer.send({
        topic: kafkaConfig.topics.deadLetter,
        messages: [
          {
            key: payload.message.key ?? undefined,
            value: payload.message.value ?? null,
            headers: {
              ...payload.message.headers,
              "error-timestamp": Date.now().toString(),
            },
          },
        ],
      });
    } finally {
      await deadLetterProducer.disconnect();
    }
  }
}
