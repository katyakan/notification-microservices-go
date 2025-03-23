import { Module } from "@nestjs/common";
import { ConsumerController } from "./controllers/consumer.controller";
import { KafkaService } from "./services/kafka.service";

@Module({
  imports: [],
  controllers: [ConsumerController],
  providers: [KafkaService],
})
export class ConsumerModule {}
