import { Module } from "@nestjs/common";
import { ProducerController } from "./controllers/producer.controller";
import { KafkaService } from "./services/kafka.service";

@Module({
  imports: [],
  controllers: [ProducerController],
  providers: [KafkaService],
})
export class ProducerModule {}
