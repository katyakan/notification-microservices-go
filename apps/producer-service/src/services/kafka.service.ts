import { Injectable, OnModuleInit, OnModuleDestroy } from "@nestjs/common";
import { Kafka, Producer, ProducerRecord } from "kafkajs";
import { v4 as uuidv4 } from "uuid";
import { kafkaConfig } from "../config/kafka.config";
import { CreateMessageDto } from "../dto/create-message.dto";
import { KafkaMessage } from "@libs/shared";

@Injectable()
export class KafkaService implements OnModuleInit, OnModuleDestroy {
  private readonly kafka: Kafka;
  private readonly producer: Producer;

  constructor() {
    this.kafka = new Kafka({
      clientId: kafkaConfig.clientId,
      brokers: kafkaConfig.brokers,
      retry: kafkaConfig.retry,
    });

    this.producer = this.kafka.producer();
  }

  async onModuleInit() {
    await this.producer.connect();
  }

  async onModuleDestroy() {
    await this.producer.disconnect();
  }

  async sendMessage(dto: CreateMessageDto): Promise<{ id: string }> {
    const message: KafkaMessage = {
      id: uuidv4(),
      type: dto.type,
      payload: dto.payload,
      timestamp: Date.now(),
    };

    const record: ProducerRecord = {
      topic: kafkaConfig.topics.notifications,
      messages: [
        {
          key: message.id,
          value: JSON.stringify(message),
          headers: {
            'message-type': dto.type,
          },
        },
      ],
    };

    try {
      await this.producer.send(record);
      return { id: message.id };
    } catch (error) {
      throw new Error(`Failed to send message: ${error.message}`);
    }
  }
}