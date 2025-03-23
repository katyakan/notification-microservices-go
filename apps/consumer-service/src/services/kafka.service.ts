import { Injectable, OnModuleInit, OnModuleDestroy } from "@nestjs/common";
import { Kafka, Consumer, EachMessagePayload } from "kafkajs";
import { kafkaConfig } from "../config/kafka.config";
import { KafkaMessage } from "@libs/shared";

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

          const message = JSON.parse(payload.message.value.toString()) as KafkaMessage;
          await this.processMessage(message);
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

  private async processMessage(message: KafkaMessage) {
    console.log("Processing message:", {
      id: message.id,
      type: message.type,
      timestamp: message.timestamp,
    });

    switch (message.type) {
      case "notification":
        await this.handleNotification(message);
        break;
      default:
        console.warn(`Unknown message type: ${message.type}`);
    }
  }

  private async handleNotification(message: KafkaMessage) {
    const { chatId, text } = message.payload;
    console.log(`Sending notification to chat ${chatId}: ${text}`);
    // Здесь будет логика отправки уведомления
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
